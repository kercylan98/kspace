package main

import (
	"fmt"
	"math/rand"
)

func main() {

	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		for i := 0; i < 5; i++ {
			ch <- rand.Intn(100)
		}

		close(ch)
	}()

	go func() {
		for {
			if data, isClose := <-ch; !isClose {
				break
			} else {
				fmt.Println(data)
			}
		}
		done <- struct{}{}
	}()

	<-done
}
