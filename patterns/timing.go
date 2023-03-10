package patterns

import (
	"fmt"
	"time"
)

func TimeCost(start time.Time, name string) {
	fmt.Printf("[%v] cost %v\n", name, time.Since(start))
}

func LongTimeWorking() {
	defer TimeCost(time.Now(), "LongTimeWorking")
	time.Sleep(time.Millisecond * 500)
}

func TimeoutWithSelect() {
	c := make(chan int)
	timeout := time.After(2 * time.Second)
	for {
		select {
		case v := <-c:
			fmt.Printf("recv %v\n", v)
		case <-timeout:
			fmt.Println("timeout")
			return
		}
	}
}
