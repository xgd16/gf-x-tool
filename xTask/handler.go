package xTask

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"sync"
	"time"
)

// TaskItem 任务设置
type TaskItem struct {
	taskName string
	runTime  int64
	retry    int
	handler  TaskHandler
}

func (t *TaskItem) run(runnerKey string) {
	err := t.handler()
	// 当出现异常且 重试次数大于 0
	if err != nil && t.retry > 0 {
		t.retry -= 1
		return
	}
	// 执行完成后删除任务
	taskList.Remove(runnerKey)
}

func (t *TaskItem) asyncRun(runnerKey string, wg *sync.WaitGroup, c chan int) {
	defer func() {
		wg.Done()
		<-c
	}()

	t.run(runnerKey)
}

// TaskHandler 执行定义
type TaskHandler func() error

// 任务列表
var taskList = gmap.NewListMap(true)

// SetDelTask 删除任务
func SetDelTask(taskName string) {
	if taskList.Contains(taskName) {
		taskList.Remove(taskName)
	}
}

// SetTask 设置任务
func SetTask(taskName string, handler TaskHandler, runTime int64, retry int) {
	taskList.Set(taskName, &TaskItem{
		taskName: taskName,
		runTime:  runTime,
		retry:    retry,
		handler:  handler,
	})
}

// Run 运行任务管理
func Run(async bool, forTime time.Duration) {
	ctx := gctx.New()
	// 异常推出处理
	defer func() {
		if err := recover(); err != nil {
			g.Log().Error(ctx, "检测到运行任务异常推出", err)
			// 睡眠三秒
			time.Sleep(3 * time.Second)
			// 重启
			Run(async, forTime)
		}
	}()

	var (
		wg      *sync.WaitGroup
		channel chan int
	)

	for {
		// 获取当前时间戳
		nowTime := gtime.Timestamp()
		// 处理异步运行
		if async {
			wg = &sync.WaitGroup{}
			channel = make(chan int, 30)
		}
		// 循环处理
		taskList.IteratorAsc(func(key interface{}, value interface{}) bool {
			v := value.(*TaskItem)
			// 判断是否到运行时间
			if nowTime < v.runTime {
				time.Sleep(500 * time.Millisecond)
				return true
			}
			// 运行
			if async {
				wg.Add(1)
				channel <- 1

				go v.asyncRun(key.(string), wg, channel)
			} else {
				v.run(key.(string))
			}

			return true
		})
		// 如果开启了异步处理 等待处理完成
		if async {
			wg.Wait()
		}
		// 一秒检测 2 次
		time.Sleep(forTime * time.Millisecond)
	}
}
