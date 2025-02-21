package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 定义任务类型
type Task func()

// TimerTask 定义定时器任务
type TimerTask struct {
	time     int64 // 到期时间（毫秒）
	callback Task  // 回调函数
}

// EventLoop 定义事件循环结构体
type EventLoop struct {
	micro []Task        // 微任务队列（动态数组）
	macro []Task        // 宏任务队列（动态数组）
	timer []TimerTask   // 定时器队列（动态数组）
	wg    sync.WaitGroup // 用于等待所有任务完成
}

// NewEventLoop 创建事件循环实例
func NewEventLoop() *EventLoop {
	return &EventLoop{
		micro: make([]Task, 0),
		macro: make([]Task, 0),
		timer: make([]TimerTask, 0),
	}
}

// QueueMicrotask 模拟 JS 的 queueMicrotask
func (e *EventLoop) QueueMicrotask(callback Task) {
	e.micro = append(e.micro, callback)
}

// SetTimeout 模拟 JS 的 setTimeout
func (e *EventLoop) SetTimeout(callback Task, timeout int64) {
	e.timer = append(e.timer, TimerTask{
		time:     time.Now().UnixMilli() + timeout,
		callback: callback,
	})
}

// runMicrotasks 清空微任务队列
func (e *EventLoop) runMicrotasks() {
	for len(e.micro) > 0 {
		task := e.micro[0]           // 取出第一个任务
		e.micro = e.micro[1:]       // 移除已执行的任务
		task()                      // 执行任务
	}
}

// checkTimers 检查定时器，将到期任务放入宏任务队列
func (e *EventLoop) checkTimers() {
	now := time.Now().UnixMilli()
	for i := 0; i < len(e.timer); {
		if e.timer[i].time <= now {
			e.macro = append(e.macro, e.timer[i].callback) // 到期任务放入宏任务队列
			e.timer = append(e.timer[:i], e.timer[i+1:]...) // 移除已触发的任务
		} else {
			i++ // 只有未移除时递增
		}
	}
}

// Run 启动事件循环
func (e *EventLoop) Run() {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		for {
			// 1. 先执行所有微任务
			e.runMicrotasks()

			// 2. 检查 timer，将到期任务放入宏任务队列
			e.checkTimers()

			// 3. 如果有微任务，优先执行
			if len(e.micro) > 0 {
				e.runMicrotasks()
			} else if len(e.macro) > 0 {
				// 4. 执行一个宏任务
				task := e.macro[0]
				e.macro = e.macro[1:] // 移除已执行的任务
				task()
			} else {
				// 5. 无任务时，检查是否有未到期定时器
				if len(e.timer) == 0 && len(e.micro) == 0 && len(e.macro) == 0 {
					fmt.Println("Event Loop 结束，所有任务已执行")
					return // 所有任务完成，退出
				}
				// 计算最近的定时器到期时间，休眠等待
				if len(e.timer) > 0 {
					nearest := e.timer[0].time
					for _, t := range e.timer {
						if t.time < nearest {
							nearest = t.time
						}
					}
					waitTime := nearest - time.Now().UnixMilli()
					if waitTime > 0 {
						time.Sleep(time.Duration(waitTime) * time.Millisecond)
					}
				} else {
					time.Sleep(10 * time.Millisecond) // 默认短暂休眠
				}
			}
		}
	}()
}

func main() {
	// 创建事件循环实例
	loop := NewEventLoop()

	// 测试代码
	fmt.Println("Start")
	t1 := time.Now().UnixMilli()
	loop.SetTimeout(func() {
		fmt.Println("Timeout 1", time.Now().UnixMilli()-t1)
		loop.QueueMicrotask(func() {
			fmt.Println("Microtask from Timeout 1")
		})
	}, 1000)
	t2 := time.Now().UnixMilli()
	loop.SetTimeout(func() {
		fmt.Println("Timeout 2", time.Now().UnixMilli()-t2)
	}, 1000)

	loop.QueueMicrotask(func() {
		fmt.Println("Microtask 1")
	})
	fmt.Println("End")

	// 启动事件循环并等待完成
	loop.Run()
	loop.wg.Wait()
}
