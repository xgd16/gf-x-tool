//go:build windows

package xTool

import "fmt"

func SafeExit() {
	fmt.Println("windows 下不支持安全退出")
}
