// ProducerConsumer.go
// Copyright (C) 2024 Dr. Joseph Kehoe

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

//--------------------------------------------
// Author: Joseph Kehoe (Joseph.Kehoe@setu.ie)
// Modified by: Shuai Jiang
// Description: Finite Buffer Producer-Consumer using Mutex and Semaphores
//--------------------------------------------

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type Event struct {
	ID int
}

// SafeBuffer implements the thread-safe queue.
type SafeBuffer struct {
	queue     []Event             // The actual buffer
	capacity  int64               // Max size
	mutex     sync.Mutex          // Protects access to the queue
	semItems  *semaphore.Weighted // Counts items available (initially 0)
	semSpaces *semaphore.Weighted // Counts spaces available (initially capacity)
}

// NewSafeBuffer initializes the buffer and semaphores
func NewSafeBuffer(size int64) *SafeBuffer {
	return &SafeBuffer{
		queue:     make([]Event, 0, size),
		capacity:  size,
		semItems:  semaphore.NewWeighted(size),
		semSpaces: semaphore.NewWeighted(size),
	}
}

// Put adds an item
func (b *SafeBuffer) Put(e Event) {
	ctx := context.TODO()

	// 1. Wait for Space
	// We acquire 1 token from semSpaces. If buffer is full, this blocks.
	if err := b.semSpaces.Acquire(ctx, 1); err != nil {
		fmt.Printf("Failed to acquire space: %v\n", err)
		return
	}

	// 2. Lock Mutex
	b.mutex.Lock()

	// 3. Add to Buffer
	b.queue = append(b.queue, e)
	fmt.Printf("Producer added Event %d. (Buffer Size: %d)\n", e.ID, len(b.queue))

	// 4. Unlock Mutex
	b.mutex.Unlock()

	// 5. Signal Item Available
	// We release 1 token to semItems, waking up a consumer if one is waiting.
	b.semItems.Release(1)
}

// Get removes an item
func (b *SafeBuffer) Get() Event {
	ctx := context.TODO()

	// 1. Wait for Item
	if err := b.semItems.Acquire(ctx, 1); err != nil {
		fmt.Printf("Failed to acquire item: %v\n", err)
		return Event{}
	}

	// 2. Lock Mutex
	b.mutex.Lock()

	// 3. Get from Buffer
	event := b.queue[0]
	b.queue = b.queue[1:] // Dequeue
	fmt.Printf("\tConsumer processed Event %d. (Buffer Size: %d)\n", event.ID, len(b.queue))

	// 4. Unlock Mutex
	b.mutex.Unlock()

	// 5. Signal Space Available
	b.semSpaces.Release(1)

	return event
}

func producer(id int, buf *SafeBuffer, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		// Simulate production time
		time.Sleep(time.Millisecond * 100)
		buf.Put(Event{ID: (id * 100) + i})
	}
}

func consumer(id int, buf *SafeBuffer, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		// Simulate processing time
		time.Sleep(time.Millisecond * 150)
		_ = buf.Get()
	}
}

func main() {
	bufferSize := int64(5) // Finite buffer size
	numProducers := 2
	numConsumers := 2
	itemsPerThread := 5

	fmt.Printf("Starting Producer-Consumer (Buffer Size: %d)\n", bufferSize)

	// Hack to make semItems start at 0:
	// 1. Create with max capacity
	// 2. Acquire ALL tokens immediately so it is "empty"
	buf := NewSafeBuffer(bufferSize)
	// Drain the items semaphore so it starts at 0 available
	if err := buf.semItems.Acquire(context.TODO(), bufferSize); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(numProducers + numConsumers)

	// Start Producers
	for i := 1; i <= numProducers; i++ {
		go producer(i, buf, itemsPerThread, &wg)
	}

	// Start Consumers
	for i := 1; i <= numConsumers; i++ {
		go consumer(i, buf, itemsPerThread, &wg)
	}

	wg.Wait()
	fmt.Println("All items produced and consumed. Exiting.")
}
