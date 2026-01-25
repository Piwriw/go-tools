package common

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type R struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Detail string      `json:"detail"`
}

// 响应状态码常量
const (
	CodeSuccess = 200
	CodeFailed  = 400
)

// WriteErrorResponse 写入错误响应到HTTP响应中
// 设置响应头为JSON格式，写入指定的状态码，并将当前对象编码为JSON格式写入响应体。
// 如果编码过程中发生错误，会记录错误日志。
//
// 参数:
//   - w: http.ResponseWriter 类型，用于写入HTTP响应
//   - statusCode: int 类型，表示要写入的HTTP状态码
func (r *R) WriteErrorResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		slog.Error("write error response failed", slog.Any("err", err))
	}
}

// WriteSuccessResponse 写入成功的HTTP响应
// 设置响应头为JSON格式，状态码为200，并将当前对象编码为JSON格式写入响应体。
// 如果编码过程中发生错误，会记录错误日志。
//
// 参数:
//   - w: http.ResponseWriter 类型，用于写入HTTP响应
func (r *R) WriteSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		slog.Error("write success response failed", slog.Any("err", err))
	}
}

// SetData 设置数据
func (r *R) SetData(data interface{}) *R {
	r.Data = data
	return r
}

// SetCode 设置状态码
func (r *R) SetCode(code int) *R {
	r.Code = code
	return r
}

// SetMsg 设置消息
func (r *R) SetMsg(msg string) *R {
	r.Msg = msg
	return r
}

// SetFailed 设置失败信息
func (r *R) SetFailed(msg string, detail string) *R {
	r.Code = CodeFailed
	r.Msg = msg
	r.Detail = detail
	return r
}

// SetSuccess 设置成功信息
func (r *R) SetSuccess(msg string) *R {
	r.Code = CodeSuccess
	r.Msg = msg
	return r
}

func NewR() *R {
	return &R{}
}
