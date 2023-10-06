package xhttp

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/grand"
	"sync"
	"testing"
	"time"
)

type AVerify struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Nickname string `json:"nickname"`
}

func (A AVerify) Rules() map[string]string {
	return g.MapStrStr{
		"name":     "required",
		"age":      "required",
		"nickname": "required",
	}
}

func (A AVerify) Message() map[string]string {
	return nil
}

func TestHttpVerify(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		wg := new(sync.WaitGroup)
		host := "127.0.0.1"
		port := grand.N(10000, 15000)
		addr := fmt.Sprintf("%s:%d", host, port)
		go func() {
			wg.Add(1)
			s := g.Server()
			s.SetDumpRouterMap(false)
			s.SetAddr(addr)
			s.BindHandler("/testVerify", func(r *ghttp.Request) {
				go func() {
					time.Sleep(100 * time.Millisecond)
					wg.Done()
				}()
				g.DumpWithType(VerifyDataToStruct(r, new(AVerify)))
				FastResp(r).Resp()
			})
			s.Run()
		}()
		time.Sleep(500 * time.Millisecond)
		s := g.Client().Discovery(nil).ContentJson().PostVar(ctx, fmt.Sprintf("http://%s/testVerify", addr), g.Map{
			"name":     "kaka",
			"age":      15,
			"nickname": "a",
		})
		t.Assert(s.MapStrVar()["code"].Int(), 1000)
		wg.Wait()
	})
}
