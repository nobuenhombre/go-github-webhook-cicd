package queue

import (
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	exec := func(data interface{}) error {
		t.Logf("Exec Data: %v", data)

		return nil
	}

	queue := NewQueue(exec)

	queue.Run()

	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		queue.Push(i)
	}

	queue.Stop()

	time.Sleep(1 * time.Second)
}
