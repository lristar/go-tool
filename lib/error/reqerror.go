package myerror

import "fmt"

type MyError interface {
	Error() string
	Code() int
}

// ReqError 请求类型错误
type ReqError struct {
	msg  string
	code int
}

func (r ReqError) Error() string {
	return r.msg
}

// New 新建error
func New(msg string, code int) ReqError {
	return ReqError{msg: msg, code: code}
}

// Errorf 输出msg并新建error
func Errorf(format string, code int, a ...interface{}) ReqError {
	return ReqError{msg: fmt.Sprintf(format, a...), code: code}
}

// Code 返回code
func (r ReqError) Code() int {
	return r.code
}
