//go:build windows

package xlib

import "fmt"

func SafeExit() {
	fmt.Println("windows 下不支持安全退出")
}
