package xhttp

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
)

// SetCtxVar 向 ctx 中设置数据
func SetCtxVar(ctx context.Context, key, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetCtxVar 从 ctx 中获取数据
func GetCtxVar(ctx context.Context, key any, def ...any) *gvar.Var {
	value := ctx.Value(key)
	if value == nil && len(def) > 0 {
		value = def[0]
	}
	return gvar.New(value)
}
