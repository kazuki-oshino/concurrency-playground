package concurrency

import (
	"fmt"
)

func P72() {
	ch := make(chan string, 2)
	ch <- "aaa"
	ch <- "bbb"
	close(ch)
	for s := range ch {
		fmt.Println(s)
	}
	// channelのバッファを0にしている場合、同時消化できないと死ぬ
}
