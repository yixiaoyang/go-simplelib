package patterns

func FanIn(input1, input2 <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			c <- <-input1
		}
	}()

	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}

func FanIn2(input1, input2 <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-input1:
				c <- v
			case v := <-input2:
				c <- v
			}
		}
	}()
	return c
}

func FanOut(input <-chan int, outputs []chan<- int, exitChan <-chan int) {
	for _, output := range outputs {
		go func(out chan<- int) {
			for {
				select {
				case v, ok := <-input:
					{
						if !ok {
							return
						}
						out <- v
					}
				case <-exitChan:
					{
						return
					}
				}
			}
		}(output)
	}
}
