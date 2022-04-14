package concurrency

import (
	"fmt"
	"math/rand"
	"time"
)

func doWorkT(done <-chan interface{}) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	workStream := make(chan int)

	go func() {
		defer close(heartbeatStream)
		defer close(workStream)

		for i := 0; i < 10; i++ {

			select {
			case heartbeatStream <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case workStream <- rand.Intn(10):
			}
		}
	}()

	return heartbeatStream, workStream
}

func DoWorkT(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	intStream := make(chan int)

	go func() {
		defer close(heartbeatStream)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		for _, n := range nums {

			select {
			case heartbeatStream <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()

	return heartbeatStream, intStream
}

func P169HeartBeat() {
	done := make(chan interface{})
	defer close(done)

	heartbeat, results := doWorkT(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				fmt.Println("heartbeat die")
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				fmt.Println("results die")
				return
			}
			fmt.Printf("results %v\n", r)
		}
	}
}
