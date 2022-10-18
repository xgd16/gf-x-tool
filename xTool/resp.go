package xTool

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/net/ghttp"
)

// CreateFastResponse 创建快速返回
// @param r *ghttp.Request
func CreateFastResponse(r *ghttp.Request, args ...any) *FastResponse {
	ifMode := true

	if len(args) > 0 {
		ifMode = gvar.New(args[0], true).Bool()
	}

	return &FastResponse{
		r:      r,
		Status: 1,
		Msg:    "",
		IfMode: ifMode,
		Data:   []int{},
	}
}

// FastResponse 快速返回
type FastResponse struct {
	r          *ghttp.Request
	Status     int
	StatusCode int
	Msg        string
	Data       any
	IfMode     bool
}

// SetStatus 设置状态
func (t *FastResponse) SetStatus(status int) *FastResponse {
	t.Status = status

	return t
}

// SetStatusCode 设置返回的 HTTP 状态码
func (t *FastResponse) SetStatusCode(statusCode int) *FastResponse {
	t.StatusCode = statusCode

	return t
}

// ErrorStatus 错误状态
func (t *FastResponse) ErrorStatus() *FastResponse {
	return t.SetStatus(1001).SetStatusCode(500).SetMsg("ERROR !!!")
}

// SuccessStatus 成功状态
func (t *FastResponse) SuccessStatus() *FastResponse {
	return t.SetStatus(1000).SetStatusCode(200).SetMsg("SUCCESS !!!")
}

// SetMsg 设置返回
func (t *FastResponse) SetMsg(msg string) *FastResponse {
	t.Msg = msg

	return t
}

// SetData 设置返回数据
func (t *FastResponse) SetData(data any) *FastResponse {
	t.Data = data

	return t
}

// Callback 当条件为 true 时触发
func (t *FastResponse) Callback(fn func(t *FastResponse)) *FastResponse {
	if t.IfMode {
		fn(t)
	}

	return t
}

func (t *FastResponse) Response(args ...any) {
	resp := make(map[string]any, 1)

	if len(args) >= 1 {
		t.Msg = gvar.New(args[0], true).String()
	}

	if len(args) >= 2 {
		t.Data = args[1]
	}

	if len(args) >= 3 {
		t.Status = gvar.New(args[2], true).Int()
	}

	resp["code"] = t.Status
	resp["msg"] = t.Msg
	resp["data"] = t.Data

	if t.IfMode {
		t.r.Response.Status = t.StatusCode
		t.r.Response.WriteJsonExit(resp)
	}
}
