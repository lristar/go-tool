package sentry

import (
	logger "github.com/lristar/go-tool/server/logger"
	"io"
)

var cli *raven.Client

type Client struct {
	url string
}

func (c *Client) Start() (closer io.Closer, err error) {
	cli, err = raven.New(c.url)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c *Client) Close() error {
	cli.Close()
	return nil
}

var tags = map[string]string{"server_name": "platform-backend"}

// Panic 执行函数并捕获+报告错误
func Panic(f func()) {
	raven.CapturePanic(f, tags)
}

// Sentry 报告错误
func Sentry(err error, args ...raven.Interface) {
	raven.CaptureError(err, tags, args...)
}

// LogAndSentry 报告并打印错误
func LogAndSentry(e error, args ...raven.Interface) {
	logger.Errorf("%s\n", e.Error())
	Sentry(e, args...)
}
