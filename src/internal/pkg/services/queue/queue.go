package queue

import (
	"context"
	"fmt"
	"log"
)

type ExecFunc func(data interface{}) error

type Conn struct {
	items  chan interface{}
	ctx    context.Context
	cancel context.CancelFunc
	exec   ExecFunc
}

func NewQueue(exec ExecFunc) Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Conn{
		items:  make(chan interface{}, 5),
		ctx:    ctx,
		cancel: cancel,
		exec:   exec,
	}
}

func (q *Conn) Push(item interface{}) {
	q.items <- item
}

func (q *Conn) Run() {
	runner := func() {
		for {
			select {
			case item := <-q.items:
				err := q.exec(item)
				if err != nil {
					log.Println(fmt.Sprintf("Exec error: %v", err))
				}
			case <-q.ctx.Done():
				close(q.items)
				log.Println("Queue closed")
				return
			default:
				// no one option worked
			}
		}
	}

	go runner()
}

func (q *Conn) Stop() {
	q.cancel()
}
