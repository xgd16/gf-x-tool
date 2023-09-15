package xTool

import (
	"bytes"
	"context"
	"fmt"
	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"runtime"
	"strings"
	"time"
)

type GrayLogConfigType struct {
	Host string
	Port int
}

var GrayLogConfig *GrayLogConfigType = nil

type JsonOutputsForLogger struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Content string `json:"content"`
}

// SetGrayLogConfig 配置 GrayLog 基础信息
func SetGrayLogConfig(host string, port int) {
	GrayLogConfig = &GrayLogConfigType{
		Host: host,
		Port: port,
	}
}

// SwitchToGraylog 转换到 Graylog 日志
func SwitchToGraylog(name string) glog.Handler {
	return func(ctx context.Context, in *glog.HandlerInput) {
		// 如果没有配置那么按照正常配置写入
		if GrayLogConfig == nil {
			in.Next(ctx)
			return
		}
		// init server config to graylog
		hook := graylog.NewGraylogHook(fmt.Sprintf("%s:%d", GrayLogConfig.Host, GrayLogConfig.Port), map[string]interface{}{})
		// get logger file path and code line
		file, line := getCallerIgnoringLogMulti(5)
		// get logger content message
		p := gconv.Bytes(gstr.Join(gconv.Strings(in.Values), "\n"))
		// init short message
		short := p
		full := []byte("")
		if i := bytes.IndexRune(p, '\n'); i > 0 {
			short = p[:i]
			full = p
		}
		// writer logger
		err := hook.Writer().WriteMessage(&graylog.Message{
			Version:  "1.1",
			Host:     name,
			Short:    string(short),
			Full:     string(full),
			TimeUnix: float64(time.Now().UnixNano()/1000000) / 1000.,
			Level:    int32(in.Level),
			File:     file,
			Line:     line,
			Extra: g.Map{ // custom logger param
				"levelFormat": in.LevelFormat,
				"stack":       in.Stack,
				"traceId":     in.TraceId,
				"device":      hook.Host,
			},
		})
		if err != nil {
			fmt.Printf("[记录日志出错] %s \n", err.Error())
		}
	}
}

// getCaller 返回函数的文件名和行信息
func getCaller(callDepth int, suffixesToIgnore ...string) (file string, line int) {
	// bump by 1 to ignore the getCaller (this) stackframe
	callDepth++
outer:
	for {
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
			break
		}

		for _, s := range suffixesToIgnore {
			if strings.HasSuffix(file, s) {
				callDepth++
				continue outer
			}
		}
		break
	}
	return
}

func getCallerIgnoringLogMulti(callDepth int) (string, int) {
	// the +1 is to ignore this (getCallerIgnoringLogMulti) frame
	return getCaller(callDepth+1, "logrus/hooks.go", "logrus/entry.go", "logrus/logger.go", "logrus/exported.go", "asm_amd64.s")
}
