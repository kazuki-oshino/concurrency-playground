package concurrency

import (
	"fmt"
	"sync"
	"time"
)

func toInt(done <-chan interface{}, ch <-chan interface{}) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for c := range ch {
			select {
			case <-done:
				return
			case intStream <- c.(int):
			}
		}
	}()
	return intStream
}

func longTime(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
	lStream := make(chan interface{})
	go func() {
		defer close(lStream)
		for i := range intStream {
			time.Sleep(2 * time.Second)
			select {
			case <-done:
				return
			case lStream <- i:
			}
		}
	}()
	return lStream
}

func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	stream := make(chan interface{})

	l := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case stream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go l(c)
	}

	go func() {
		wg.Wait()
		close(stream)
	}()

	return stream
}

func P119FanOutFanIn() {
	done := make(chan interface{})
	defer close(done)
	start := time.Now()

	intStream := toInt(done, repeat(done, 2))

	numLongTimes := 10
	fmt.Printf("Spinning up %d long times.", numLongTimes)
	longTimes := make([]<-chan interface{}, numLongTimes)
	for i := 0; i < numLongTimes; i++ {
		longTimes[i] = longTime(done, intStream)
	}

	for x := range take(done, fanIn(done, longTimes...), 10) {
		fmt.Println(x)
	}

	fmt.Printf("Time: %v", time.Since(start))
}
