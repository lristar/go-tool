package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/lristar/go-tool/lib/stringx"
	"github.com/lristar/go-tool/server/sentry"
	"github.com/opentracing/opentracing-go"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"io"
	"io/ioutil"
	"moul.io/http2curl"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// RequestHook 新增的切面方式添加中间件
type RequestHook interface {
	Before(ctx context.Context, req *http.Request) error
	AfterProcess(ctx context.Context, resp *http.Response) error
}

// MyRequest http请求
type MyRequest struct {
	Method  string            // 请求方式POST,GET,PUT,DELETE
	URL     string            // 请求url
	Header  map[string]string // header
	Cookies map[string]string // cookies
	Data    interface{}       // 请求body
	Timeout time.Duration     // 超时时间
	Tr      *http.Transport   // 承载体
	Alert   bool              // 报警
	Ctx     context.Context
	hooks   []RequestHook
}

// NewRequest 创建一个新的request
func NewRequest() *MyRequest {
	// 设置默认中间件
	r := &MyRequest{Timeout: 120, Alert: true, Tr: new(http.Transport), hooks: make([]RequestHook, 0)}
	return r
}

// initHeaders 设置header
func (r *MyRequest) initHeaders(cli *http.Request) {
	cli.Header.Set("Content-Type", "application/json")
	for k, v := range r.Header {
		cli.Header.Set(k, v)
	}
}

func (r *MyRequest) initCookies(cli *http.Request) {
	for k, v := range r.Cookies {
		cli.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
}

// Check "x-www-form-urlencoded"
func (r *MyRequest) isX() bool {
	if len(r.Header) > 0 {
		for _, v := range r.Header {
			if strings.Contains(strings.ToLower(v), "x-www-form-urlencoded") {
				return true
			}
		}
	}
	return false
}

// Build body
func (r *MyRequest) buildBody() (io.Reader, error) {
	if r.Method == "GET" || r.Method == "DELETE" {
		return nil, nil
	}

	if r.Data == nil {
		return strings.NewReader(""), nil
	}

	if r.isX() {
		data := make([]string, 0)
		datas := r.Data.(map[string]interface{})
		for k, v := range datas {
			if s, ok := v.(string); ok {
				data = append(data, fmt.Sprintf("%s=%v", k, s))
				continue
			}
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			data = append(data, fmt.Sprintf("%s=%s", k, string(b)))
		}
		return strings.NewReader(strings.Join(data, "&")), nil
	}

	b, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// HTTPDo 发起http请求并返回结果或error
func (r *MyRequest) HTTPDo() ([]byte, error) {
	client := &http.Client{Timeout: r.Timeout * time.Second, Transport: r.Tr}
	if r.Method == "" {
		r.Method = "GET"
	}

	dataNew, err := r.buildBody()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.Method, r.URL, dataNew)
	if err != nil {
		return nil, err
	}

	r.initHeaders(req)
	r.initCookies(req)
	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	for _, hook := range r.hooks {
		if err := hook.Before(r.Ctx, req); err != nil {
			return nil, err
		}
	}
	var span opentracing.Span
	if r.Ctx != nil {
		c, ok := r.Ctx.Value("ctx").(context.Context)
		if ok {
			span, _ = opentracing.StartSpanFromContext(c, fmt.Sprintf("%s  %s", r.Method, r.URL))
			defer span.Finish()
			span.SetTag("http.req.method", r.Method)
			if r.Header != nil {
				span.SetTag("http.req.header", stringx.Marshal(r.Header))
			}
			if r.Data != nil {
				span.LogFields(tracerLog.String("http.req.body", stringx.Marshal(r.Data)))
			}
			span.SetTag("http.req.url", r.URL)
			span.LogFields(tracerLog.String("curl", command.String()))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	for _, hook := range r.hooks {
		if err := hook.AfterProcess(r.Ctx, resp); err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if span != nil {
		span.LogFields(tracerLog.String("http.res.body", string(body)))
		span.SetTag("http.res.status_code", resp.StatusCode)
	}
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		if span != nil {
			span.SetTag("error", true)
		}
		var errT map[string]interface{}
		err = json.Unmarshal(body, &errT)
		if err != nil {
			e := fmt.Errorf("[S=%d] err= %s", resp.StatusCode, string(body))
			if r.Alert {
				r.Sentry(e)
			}
			return nil, e
		}
		errKey := []string{"error_description", "message", "msg", "stack"}
		for _, key := range errKey {
			if v, ok := errT[key]; ok {
				e := fmt.Errorf("[S=%d] err= %s", resp.StatusCode, v.(string))
				if r.Alert {
					r.Sentry(e)
				}
				return nil, e
			}
		}
		e := fmt.Errorf("[S=%d] err= %s", resp.StatusCode, string(body))
		if r.Alert {
			r.Sentry(e)
		}
		if span != nil {
			span.SetTag("error", true)
		}
		return nil, e
	}

	return body, nil
}

// Sentry 如果状态码为200，但返回的数据不符合预期，手动报警
func (r *MyRequest) Sentry(e error) {
	var ravenHttp = raven.Http{
		URL:     r.URL,
		Method:  r.Method,
		Query:   "",
		Cookies: stringx.Marshal(r.Cookies),
		Headers: r.Header,
		Data:    r.Data,
	}
	sentry.LogAndSentry(e, &ravenHttp)
}

// HTTPDoByClient 发起http请求并返回结果或error
func (r *MyRequest) HTTPDoByClient(client *http.Client) ([]byte, error) {
	if r.Method == "" {
		r.Method = "GET"
	}

	dataNew, err := r.buildBody()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.Method, r.URL, dataNew)
	if err != nil {
		return nil, err
	}

	r.initHeaders(req)
	r.initCookies(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		var errT map[string]interface{}
		err = json.Unmarshal(body, &errT)
		if err != nil {
			return nil, errors.New(string(body))
		}
		errKey := []string{"error_description", "message", "msg", "stack"}
		for _, key := range errKey {
			if v, ok := errT[key]; ok {
				return nil, errors.New(v.(string))
			}
		}
		return nil, errors.New(string(body))
	}

	return body, nil
}

func (r *MyRequest) WithCtx(ctx context.Context) *MyRequest {
	r.Ctx = ctx
	return r
}

func (r *MyRequest) Get() *MyRequest {
	r.Method = http.MethodGet
	return r
}

func (r *MyRequest) Post() *MyRequest {
	r.Method = http.MethodPost
	return r
}

func (r *MyRequest) Put() *MyRequest {
	r.Method = http.MethodPut
	return r
}

func (r *MyRequest) Delete() *MyRequest {
	r.Method = http.MethodDelete
	return r
}

type MyParam map[string]interface{}

func (m MyParam) toUrlParam() string {
	if len(m) == 0 {
		return ""
	}
	u := "?"
	for key, v := range m {
		param := ""
		switch reflect.TypeOf(v).Kind() {
		case reflect.Int:
			param = strconv.Itoa(v.(int))
		case reflect.String:
			param = v.(string)
		default:
			continue
		}
		p := key + "=" + param + "&"
		u += p
	}
	return u
}
func (r *MyRequest) SetHeader(header map[string]string) *MyRequest {
	r.Header = header
	return r
}

// param只支持int和string
func (r *MyRequest) SetUrl(base, extend string, param MyParam) *MyRequest {
	url := fmt.Sprintf("%s%s", base, extend)
	if param != nil {
		url += param.toUrlParam()
	}
	r.URL = url
	return r
}

func (r *MyRequest) AddHook(hooks ...RequestHook) *MyRequest {
	for i := range hooks {
		index := i
		if hooks[i] == nil {
			r.hooks = append(r.hooks, hooks[index])
		}
	}
	return r
}

func (r *MyRequest) SetTimeout(timeout int) *MyRequest {
	r.Timeout = time.Duration(timeout)
	return r
}

func (r *MyRequest) Body(values interface{}) *MyRequest {
	val := make(map[string]interface{})
	if values != nil {
		bf, _ := json.Marshal(values)
		json.Unmarshal(bf, &val)
	}
	r.Data = val
	return r
}

func (r *MyRequest) Result(res interface{}) error {
	buf, err := r.HTTPDo()
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, res)
}
