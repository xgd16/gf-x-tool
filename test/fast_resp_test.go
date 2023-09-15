package test

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/xgd16/gf-x-tool/xTool"
	"testing"
	"time"
)

func TestFastRespFileCfg(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		go func() {
			_ = xTool.FastRespJsonConfigOn()
			s := g.Server()
			s.BindHandler("/test", func(r *ghttp.Request) {
				xTool.FastResp(r).Resp()
			})
			s.SetAddr("127.0.0.1:9441")
			s.Run()
		}()
		time.Sleep(1 * time.Second)
		get, _ := g.Client().Get(gctx.New(), "http://127.0.0.1:9441/test")
		fmt.Println(get.ReadAllString(), 1)
	})
}
