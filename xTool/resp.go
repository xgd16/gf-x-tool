package xTool

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// GetMessageToStatus 根据状态获取 Message
var GetMessageToStatus func(status int) string = nil

// RespAdvancedOption 高级配置
var RespAdvancedOption func(t *FastResponse, m *map[string]any) = nil

// RespErrorMsg 当获取到Error时返回错误信息
var RespErrorMsg = false

// FastResponseOptionType 快速返回设置类型
type FastResponseOptionType struct {
	// 返回 code
	DefaultSuccessCode int `json:"defaultSuccessCode"`
	DefaultErrorCode   int `json:"defaultErrorCode"`
	// http code
	DefaultSuccessStatusCode int `json:"defaultSuccessStatusCode"`
	DefaultErrorStatusCode   int `json:"defaultErrorStatusCode"`
	// 返回字段
	MsgName  string `json:"msgName"`
	CodeName string `json:"codeName"`
	DataName string `json:"dataName"`
	TimeName string `json:"timeName"`
	// 返回 message
	SuccessMsg string `json:"successMsg"`
	ErrorMsg   string `json:"errorMsg"`
}

// FastResponseOption 快速返回设置
var FastResponseOption = &FastResponseOptionType{
	DefaultSuccessCode:       1000,
	DefaultErrorCode:         1001,
	DefaultSuccessStatusCode: 200,
	DefaultErrorStatusCode:   500,
	MsgName:                  "msg",
	CodeName:                 "code",
	DataName:                 "data",
	TimeName:                 "time",
	SuccessMsg:               "SUCCESS !!!",
	ErrorMsg:                 "ERROR !!!",
}

func FastResp(r *ghttp.Request, args ...any) *FastResponse {
	return CreateFastResponse(r, args...)
}

const (
	JsonRespMode uint8 = iota
	MsgRespMode
)

// CreateFastResponse 创建快速返回
func CreateFastResponse(r *ghttp.Request, args ...any) *FastResponse {
	ifMode := true
	isErr := false
	var err error
	// 判断是否进行返回 默认进行返回
	if len(args) >= 1 {
		if args[0] != nil {
			switch args[0].(type) {
			case error:
				err = args[0].(error)
				g.Log().Error(r.GetCtx(), err)
				ifMode = true
				isErr = true
			default:
				ifMode = gconv.Bool(args[0])
			}
		} else {
			ifMode = false
		}
	}
	if !ifMode {
		return &FastResponse{}
	}
	// 创建 基础数据
	f := &FastResponse{
		r:        r,
		Status:   FastResponseOption.DefaultSuccessCode,
		Msg:      "",
		IfMode:   ifMode,
		Err:      err,
		respMode: JsonRespMode,
		Data:     []int{},
	}
	// 设置返回 类型 默认成功
	if len(args) >= 2 {
		if gconv.Bool(args[1]) {
			f.SuccessStatus()
		} else {
			f.ErrorStatus()
		}
	} else {
		if isErr {
			f.ErrorStatus()
		} else {
			f.SuccessStatus()
		}
	}

	return f
}

// FastResponse 快速返回
type FastResponse struct {
	r          *ghttp.Request
	Status     int
	StatusCode int
	Msg        string
	Data       any
	Err        error
	respMode   uint8
	IfMode     bool
	respMap    map[string]any
}

// SetStatus 设置状态
func (t *FastResponse) SetStatus(status int) *FastResponse {
	if !t.IfMode {
		return t
	}
	t.Status = status
	return t
}

// SetStatusCode 设置返回的 HTTP 状态码
func (t *FastResponse) SetStatusCode(statusCode int) *FastResponse {
	if !t.IfMode {
		return t
	}
	t.StatusCode = statusCode
	return t
}

func (t *FastResponse) SetRespMode(mode uint8) *FastResponse {
	if mode > MsgRespMode {
		mode = JsonRespMode
	}
	t.respMode = mode
	return t
}

// ErrorStatus 错误状态
func (t *FastResponse) ErrorStatus() *FastResponse {
	if !t.IfMode {
		return t
	}
	return t.SetStatus(FastResponseOption.DefaultErrorCode).SetStatusCode(FastResponseOption.DefaultErrorStatusCode).SetMsg(FastResponseOption.ErrorMsg)
}

// SuccessStatus 成功状态
func (t *FastResponse) SuccessStatus() *FastResponse {
	if !t.IfMode {
		return t
	}
	return t.SetStatus(FastResponseOption.DefaultSuccessCode).SetStatusCode(FastResponseOption.DefaultSuccessStatusCode).SetMsg(FastResponseOption.SuccessMsg)
}

// TxRollBack 事务推回
func (t *FastResponse) TxRollBack(tx gdb.TX) *FastResponse {
	if t.IfMode {
		_ = tx.Rollback()
	}
	return t
}

// SetMsg 设置返回
func (t *FastResponse) SetMsg(msg string) *FastResponse {
	if !t.IfMode {
		return t
	}
	t.Msg = msg
	return t
}

// SetData 设置返回数据
func (t *FastResponse) SetData(data any) *FastResponse {
	if !t.IfMode {
		return t
	}
	t.Data = data
	return t
}

// Callback 当条件为 true 时触发
func (t *FastResponse) Callback(fn func(t *FastResponse)) *FastResponse {
	if t.IfMode {
		fn(t)
	}
	return t
}

func (t *FastResponse) Resp(args ...any) {
	t.Response(args...)
}

func (t *FastResponse) respHandler() {
	switch t.respMode {
	case JsonRespMode:
		t.r.Response.WriteJsonExit(t.respMap)
		break
	case MsgRespMode:
		t.r.Response.WriteExit(t.Msg)
		break
	}
}

// Response 执行返回操作 (old)
// Msg 参数 1 要返回的文本参数
// Data 参数 2 要返回的数据
// Status 要返回的 code 码
func (t *FastResponse) Response(args ...any) {
	if !t.IfMode {
		return
	}
	t.respMap = make(map[string]any, 1)

	if len(args) >= 1 {
		s := gconv.String(args[0])
		if s != "" {
			t.Msg = s
		}
	}

	if len(args) >= 2 && args[1] != nil {
		t.Data = args[1]
	}

	if len(args) >= 3 {
		t.Status = gvar.New(args[2], true).Int()
	}
	// 将状态数据转换为 对应的消息
	if GetMessageToStatus != nil && t.Status != FastResponseOption.DefaultSuccessCode && t.Status != FastResponseOption.DefaultErrorCode {
		t.Msg = GetMessageToStatus(t.Status)
	}
	// 调用高级配置
	if RespAdvancedOption != nil {
		RespAdvancedOption(t, &t.respMap)
	}
	// 返回错误信息
	if RespErrorMsg && t.Err != nil {
		t.Msg = fmt.Sprintf("%s", t.Err)
	}
	if gvar.New(t.Data).IsEmpty() {
		t.Data = []int{}
	}
	// 制造返回数据
	t.respMap[FastResponseOption.CodeName] = t.Status
	t.respMap[FastResponseOption.MsgName] = t.Msg
	t.respMap[FastResponseOption.DataName] = t.Data
	t.respMap[FastResponseOption.TimeName] = gtime.Now().UnixMilli()
	// 写入状态
	t.r.Response.Status = t.StatusCode
	// 触发返回
	t.respHandler()
}

// FastRespJsonConfigOn 快速返回 json 配置开启
func FastRespJsonConfigOn() (err error) {
	cfg, err := g.Cfg().Get(gctx.New(), "fastResp")
	if err != nil {
		return
	}
	if !cfg.IsEmpty() {
		err = cfg.Struct(FastResponseOption)
		if err != nil {
			return
		}
	}
	return
}
