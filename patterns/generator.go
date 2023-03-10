package patterns

// Generator: function that returns a channel
// https://go.dev/talks/2012/concurrency.slide
func IntGenerator(from, to int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for i := from; i < to; i++ {
			c <- i
		}
	}()
	return c
}
