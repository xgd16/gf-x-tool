package test

import (
	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	log "github.com/sirupsen/logrus"
	"github.com/xgd16/gf-x-tool/xTool"
	"testing"
)

func TestSendLog(t *testing.T) {
	xTool.SetGrayLogConfig(
		"127.0.0.1",
		9802,
	)
	glog.SetDefaultHandler(xTool.SwitchToGraylog("test"))

	g.Log().Error(gctx.New(), "test123123123123123123123123", "12312312312312aaaaa")
}

func TestOSendLog(t *testing.T) {
	hook := graylog.NewGraylogHook("127.0.0.1:9802", map[string]interface{}{"this": "is logged every time"})
	log.AddHook(hook)
	log.Info("some logging message")
}
