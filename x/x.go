package x

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xgd16/gf-x-tool/xcli"
	"github.com/xgd16/gf-x-tool/xhttp"
	"github.com/xgd16/gf-x-tool/xstorage"
)

// Resp HTTP 快速定义返回结果
func Resp(r *ghttp.Request, args ...any) {
	xhttp.CreateFastResponse(r, args...)
}

// XDB 创建简易文件存储
func XDB() *xstorage.XDB {
	return xstorage.CreateXDB()
}

// TerminalPrintView 终端打印显示
func TerminalPrintView(callback func(updater *xcli.TerminalPrint) error) (err error) {
	err = xcli.TerminalPrintView(callback)
	return
}

// VerifyDataToStruct 验证数据后转换为结构体
func VerifyDataToStruct[T xhttp.VerifyHandlerInterface](r *ghttp.Request, i T) T {
	return xhttp.VerifyDataToStruct(r, i)
}

// VerifyDataToJson 验证数据后转换为 gjson
func VerifyDataToJson(r *ghttp.Request, i xhttp.VerifyHandlerInterface) *gjson.Json {
	return xhttp.VerifyDataToJson(r, i)
}
