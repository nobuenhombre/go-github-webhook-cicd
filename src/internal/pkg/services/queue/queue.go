package queue

import (
	"container/list"
	"fmt"
	"github.com/nobuenhombre/suikat/pkg/ge"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

//---------------------------------------------------------------------------
// 1. Push in queue - fast, do not wait Exec func result
// 2. All queue items call Exec func by FIFO - not parallel
// 3. Can stop in every time moment, but if Exec func start - it must be done
//---------------------------------------------------------------------------

type ExecFunc func(data interface{}) error

type Conn struct {
	items          *list.List
	exec           ExecFunc
	active         atomic.Bool
	itemsMutex     *sync.Mutex
	waitRunnerStop *sync.WaitGroup
}

func NewQueue(exec ExecFunc) Service {
	return &Conn{
		items:          list.New(),
		exec:           exec,
		itemsMutex:     new(sync.Mutex),
		waitRunnerStop: new(sync.WaitGroup),
	}
}

func (q *Conn) Activate() {
	// wait current goroutine exit before start new
	q.waitRunnerStop.Wait()

	q.active.CompareAndSwap(false, true)
}

func (q *Conn) DeActivate() {
	q.active.CompareAndSwap(true, false)

	// wait current goroutine exit
	q.waitRunnerStop.Wait()
}

func (q *Conn) Push(item interface{}) error {
	if !q.active.Load() {
		return ge.Pin(&InActiveError{})
	}

	q.itemsMutex.Lock()
	defer q.itemsMutex.Unlock()

	q.items.PushBack(item)

	return nil
}

func (q *Conn) getNextItem() *list.Element {
	q.itemsMutex.Lock()
	defer q.itemsMutex.Unlock()

	var item *list.Element
	if q.items.Len() > 0 {
		item = q.items.Front()
		q.items.Remove(item)
	}

	return item
}

func (q *Conn) Run() error {
	if q.active.Load() {
		return ge.Pin(&AlreadyActiveError{})
	}

	q.Activate()

	go func() {
		q.waitRunnerStop.Add(1)
		for {
			if !q.active.Load() {
				q.waitRunnerStop.Done()
				return
			}

			item := q.getNextItem()
			if item != nil {
				err := q.exec(item.Value)
				if err != nil {
					log.Println(fmt.Sprintf("Run exec error: %v", err))
				}
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	return nil
}

func (q *Conn) Stop() {
	q.DeActivate()
}
