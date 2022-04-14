package concurrency

import (
	"testing"
)

func TestDoWork(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := DoWorkT(done, intSlice...)

	<-heartbeat

	i := 0
	for r := range results {

		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		i++
	}
}
