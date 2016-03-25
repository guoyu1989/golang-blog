package main

import (
	"flag"
	"fmt"
	"sync"
)

const (
	secSquaringNumber = 1
	secFanInFanOut    = 2
	secExplicitCancel = 3
	numFans           = 2
)

var sec = flag.Int("sec", secFanInFanOut, `
    The integer number indicate sections within chapter
      1 - Squaring numbers
      2 - Fan in, fan out
      3 - Explicit Cancellation
`)

var fans = flag.Int("fans", numFans, "Number of fans")

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()
	return out
}

func merge(done *chan struct{}, ins ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			if done != nil {
				select {
				case out <- n:
				case <-*done:
					fmt.Println("Closed")
					return
				}
			} else {
				out <- n
			}
		}
	}

	wg.Add(len(ins))
	for _, in := range ins {
		go output(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	flag.Parse()
	if *sec == secSquaringNumber {
		for n := range sq(sq(gen(1, 2, 3))) {
			fmt.Println(n)
		}
	} else {
		ins := gen(1, 2, 3)
		fans := make([]<-chan int, numFans)
		for i := 0; i < numFans; i++ {
			fans[i] = sq(ins)
		}
		if *sec == secFanInFanOut {
			for n := range merge(nil, fans...) {
				fmt.Println(n)
			}
		} else if *sec == secExplicitCancel {
			done := make(chan struct{})
			defer close(done)
			for n := range merge(&done, fans...) {
				fmt.Println(n)
			}
		}
	}
}
