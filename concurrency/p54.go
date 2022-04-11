package concurrency

import (
	"fmt"
	"sync"
	"time"
)

func P54() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration, i int) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Remove from queue : ", i)
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to queue : ", i, queue)
		queue = append(queue, struct{}{})
		go removeFromQueue(time.Second*1, i)
		c.L.Unlock()
	}

}
