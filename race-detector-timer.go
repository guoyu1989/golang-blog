package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// The demo with race happens
	// raceDemo()

	// The demo without race
	raceFreeDemo()
}

func raceDemo() {
	start := time.Now()
	var t *time.Timer

	t = time.AfterFunc(randomDuration(), func() {
		fmt.Println(time.Now().Sub(start))
		t.Reset(randomDuration())
	})

	time.Sleep(10 * time.Second)
}

func raceFreeDemo() {
	start := time.Now()
	reset := make(chan bool)
	var t *time.Timer

	t = time.AfterFunc(randomDuration(), func() {
		fmt.Println(time.Now().Sub(start))
		reset <- true
	})
	for time.Since(start) < 5*time.Second {
		<-reset
		t.Reset(randomDuration())
	}
}

func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}
