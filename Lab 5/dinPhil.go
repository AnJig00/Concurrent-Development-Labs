// dinPhil.go
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
// Description: Dining Philosophers solution using Asymmetry to prevent Deadlock
//--------------------------------------------

package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func think(index int) {
	var X time.Duration
	X = time.Duration(rand.IntN(5))
	time.Sleep(X * time.Second) // wait random time amount
	fmt.Println("Phil: ", index, "was thinking")
}

func eat(index int) {
	var X time.Duration
	X = time.Duration(rand.IntN(5))
	time.Sleep(X * time.Second) // wait random time amount
	fmt.Println("Phil: ", index, "was eating")
}

// getForks implements the solution to avoid Deadlock.
// We will use the "Asymmetry Solution":
// One philosopher (e.g., Phil 0) picks up the right fork first.
// Everyone else picks up the left fork first.
func getForks(index int, forks map[int]chan bool) {
	if index == 0 {
		// Phil 0 picks up right (index+1) then left (index)
		forks[(index+1)%5] <- true
		forks[index] <- true
	} else {
		// Everyone else picks up left (index) then right (index+1)
		forks[index] <- true
		forks[(index+1)%5] <- true
	}
}

func putForks(index int, forks map[int]chan bool) {
	// Release forks
	<-forks[index]
	<-forks[(index+1)%5]
}

func doPhilStuff(index int, wg *sync.WaitGroup, forks map[int]chan bool) {
	for i := 0; i < 3; i++ {
		think(index)
		getForks(index, forks)
		eat(index)
		putForks(index, forks)
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	philCount := 5
	wg.Add(philCount)

	// Create channels to act as Mutexes for the forks
	forks := make(map[int]chan bool)
	for k := range philCount {
		forks[k] = make(chan bool, 1)
	} //set up forks

	fmt.Println("Dining Philosophers Problem")

	for N := range philCount {
		go doPhilStuff(N, &wg, forks)
	} //start philosophers

	wg.Wait() //wait here until everyone is done
	fmt.Println("All philosophers are full.")

} //main
