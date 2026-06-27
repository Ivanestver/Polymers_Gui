package base

import "testing"

func TestQueuePush(t *testing.T) {
	queue := Queue[int]{}
	queue.Push(0)
	if queue[0] != 0 {
		t.Fatalf("Wrong number at position 0. Expected 0, got:  %v", queue[0])
	}
	queue.Push(1)
	if queue[0] != 0 {
		t.Fatalf("Wrong number at position 0. Expected 0, got:  %v", queue[0])
	}
	if queue[1] != 1 {
		t.Fatalf("Wrong number at position 1. Expected 1, got:  %v", queue[1])
	}
}

func TestQueuePopNoError(t *testing.T) {
	queue := Queue[int]{}
	queue.Push(0)
	queue.Push(1)
	if res, ok := queue.Pop(); !ok {
		t.Fatal("Could not pop an existing element")
	} else if resInt, ok := res.(int); !ok {
		t.Fatal("The given value is not int")
	} else {
		if resInt != 0 {
			t.Fatalf("Wrong number. Expected: 0, got: %v", resInt)
		}
	}

	if len(queue) == 0 {
		t.Fatal("Queue must not be empty")
	}
}

func TestQueuePopError(t *testing.T) {
	queue := Queue[int]{}
	if _, ok := queue.Pop(); ok {
		t.Fatal("Empty queue cannot return true")
	}
	queue.Push(0)
	if _, ok := queue.Pop(); !ok {
		t.Fatal("Non-empty queue must return an element")
	}
	if _, ok := queue.Pop(); ok {
		t.Fatal("Empty queue cannot return true")
	}
}

func TestQueueIsEmpty(t *testing.T) {
	queue := Queue[int]{}
	if !queue.IsEmpty() {
		t.Fatalf("Queue must be empty")
	}
	queue.Push(0)
	if queue.IsEmpty() {
		t.Fatalf("Queue must not be empty")
	}
	queue.Pop()
	if !queue.IsEmpty() {
		t.Fatalf("Queue must be empty")
	}
}

func TestQueuePeekNoError(t *testing.T) {
	queue := Queue[int]{}
	queue.Push(0)
	if res, ok := queue.Peek(); !ok {
		t.Fatal("Could not peek an existing element")
	} else if resInt, ok := res.(int); !ok {
		t.Fatal("The given value is not int")
	} else {
		if resInt != 0 {
			t.Fatalf("Wrong number. Expected: 0, got: %v", resInt)
		}
	}

	if queue.IsEmpty() {
		t.Fatal("Queue must not be empty")
	}
}

func TestQueuePeekError(t *testing.T) {
	queue := Queue[int]{}
	if _, ok := queue.Peek(); ok {
		t.Fatal("Empty queue cannot return true")
	}
	queue.Push(0)
	if _, ok := queue.Peek(); !ok {
		t.Fatal("Non-empty queue must return an element")
	}
}

func TestQueuePeekNotSafe(t *testing.T) {
	queue := Queue[int]{}
	queue.Push(0)
	res := queue.PeekNotSafe()
	if resInt, ok := res.(int); !ok {
		t.Fatal("The given value is not int")
	} else {
		if resInt != 0 {
			t.Fatalf("Wrong number. Expected: 0, got: %v", resInt)
		}
	}

	if queue.IsEmpty() {
		t.Fatal("Queue must not be empty")
	}
}
