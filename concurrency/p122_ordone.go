package concurrency

import "fmt"

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func intToInterface(done <-chan interface{}, target <-chan int) <-chan interface{} {
	interfaceStream := make(chan interface{})
	go func() {
		defer close(interfaceStream)
		for t := range target {
			select {
			case <-done:
				return
			case interfaceStream <- t:
			}
		}
	}()
	return interfaceStream
}

func P122OrDone() {
	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 3), 3)

	for v := range orDone(done, intToInterface(done, pipeline)) {
		fmt.Println(v)
	}
}
