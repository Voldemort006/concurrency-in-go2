package replicating_requests

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {

	doWork := func(done <-chan interface{},
		id int, wg *sync.WaitGroup, result chan<- int) {
		started := time.Now()
		defer wg.Done()

		simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
		select {
		case <-done:
			fmt.Println(id, "done closed")
		case <-time.After(simulatedLoadTime):
		}

		select {
		case <-done:
			fmt.Println(id, "done closed again")
		case result <- id:
		}

		took := time.Since(started)
		if took < simulatedLoadTime {
			fmt.Println("less time = ", id, took, simulatedLoadTime)
			took = simulatedLoadTime
		}
		fmt.Printf("%v took %v %v\n", id, took, simulatedLoadTime)
	}

	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 1; i <= 10; i++ {
		go doWork(done, i, &wg, result)
	}

	firstReturned := <-result

	close(done)
	wg.Wait()

	fmt.Println("firstReturned = ", firstReturned)
}
