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
	Detail string      `json:"detail,omitempty"`
}

const (
	CodeSuccess = 200
	CodeFailed  = 400
)

// WriteJSON 写入JSON响应
func (r *R) writeJSON(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		slog.Error("write response failed", slog.Any("err", err))
	}
}

// WriteErrorResponse 写入错误响应
func (r *R) WriteErrorResponse(w http.ResponseWriter, statusCode int) {
	r.writeJSON(w, statusCode)
}

// WriteSuccessResponse 写入成功响应
func (r *R) WriteSuccessResponse(w http.ResponseWriter) {
	r.writeJSON(w, http.StatusOK)
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

// SetFailed 设置失败信息，清空Data和Detail
func (r *R) SetFailed(msg string, detail string) *R {
	r.Code = CodeFailed
	r.Msg = msg
	r.Detail = detail
	r.Data = nil
	return r
}

// SetSuccess 设置成功信息，清空Detail
func (r *R) SetSuccess(msg string) *R {
	r.Code = CodeSuccess
	r.Msg = msg
	r.Detail = ""
	return r
}

// NewR 创建新的响应对象
func NewR() *R {
	return &R{}
}

// Success 创建成功响应
func Success(data interface{}) *R {
	return &R{Code: CodeSuccess, Data: data, Msg: "success"}
}

// Failed 创建失败响应
func Failed(msg string, detail string) *R {
	return &R{Code: CodeFailed, Msg: msg, Detail: detail}
}

// Example_R 使用示例
func Example_R() {
	r := Success("data")
	_ = r
	// Output: {"code":200,"data":"data","msg":"success"}
}

// Example_R_SetData 使用示例
func Example_R_SetData() {
	r := NewR().SetData("hello").SetCode(CodeSuccess)
	_ = r
	// Output: {"code":200,"data":"hello","msg":""}
}

// Example_R_SetFailed 使用示例
func Example_R_SetFailed() {
	r := NewR().SetFailed("invalid params", "userId is required")
	_ = r
	// Output: {"code":400,"data":null,"msg":"invalid params","detail":"userId is required"}
}

// Example_Success 使用示例
func Example_Success() {
	r := Success(map[string]string{"name": "test"})
	_ = r
	// Output: {"code":200,"data":{"name":"test"},"msg":"success"}
}

// Example_Failed 使用示例
func Example_Failed() {
	r := Failed("database error", "connection timeout")
	_ = r
	// Output: {"code":400,"msg":"database error","detail":"connection timeout"}
}
