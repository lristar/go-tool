package stringx

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// genericity 基于泛型的string

type IStringG interface {
	// 类型底层是string
	~string
}

type IInteger interface {
	constraints.Integer
}

type IFloat interface {
	constraints.Float
}

func ToInteger[T IStringG, V IInteger](src T) (res V, err error) {
	c, _ := regexp.Compile("^(0|[1-9][0-9]*|-[1-9][0-9]*)$")
	if ok := c.MatchString(string(src)); !ok {
		return res, fmt.Errorf("非数字类型")
	}
	switch reflect.TypeOf(res).Kind() {
	case reflect.Int:
		r, err := strconv.ParseInt(string(src), 0, 64)
		return V(r), err
	case reflect.Int8:
		r, err := strconv.ParseInt(string(src), 0, 8)
		return V(r), err
	case reflect.Int16:
		r, err := strconv.ParseInt(string(src), 0, 16)
		return V(r), err
	case reflect.Int32:
		r, err := strconv.ParseInt(string(src), 0, 32)
		return V(r), err
	case reflect.Int64:
		r, err := strconv.ParseInt(string(src), 0, 64)
		return V(r), err
	case reflect.Uint:
		r, err := strconv.ParseUint(string(src), 0, 64)
		return V(r), err
	case reflect.Uint8:
		r, err := strconv.ParseUint(string(src), 0, 8)
		return V(r), err
	case reflect.Uint16:
		r, err := strconv.ParseUint(string(src), 0, 16)
		return V(r), err
	case reflect.Uint32:
		r, err := strconv.ParseUint(string(src), 0, 32)
		return V(r), err
	case reflect.Uint64:
		r, err := strconv.ParseUint(string(src), 0, 64)
		return V(r), err
	default:
		return res, fmt.Errorf("类型未设定")
	}
}

func ToFloat[T IStringG, V IFloat](src T) (res V, err error) {
	c, _ := regexp.Compile("^(0\\.?[0-9]*|[1-9][0-9]*\\.?[0-9]*|-[1-9][0-9]*\\.?[0-9]*)$")
	if ok := c.MatchString(string(src)); !ok {
		return res, fmt.Errorf("非浮点数类型")
	}
	switch reflect.TypeOf(src).Kind() {
	case reflect.Float32:
		r, err := strconv.ParseFloat(string(src), 32)
		return V(r), err
	case reflect.Float64:
		r, err := strconv.ParseFloat(string(src), 64)
		return V(r), err
	default:
		return res, fmt.Errorf("类型未设定")
	}
}

func ToStringArrayComma[T IStringG](src T) []string {
	return strings.Split(string(src), COMMA)
}

func ArrayCommaToStringG[T IStringG](src []string) T {
	return T(strings.Join(src, COMMA))
}

func ToStringArraySemicolon[T IStringG](src T) []string {
	return strings.Split(string(src), SEMICOLON)
}

func ArraySemicolonToStringG[T IStringG](src []string) T {
	return T(strings.Join(src, SEMICOLON))
}

// 用于字符串拼接
type Strings struct {
	buf *bytes.Buffer
}

func GetStrings() *Strings {
	return &Strings{buf: &bytes.Buffer{}}
}

//go test -bench='Strings_Add'
//cpu: Intel(R) Core(TM) i3-8100B CPU @ 3.60GHz
//BenchmarkStrings_Add-4           2520090               494.1 ns/op
//BenchmarkStrings_Add2-4           523032              2249 ns/op

// 第一个优于第二个
func (s *Strings) Add(v ...string) error {
	for i := range v {
		_, err := s.buf.Write([]byte(v[i]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Strings) Add2(v ...string) error {
	_, err := fmt.Fprint(s.buf, v)
	if err != nil {
		return err
	}
	return nil
}

func (s *Strings) Addf(format string, v ...any) error {
	_, err := fmt.Fprintf(s.buf, format, v...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Strings) ToString() string {
	return s.buf.String()
}
