package xlib

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// ErrorIsSystem 是否为系统错误
func ErrorIsSystem(err error) bool {
	return gerror.HasCode(gerror.NewCode(gerror.Code(err)), gcode.CodeNil)
}

// GetErrorCode 获取 error 的 g code
func GetErrorCode(err error) gcode.Code {
	gCodeData := gerror.Code(err)
	if gerror.HasCode(gerror.NewCode(gCodeData), gcode.CodeNil) {
		return gcode.New(gCodeData.Code(), fmt.Sprintf("%s", err), gCodeData.Detail())
	}
	return gCodeData
}
