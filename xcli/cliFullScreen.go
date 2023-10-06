package xcli

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"sync"
	"time"
)

//func main() {
//	RunFullScreenUpdate(func(updater *TerminalPrint) error {
//		// 在这里编写你的逻辑，根据需要更新屏幕内容
//		// 调用 updater.Stop() 来结束全屏内容更新
//		// 调用 updater.DrawText() 来绘制文本
//		// 示例：每一帧显示递增的数字
//		i := 0
//		for {
//			text := fmt.Sprintf("Count: %d", i)
//			updater.DrawText(10, 2, text)
//			updater.DrawText(11, 3, text)
//			i++
//			time.Sleep(100 * time.Millisecond)
//			if i > 100 {
//				updater.Stop()
//				break
//			}
//		}
//
//		return nil
//	})
//}

type bufferType struct {
	X    int
	Y    int
	Text string
}

// TerminalPrint 定义全屏更新器
type TerminalPrint struct {
	stopFlag   bool         // 停止标志
	buffer     []bufferType // 更新缓冲区
	display    []bufferType // 显示缓冲区
	bufferLock sync.Mutex
}

// Stop 停止全屏更新
func (u *TerminalPrint) Stop() {
	u.stopFlag = true
}

// DrawText 在屏幕上绘制文本
func (u *TerminalPrint) DrawText(x, y int, text string) {
	u.bufferLock.Lock()
	defer u.bufferLock.Unlock()

	// 将文本添加到更新缓冲区
	if y >= len(u.buffer) {
		for i := len(u.buffer); i <= y; i++ {
			u.buffer = append(u.buffer, *new(bufferType))
		}
	}
	u.buffer[y] = bufferType{
		X:    x,
		Y:    y,
		Text: text,
	}
}

// TerminalPrintView 终端输出显示
// callback 是回调函数，用于控制每一帧的显示内容
// 结束时，通过调用 callback 参数的 Stop() 方法来退出全屏更新
func TerminalPrintView(callback func(updater *TerminalPrint) error) (err error) {
	// 初始化数据
	if err = termbox.Init(); err != nil {
		return
	}
	defer termbox.Close()
	// 创建全屏更新器实例
	updater := &TerminalPrint{}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("发生错误：%v", r)
				termbox.Close() // 关闭 term box
			}
		}()
		err = callback(updater)
	}()

	for !updater.stopFlag {
		if err = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
			return
		}

		updater.bufferLock.Lock()
		updater.display = make([]bufferType, len(updater.buffer))
		copy(updater.display, updater.buffer) // 复制更新缓冲区内容到显示缓冲区
		updater.bufferLock.Unlock()

		for _, item := range updater.display {
			updater.drawText(item.X, item.Y, item.Text) // 修改这一行
		}

		if err = termbox.Flush(); err != nil {
			return
		}
		time.Sleep(16 * time.Millisecond)
	}
	// 刷新终端内容
	if err = termbox.Sync(); err != nil {
		return
	}
	return
}

func (u *TerminalPrint) drawText(x, y int, text string) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}
