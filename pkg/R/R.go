package common

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
