package test

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/xgd16/gf-x-tool/xTool"
	"testing"
	"time"
)

// TestFastRespMsgRespMode 快速返回测试
func TestFastRespMsgRespMode(t *testing.T) {
	go gtest.C(t, func(t *gtest.T) {
		s := g.Server()
		s.SetAddr("127.0.0.1:19321")
		s.BindHandler("/test", func(r *ghttp.Request) {
			xTool.FastResp(r).SetRespMode(xTool.MsgRespMode).Resp("testMsg")
		})
		s.BindHandler("/testJson", func(r *ghttp.Request) {
			xTool.FastResp(r).Resp("testMsg")
		})
		s.Run()
	})
	time.Sleep(1 * time.Second)
	// test msg mode
	get, err := g.Client().Get(gctx.New(), "http://127.0.0.1:19321/test")
	gtest.AssertNil(err)
	gtest.Assert(get.ReadAllString(), "testMsg")
	// test json mode
	get, err = g.Client().Get(gctx.New(), "http://127.0.0.1:19321/testJson")
	gtest.AssertNil(err)
	jsonData, err := gjson.DecodeToJson(get.ReadAllString())
	gtest.AssertNil(err)
	gtest.Assert(jsonData.Get("code").Int(), 1000)
	gtest.Assert(jsonData.Get("msg").String(), "testMsg")
}
