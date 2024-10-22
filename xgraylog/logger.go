package xgraylog

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type GrayLogConfigType struct {
	Host string
	Port int
}

var GrayLogConfig map[string]*GrayLogConfigType = nil
var GrayLogConfigArr []*GrayLogConfigType = nil
var nextKey = 0

type JsonOutputsForLogger struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Content string `json:"content"`
}

// SetGrayLogConfig 配置 GrayLog 基础信息
func SetGrayLogConfig(config *gvar.Var) {
	GrayLogConfig = make(map[string]*GrayLogConfigType)
	mapConfig := config.MapStrVar()
	if !mapConfig["host"].IsEmpty() && !mapConfig["port"].IsEmpty() {
		newConfig := &GrayLogConfigType{
			Host: mapConfig["host"].String(),
			Port: mapConfig["port"].Int(),
		}
		GrayLogConfig["default"] = newConfig
		GrayLogConfigArr = append(GrayLogConfigArr, newConfig)
		return
	}
	for k, v := range mapConfig {
		item := v.MapStrVar()
		if item["host"].IsEmpty() || item["port"].IsEmpty() {
			continue
		}
		newConfig := &GrayLogConfigType{
			Host: item["host"].String(),
			Port: item["port"].Int(),
		}
		GrayLogConfig[k] = newConfig
		GrayLogConfigArr = append(GrayLogConfigArr, newConfig)
	}
}

func GetConfig() (conf *GrayLogConfigType, err error) {
	if len(GrayLogConfig) == 0 || GrayLogConfig == nil {
		err = fmt.Errorf("GrayLogConfig is nil")
		return
	}
	if len(GrayLogConfigArr) == 1 {
		conf = GrayLogConfigArr[0]
		return
	}
	configLen := len(GrayLogConfigArr) - 1
	conf = GrayLogConfigArr[nextKey]
	nextKey += 1
	if nextKey >= configLen {
		nextKey = 0
	}
	return
}

// SwitchToGraylog 转换到 Graylog 日志
func SwitchToGraylog(name string, optFunc func(ctx context.Context, m g.Map)) glog.Handler {
	return func(ctx context.Context, in *glog.HandlerInput) {
		// 如果没有配置那么按照正常配置写入
		if GrayLogConfig == nil {
			in.Next(ctx)
			return
		}
		config, err := GetConfig()
		if err != nil {
			fmt.Println("[记录日志出错]", err.Error())
		}
		// init server config to graylog
		hook := graylog.NewGraylogHook(fmt.Sprintf("%s:%d", config.Host, config.Port), map[string]interface{}{})
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
		extra := g.Map{ // custom logger param
			"levelFormat": in.LevelFormat,
			"stack":       in.Stack,
			"traceId":     in.TraceId,
			"device":      hook.Host,
		}
		optFunc(ctx, extra)
		// writer logger
		err = hook.Writer().WriteMessage(&graylog.Message{
			Version:  "1.1",
			Host:     name,
			Short:    string(short),
			Full:     string(full),
			TimeUnix: float64(time.Now().UnixNano()/1000000) / 1000.,
			Level:    int32(in.Level),
			File:     file,
			Line:     line,
			Extra:    extra,
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
