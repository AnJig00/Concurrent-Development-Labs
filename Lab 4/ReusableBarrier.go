// ReusableBarrier.go
// Copyright (C) 2024 Dr. Joseph Kehoe

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// --------------------------------------------
// Author: Joseph Kehoe (Joseph.Kehoe@setu.ie)
// Modified by: Shuai Jiang
// Description: Implements a Reusable Two-Phase Barrier
// --------------------------------------------
package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// Create a barrier data type
// Add two channels to act as the "Two Turnstiles" mentioned in the slides
type barrier struct {
	turnstile1 chan bool // First gate (Arrival)
	turnstile2 chan bool // Second gate (Departure)
	theLock    sync.Mutex
	total      int
	count      int
}

// creates a properly initialised barrier
// N == number of threads (go Routines)
func createBarrier(N int) barrier {
	theBarrier := barrier{
		turnstile1: make(chan bool, N),
		turnstile2: make(chan bool, N),
		total:      N,
		count:      0,
	}
	return theBarrier
}

// Method belonging to barrier data type
// Implements the Two-Phase Barrier pattern
func (b *barrier) wait(id int) {
	//Phase 1: Arrival
	b.theLock.Lock()
	b.count++
	if b.count == b.total {
		// If this is the LAST routine to arrive:
		// Unlock Turnstile 1 for everyone
		for i := 0; i < b.total; i++ {
			b.turnstile1 <- true
		}
	}
	b.theLock.Unlock()

	// Everyone waits here until the channel has tokens (Turnstile 1)
	<-b.turnstile1

	//Phase 2: Departure
	// We need this second phase to ensure everyone has passed the first barrier
	// before we reset the counter for the next loop iteration.
	b.theLock.Lock()
	b.count--
	if b.count == 0 {
		// If this is the LAST routine to leave Phase 1:
		// Unlock Turnstile 2 for everyone
		for i := 0; i < b.total; i++ {
			b.turnstile2 <- true
		}
	}
	b.theLock.Unlock()

	// Everyone waits here until the channel has tokens (Turnstile 2)
	<-b.turnstile2
}

func WorkWithRendezvous(wg *sync.WaitGroup, Num int, theBarrier *barrier) bool {
	//Loop 3 times to prove the barrier is reusable
	for i := 0; i < 3; i++ {
		var X time.Duration
		X = time.Duration(rand.IntN(3)) // Reduced sleep time for faster testing
		time.Sleep(X * time.Second)

		fmt.Printf("Routine %d: Arrived at Barrier (Iteration %d)\n", Num, i+1)

		// Rendezvous here
		theBarrier.wait(Num)

		fmt.Printf("Routine %d: Passed Barrier (Iteration %d)\n", Num, i+1)
	}
	wg.Done()
	return true
}

func main() {
	var wg sync.WaitGroup
	threadCount := 5

	fmt.Println("Starting Reusable Barrier with 3 Iteration")

	// Initialize barrier for 5 routines
	barrier := createBarrier(threadCount)

	wg.Add(threadCount)
	for N := range threadCount {
		go WorkWithRendezvous(&wg, N, &barrier)
	}
	wg.Wait() //wait here until everyone is done with all iterations

	fmt.Println("All routines finished all iterations.")
}
