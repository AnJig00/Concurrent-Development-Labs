//Barrier.go Template Code
//Copyright (C) 2024 Dr. Joseph Kehoe

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

//--------------------------------------------
// Author: Joseph Kehoe (Joseph.Kehoe@setu.ie)
// Created on 30/9/2024
// Modified by: Shuai Jiang
//--------------------------------------------

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// Global variables for the barrier
// In a production app, these might be in a struct, but we use globals/pointers here for the lab structure.
var (
	count int        // Keeps track of how many routines have reached the barrier
	lock  sync.Mutex // Protects access to the 'count' variable
)

// doStuff performs part A, waits at the barrier, and then performs part B.
// Passing the synchronization primitives as pointers.
func doStuff(goNum int, total int, wg *sync.WaitGroup, sem *semaphore.Weighted) {
	defer wg.Done()

	// 1. EXECUTE PART A
	// Simulate work with Sleep
	time.Sleep(time.Second)
	fmt.Printf("Routine %d: Finished Part A\n", goNum)

	// --- BARRIER START ---
	lock.Lock()
	count++ // Increment the counter as this routine has arrived

	if count == total {
		// If this is the LAST routine to arrive:
		// Release (total - 1) tokens into the semaphore.
		// This "opens the gate" for all the other waiting routines.
		// We subtract 1 because the current routine doesn't need to wait on the semaphore.
		fmt.Println("\n--- All routines arrived. Barrier Releasing! ---\n")
		sem.Release(int64(total - 1))
		lock.Unlock()
	} else {
		// If this is NOT the last routine:
		lock.Unlock() // Release the mutex so others can increment the count

		// Attempt to acquire a token from the semaphore.
		// Since main() drained the semaphore initially, this will BLOCK
		// until the last routine performs the Release().
		if err := sem.Acquire(context.TODO(), 1); err != nil {
			fmt.Printf("Routine %d: Failed to acquire semaphore: %v\n", goNum, err)
		}
	}
	// --- BARRIER END ---

	// 2. EXECUTE PART B
	// All routines execute this only after the barrier is lifted.
	fmt.Printf("Routine %d: Starting Part B\n", goNum)
}

func main() {
	totalRoutines := 10
	var wg sync.WaitGroup

	// Context is required for the semaphore Acquire method
	ctx := context.TODO()

	// Initialize the weighted semaphore with capacity = totalRoutines
	sem := semaphore.NewWeighted(int64(totalRoutines))

	// We acquire all tokens now so that the semaphore is "empty".
	// This ensures that any goroutine trying to Acquire(1) later will block (wait).
	if err := sem.Acquire(ctx, int64(totalRoutines)); err != nil {
		fmt.Printf("Failed to acquire initial semaphore: %v", err)
		return
	}

	fmt.Println("Starting Goroutines...")

	wg.Add(totalRoutines)
	for i := 0; i < totalRoutines; i++ {
		// Pass necessary variables to the function
		go doStuff(i, totalRoutines, &wg, sem)
	}

	wg.Wait() // Wait for everyone to finish Part B before exiting
	fmt.Println("All routines finished. Exiting.")
}
