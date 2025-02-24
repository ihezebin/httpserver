package httpserver

// Body 统一的响应体结构，code 状态码，成功为 0，message 为错误消息, data 为响应的数据。
// Unified response structure, code status code, success is 0, message is error message, data is the response data
type Body[T any] struct {
	status  int
	Code    Code   `json:"code"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func (b *Body[T]) WithError(err error) *Body[T] {
	b.Code = CodeInternalServerError
	b.Message = err.Error()
	return b
}

func (b *Body[T]) WithErr(err *Err) *Body[T] {
	b.Code = err.Code
	b.Message = err.Error()
	if err.Status != 0 {
		b.status = err.Status
	}
	return b
}

type EmptyType struct{}

var EmptyRequest EmptyType = struct{}{}
var EmptyResponse EmptyType = struct{}{}
