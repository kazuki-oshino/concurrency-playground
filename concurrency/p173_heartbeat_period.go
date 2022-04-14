package concurrency

import "time"

func DoWorkPeriod(done <-chan interface{}, pulseInterval time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	intStream := make(chan int)

	go func() {
		defer close(heartbeatStream)
		defer close(intStream)

		time.Sleep(2 * time.Second)
		pulse := time.Tick(pulseInterval)

	numLoop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeatStream <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numLoop
				}
			}
		}
	}()

	return heartbeatStream, intStream
}
