package myerror

import (
	"errors"
	"strings"
)

var (
	ErrObjectID   = errors.New("ObjectID invalid, ObjectID错误")
	ErrTableName  = errors.New("table name invalid, 表名错误")
	ErrType       = errors.New("type invalid, 类型错误")
	ErrParam      = errors.New("type invalid, 参数错误")
	ErrAbnormal   = errors.New("Abnormal Error Occur, 出现非正常错误")
	ErrNotFound   = errors.New("Record Not Found, 记录不存在")
	ErrXlsxFormat = errors.New("Xlsx file format invalid, XLSX文件不符合导入格式")
)

// IsDuplicateKeyError 唯一键重复错误
func IsDuplicateKeyError(e error) bool {
	return strings.Contains(e.Error(), "E11000 duplicate key error")
}

// IsNotFoundError 无结果错误
func IsNotFoundError(e error) bool {
	return strings.Contains(e.Error(), "not found")
}

// IsDuplicateViewError 重复表修改
func IsDuplicateViewError(e error) bool {
	return strings.Contains(e.Error(), "already exists")
}
