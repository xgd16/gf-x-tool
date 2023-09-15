package xVerify

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"
	"github.com/xgd16/gf-x-tool/xTool"
)

type HandlerInterface interface {
	Rules() map[string]string
	Message() map[string]string
}

// DefaultMsgNilResp 空 msg 默认处理方式
var DefaultMsgNilResp = func(r *ghttp.Request) {
	xTool.FastResp(r).SetStatus(1001).Resp("参数错误")
}

// DefaultErrMsgResp msg 默认处理方式
var DefaultErrMsgResp = func(r *ghttp.Request, err gvalid.Error) {
	xTool.FastResp(r).SetStatus(1001).Resp(err.FirstError().Error())
}

func HandlerStruct[T any](r *ghttp.Request, i HandlerInterface) *T {
	Handler(r, i)
	data := new(T)
	err := r.GetStruct(data)
	if err != nil {
		panic("failed to convert request to struct")
	}
	return data
}

func HandlerJson(r *ghttp.Request, i HandlerInterface) *gjson.Json {
	Handler(r, i)
	return gjson.New(r.GetMap(), true)
}

func Handler(r *ghttp.Request, i HandlerInterface) {
	// 获取返回的msg
	message := i.Message()
	// 判断返回的msg如果是nil的话
	msgNil := false
	if message == nil {
		msgNil = true
		message = make(map[string]string, 1)
	}
	// 获取数据
	data := r.GetMap()
	// 调用处理验证
	err := g.Validator().Rules(i.Rules()).Messages(message).Data(data).Run(r.GetCtx())
	if err != nil {
		if msgNil {
			DefaultMsgNilResp(r)
		} else {
			DefaultErrMsgResp(r, err)
		}
	}

}
