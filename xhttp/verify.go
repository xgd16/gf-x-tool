package xhttp

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"
)

type VerifyHandlerInterface interface {
	Rules() map[string]string
	Message() map[string]string
}

// VerifyDefaultMsgNilResp 空 msg 默认处理方式
var VerifyDefaultMsgNilResp = func(r *ghttp.Request) {
	FastResp(r).SetStatus(1001).Resp("参数错误")
}

// VerifyDefaultErrMsgResp msg 默认处理方式
var VerifyDefaultErrMsgResp = func(r *ghttp.Request, err gvalid.Error) {
	FastResp(r).SetStatus(1001).Resp(err.FirstError().Error())
}

func VerifyDataToStruct[T VerifyHandlerInterface](r *ghttp.Request, i T) T {
	VerifyHandler(r, i)
	err := r.GetStruct(&i)
	if err != nil {
		panic("failed to convert request to struct")
	}
	return i
}

func VerifyDataToJson(r *ghttp.Request, i VerifyHandlerInterface) *gjson.Json {
	VerifyHandler(r, i)
	return gjson.New(r.GetMap(), true)
}

func VerifyHandler(r *ghttp.Request, i VerifyHandlerInterface) {
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
			VerifyDefaultMsgNilResp(r)
		} else {
			VerifyDefaultErrMsgResp(r, err)
		}
	}
}
