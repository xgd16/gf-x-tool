package xError

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
)

var (
	// SystemError 系统错误
	SystemError = &ErrMode{-100, "SYSTEM ERROR!!!"}
	// HttpClientRequestError 客户端发送 HTTP 请求 失败错误
	HttpClientRequestError = &ErrMode{1001, "HTTP CLIENT REQUEST ERROR!!!"}
	// MessageError 触发一个用于传递消息的错误
	MessageError = &ErrMode{1002, "MESSAGE ERROR!!!"}
)

// ErrMode 错误模式
type ErrMode struct {
	code    int
	message string
}

// Error 自定义错误
type Error struct {
	mode    *ErrMode
	Content any
}

func CreateMessageError(text string) *Error {
	return Create(MessageError, text)
}

// CreateSystemError 快速创建系统错误
func CreateSystemError(err error) *Error {
	return Create(SystemError, err)
}

// CreateHttpClientRequestError 创建客户端发送 HTTP 请求 失败错误
func CreateHttpClientRequestError[T string | *gjson.Json](content T) *Error {
	return Create(HttpClientRequestError, content)
}

// IsMode 判断模式
func IsMode(xError *Error, mode *ErrMode) bool {
	return xError != nil && xError.GetMode() == mode
}

// Create 创建错误
func Create(mode *ErrMode, content any) *Error {
	return &Error{
		mode:    mode,
		Content: content,
	}
}

// GetMode 获取错误模式
func (t *Error) GetMode() *ErrMode {
	return t.mode
}

// GetGJson 获取 gjson 类型返回值
func (t *Error) GetGJson() *gjson.Json {
	// 判断传入数据累心是否为 gjson
	if jsonData, ok := t.Content.(*gjson.Json); ok {
		return jsonData
	}
	// 转换为 gjson
	jsonData, err := gjson.DecodeToJson(t.Content)
	if err != nil {
		panic("转换 json 出错" + err.Error())
	}
	return jsonData
}

// Error 转换 错误 为 error 类型
func (t *Error) Error() error {
	// 判断传入数据是否为error如果是 直接返回 error 如果不是 进行定制处理
	if _, ok := t.Content.(error); ok {
		return t.Content.(error)
	}
	// 统一处理方法
	return errors.New(fmt.Sprintf("出现错误: Code: %d Message: %s Content: %s", t.mode.code, t.mode.message, gjson.MustEncodeString(t.Content)))
}
