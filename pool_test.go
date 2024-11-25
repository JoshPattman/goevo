package goevo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestPoolSimple(t *testing.T) {
	numCreations := 0
	// Create a new pool that creates varying sized integer lists, with a capacity to hold old objects of 2
	p := NewPool(2, func(size int) []int {
		numCreations++
		return make([]int, size)
	})

	a5a := p.Get(5)
	p.Get(5)
	p.Put(len(a5a), a5a)
	a5b := p.Get(5)
	p.Put(len(a5b), a5b)

	if numCreations != 2 {
		t.Fatalf("expected 2 creations, got %v", numCreations)
	}

	if p.CurrentSize() != 1 {
		t.Fatalf("expected 1 current size but got %v", p.CurrentSize())
	}

	a3a := p.Get(3)
	p.Put(len(a3a), a3a)

	if numCreations != 3 {
		t.Fatalf("expected 3 creations, got %v", numCreations)
	}

	if p.CurrentSize() != 2 {
		t.Fatalf("expected 2 current size but got %v", p.CurrentSize())
	}

	// At this point we should have one size 5 and one size 3

	p.Get(5)
	p.Get(5)

	if numCreations != 4 {
		t.Fatalf("expected 4 creations, got %v", numCreations)
	}

	if p.CurrentSize() != 1 {
		t.Fatalf("expected 1 current size but got %v", p.CurrentSize())
	}

	p.Get(3)
	p.Get(3)

	if numCreations != 5 {
		t.Fatalf("expected 4 creations, got %v", numCreations)
	}

	if p.CurrentSize() != 0 {
		t.Fatalf("expected 0 current size but got %v", p.CurrentSize())
	}
}

func TestPoolPerformance(t *testing.T) {
	itemOrder := make([]int, 10000)
	for i := range itemOrder {
		// 10 size classes, each ~100k allocation size
		itemOrder[i] = rand.Intn(10) + 100000
	}
	deleteQueue := make(chan []int, len(itemOrder))
	constructor := func(s int) []int {
		return make([]int, s)
	}
	// No pool
	nopoolStart := time.Now()
	for i := range itemOrder {
		deleteQueue <- constructor(itemOrder[i])
	}
	nopoolTime := time.Since(nopoolStart)
	fmt.Println("No Pool", nopoolTime)

	// Cleanup
	for len(deleteQueue) > 0 {
		<-deleteQueue
	}

	// Pool
	pool := NewPool(100, constructor)
	poolStart := time.Now()
	for i := range itemOrder {
		deleteQueue <- pool.Get(itemOrder[i])
		if len(deleteQueue) > 50 {
			dq := <-deleteQueue
			pool.Put(len(dq), dq)
		}
	}
	poolTime := time.Since(poolStart)
	fmt.Println("Pool", poolTime)

	if poolTime > nopoolTime {
		t.Fatalf("Using a pool took longer than without!")
	}
}
