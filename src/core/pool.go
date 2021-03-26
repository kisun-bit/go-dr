/*@Title 任务池实现
  @Remark
      要实现该功能需要解决下述的问题：
      1. goroutine如何重用？  回答 -> 池化，创建一个goroutine池
      2. 限制goroutine个数？  回答 -> 同上.
      3. 任务执行流程？       回答 -> "生产者 --(生产任务)--> 队列 --(消费任务)--> 消费者"
      4. 单个goroutine执行失败会导致整个进程崩溃？ 回答 -> 为每一个错误预留一个“恢复处理”的接口
      实现：
      “Talk is cheap. Show me the code.”
  @Description
      pass
*/
package core

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrInvalidJPoolCap   = errors.New("jpkt-pool-error: ErrInvalidJPoolCap")
	ErrPoolAlreadyClosed = errors.New("jpkt-pool-error: ErrPoolAlreadyClosed")
)

// 定义"任务"
type JTask struct {
	Handler func(v ...interface{}) error // 需要执行的函数
	Params  []interface{}                // 函数参数集合
}

func (jt *JTask) Run() {
	_ = jt.Handler(jt.Params)
}

func NewJTask(taskHandler func(v ...interface{}) error, handlerParams []interface{}) *JTask {
	return &JTask{
		Handler: taskHandler,
		Params:  handlerParams,
	}
}

// 定义"任务池"
type JPool struct {
	capacity       uint64                 // 容量
	runningWorkers uint64                 // 正在执行的“goroutine”个数
	isFinish       bool                   // 是否结束
	jTasks         chan *JTask            // 待处理的“任务队列”
	panicHandler   func(v ...interface{}) // 预留异常的处理接口

	sync.Mutex
	internalWg *sync.WaitGroup
}

func NewJPool(capLen uint64) (jp *JPool, err error) {
	if capLen < 0 {
		return nil, ErrInvalidJPoolCap
	}
	return &JPool{
		capacity:   capLen,
		isFinish:   false,
		jTasks:     make(chan *JTask, capLen),
		internalWg: new(sync.WaitGroup),
	}, nil
}

func (jp *JPool) Start() {
	jp.incrRunningWorkers()

	go func() {
		defer func() {
			jp.decrRunningWorkers()
			if r := recover(); r != nil {
				jse := ConvertPanic2StandardErr(r)
				jp.internalWg.Done()
				if jp.panicHandler != nil {
					jp.panicHandler(jse) // TODO 后续支持定制错误处理方法 ...
				} else {
					fmt.Println(jse.ErrorDetail())
				}
				jp.checkWorker() // 防止通道里put task后出现异常，任务通道新增的任务没有goroutine来消费
			}
		}()

		for {
			select {
			case _t, _ok := <-jp.jTasks:
				if _ok {
					_t.Run()
					jp.internalWg.Done()
				} else {
					return
				}
			}
		}
	}()
}

func (jp *JPool) incrRunningWorkers() {
	atomic.AddUint64(&jp.runningWorkers, 1)
}

func (jp *JPool) decrRunningWorkers() {
	atomic.AddUint64(&jp.runningWorkers, ^uint64(0))
}

func (jp *JPool) setFinish(isFinish bool) (isSuccess bool) {
	jp.Lock()
	defer jp.Unlock()
	if jp.isFinish == isFinish {
		return false
	}
	jp.isFinish = isFinish
	return true
}

func (jp *JPool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&jp.runningWorkers)
}

func (jp *JPool) GetCap() uint64 {
	return jp.capacity
}

func (jp *JPool) Put(task *JTask) error {
	jp.Lock()
	defer jp.Unlock()

	if jp.isFinish == true {
		return ErrPoolAlreadyClosed
	}

	if jp.GetRunningWorkers() < jp.GetCap() {
		jp.Start() // add a goroutine to handle task
	}

	if jp.isFinish == false {
		jp.internalWg.Add(1)
		jp.jTasks <- task
	}
	return nil
}

func (jp *JPool) Close() {
	if !jp.setFinish(true) {
		return
	}

	for len(jp.jTasks) > 10 {
		time.Sleep(time.Millisecond * 100)
	}

	jp.internalWg.Wait() // 保证所有任务均已执行，不关心错误

	// close ch
	_closeTaskChan := func() {
		jp.Lock()
		defer jp.Unlock()
		close(jp.jTasks)
	}
	_closeTaskChan()
}

func (jp *JPool) SetPanicHandler(panicHandler func(v ...interface{})) {
	jp.Lock()
	defer jp.Unlock()
	jp.panicHandler = panicHandler
}

func (jp *JPool) checkWorker() {
	jp.Lock()
	defer jp.Unlock()

	// 当前没有 worker 且有任务存在，运行一个 worker 消费任务
	// 没有任务无需考虑 (当前 Put 不会阻塞，下次 Put 会启动 worker)
	if jp.runningWorkers == 0 && len(jp.jTasks) > 0 {
		jp.Start()
	}
}
