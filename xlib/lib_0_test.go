package xlib

import (
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func TestInArr(t *testing.T) {
	type AType struct {
		A int
		B string
	}

	gtest.C(t, func(t *gtest.T) {
		// 基础类型验证
		arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		t.Assert(InArr(3, arr), true)
		t.Assert(InArr(100, arr), false)
		// 结构体验证
		sArr := []AType{{A: 1, B: "A"}, {A: 2, B: "B"}, {A: 3, B: "C"}, {A: 4, B: "D"}, {A: 5, B: "E"}}
		t.Assert(InArr(AType{A: 5, B: "E"}, sArr), true)
		t.Assert(InArr(AType{A: 1, B: "E"}, sArr), false)
	})
}
