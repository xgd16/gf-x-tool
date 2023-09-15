package xTool

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"sync"
	"time"
)

//func main() {
//	RunFullScreenUpdate(func(updater *FullScreenUpdater) error {
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

// FullScreenUpdater 定义全屏更新器
type FullScreenUpdater struct {
	stopFlag   bool         // 停止标志
	buffer     []bufferType // 更新缓冲区
	display    []bufferType // 显示缓冲区
	bufferLock sync.Mutex
}

// Stop 停止全屏更新
func (u *FullScreenUpdater) Stop() {
	u.stopFlag = true
}

// DrawText 在屏幕上绘制文本
func (u *FullScreenUpdater) DrawText(x, y int, text string) {
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

// RunFullScreenUpdate 开启全屏内容更新
// callback 是回调函数，用于控制每一帧的显示内容
// 结束时，通过调用 callback 参数的 Stop() 方法来退出全屏更新
func RunFullScreenUpdate(callback func(updater *FullScreenUpdater) error) {
	err := termbox.Init()
	if err != nil {
		panic(err) // 初始化错误，使用 panic 中止程序
	}
	defer termbox.Close()

	updater := &FullScreenUpdater{} // 创建全屏更新器实例

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("发生错误：%v", r)
				termbox.Close() // 关闭 termbox
				panic(err)      // 重新抛出错误
			}
		}()
		err := callback(updater)
		if err != nil {
			panic(err) // 使用 panic 中止程序
		}
	}()

	for !updater.stopFlag {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		updater.bufferLock.Lock()
		updater.display = make([]bufferType, len(updater.buffer))
		copy(updater.display, updater.buffer) // 复制更新缓冲区内容到显示缓冲区
		updater.bufferLock.Unlock()

		for _, item := range updater.display {
			updater.drawText(item.X, item.Y, item.Text) // 修改这一行
		}

		termbox.Flush()
		time.Sleep(16 * time.Millisecond)
	}

	termbox.Sync() // 刷新终端内容
}

func (u *FullScreenUpdater) drawText(x, y int, text string) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}
