package queue

import (
	"errors"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueueRun(t *testing.T) {
	exec := func(data interface{}) error {
		return nil
	}

	queue := NewQueue(exec)

	errFirst := queue.Run()
	errSecond := queue.Run()

	if errFirst != nil {
		t.Fatalf("First Run must return nil error, but returned %v", errFirst)
	}

	if !errors.Is(errSecond, &AlreadyActiveError{}) {
		t.Fatalf("Second Run must return AlreadyActiveError error, but returned %v", errSecond)
	}
}

func TestQueuePush(t *testing.T) {
	exec := func(data interface{}) error {
		return nil
	}

	queue := NewQueue(exec)

	errInactive := queue.Push(1)
	if !errors.Is(errInactive, &InActiveError{}) {
		t.Fatalf("Inactive Queue Push must return InActiveError error, but returned %v", errInactive)
	}
}

func TestQueue(t *testing.T) {
	var executed atomic.Uint32

	exec := func(data interface{}) error {
		log.Printf("Exec Start Data: %v", data)
		time.Sleep(10 * time.Millisecond)
		log.Printf("     Finish Data: %v", data)
		executed.Add(1)
		return nil
	}

	queue := NewQueue(exec)

	err := queue.Run()
	if err != nil {
		t.Fatalf("First Run must return nil error, but returned %v", err)
	}

	for i := 1; i <= 5; i++ {
		err = queue.Push(i)
		if err != nil {
			t.Logf("Push error %v", err)
		}
	}

	time.Sleep(29 * time.Millisecond) // exec func Must Executed 3 times

	queue.Stop()

	if executed.Load() != 3 {
		t.Fatalf("exec func Must Executed 3 times, but executed %v", executed.Load())
	}

	executed.Store(0)

	for i := 6; i <= 10; i++ {
		err = queue.Push(i)
		if !errors.Is(err, &InActiveError{}) {
			t.Fatalf("Inactive Queue Push must return InActiveError error, but returned %v", err)
		}
	}

	err = queue.Run()
	if err != nil {
		t.Fatalf("First Run must return nil error, but returned %v", err)
	}

	for i := 11; i <= 15; i++ {
		err = queue.Push(i)
		if err != nil {
			t.Logf("Push error %v", err)
		}
	}

	time.Sleep(49 * time.Millisecond) // exec func Must Executed 5 times

	queue.Stop()

	if executed.Load() != 5 {
		t.Fatalf("exec func Must Executed 5 times, but executed %v", executed.Load())
	}
}
