package xlib

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/frame/gins"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

// Sha256 获取数据的 sha256
func Sha256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// InArr 是否在数组中
func InArr[T any](a T, arr []T) bool {
	anyArr := make([]any, 0)
	for _, t := range arr {
		anyArr = append(anyArr, t)
	}
	return garray.NewArrayFrom(anyArr).Contains(a)
}

// Maintain 维持函数
func Maintain(handler func()) {
	// 创建一个传递信号的 channel
	sigChan := make(chan os.Signal, 1)
	// 监听退出信号
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// 等待信号
	<-sigChan
	// 退出前调用
	if handler != nil {
		handler()
	}
	// 关闭主线程
	os.Exit(0)
}

// HttpNotJsonHandler 处理 http 返回的不是个 json
var HttpNotJsonHandler = func(str string) string {
	if !gjson.Valid(str) {
		return gjson.MustEncodeString(g.Map{"resp": str})
	}
	return str
}

// HttpClient 发起一个http请求
func HttpClient(method, url string, data any, opt ...any) (gJson *gjson.Json, err error) {
	gHttp := g.Client()
	gHttp.SetDiscovery(nil)
	request, err := gHttp.SetContentType("application/json").SetTimeout(60*time.Second).DoRequest(gctx.New(), method, url, data)
	if err != nil {
		return
	}
	defer request.Close()
	text := request.ReadAllString()
	if len(opt) >= 1 {
		text = opt[0].(func(str string) string)(text)
	}
	gJson, err = gjson.DecodeToJson(text)
	return
}

func IF(condition bool, a, b any) any {
	if condition {
		return a
	} else {
		return b
	}
}

// SelectAllToStruct 查询多条转换为结构体组
func SelectAllToStruct[T any](all gdb.Result) []T {
	d := new([]T)
	if err := all.Structs(d); err != nil {
		panic("请检查 转换数据结构")
	}
	return *d
}

// SelectOneToStruct 查询单条转换为结构体
func SelectOneToStruct[T any](one gdb.Record) *T {
	d := new(T)
	if err := one.Struct(d); err != nil {
		panic("请检查 转换数据结构")
	}
	return d
}

// RedisScanData 扫描 redis 数据 (替代 keys 操作)
func RedisScanData(ctx context.Context, key string, fn func(keys []string) (err error)) (err error) {
	conn, err := gins.Redis().Conn(ctx)
	if err != nil {
		return
	}
	defer func() { _ = conn.Close(ctx) }()
	// scan 出 key 中所有的数据 -1 代表还没有执行一次操作
	index := "0"
	for {
		items, err := conn.Do(ctx, "SCAN", index, "MATCH", key, "COUNT", 10)
		if err != nil {
			return err
		}
		if items.IsEmpty() {
			break
		}
		scanData := items.Vars()
		if len(scanData) < 2 {
			return errors.New("scan data error")
		}
		index = scanData[0].String()
		arr := scanData[1].Strings()
		if len(arr) > 0 {
			if err = fn(arr); err != nil {
				return err
			}
		}
	}
	return
}

// Base642File Base64 转文件
func Base642File(base64Str, filePath, filename string) (err error) {
	suffix, err := GetFileExtension(base64Str)
	if err != nil {
		return
	}
	base64Data, err := GetBase64Data(base64Str)
	if err != nil {
		return
	}
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return
	}
	err = gfile.PutBytes(fmt.Sprintf("%s/%s.%s", filePath, filename, suffix), data)
	return
}

// GetFileExtension 获取文件后缀
func GetFileExtension(base64Str string) (string, error) {
	prefix := "data:image/"
	suffixIndex := strings.Index(base64Str, ";base64,")
	if suffixIndex == -1 {
		return "", fmt.Errorf("无效的 Base64 字符串格式")
	}
	return base64Str[len(prefix):suffixIndex], nil
}

// GetBase64Data 获取 Base64 数据
func GetBase64Data(base64Str string) (string, error) {
	dataIndex := strings.Index(base64Str, ";base64,")
	if dataIndex == -1 {
		return "", fmt.Errorf("无效的 Base64 字符串格式")
	}
	return base64Str[dataIndex+len(";base64,"):], nil
}
