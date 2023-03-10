package patterns

import (
	"container/list"
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PatternTestSuite struct {
	suite.Suite
}

func (s *PatternTestSuite) SetupTest() {
}

func (s *PatternTestSuite) TestNewSingletone() {
	var wg sync.WaitGroup
	count := 32
	singletons := make([]*Singleton, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			singletons[id] = NewSingleton()
			defer wg.Done()
		}(i)
	}

	wg.Wait()
	for i := 0; i < count; i++ {
		s.Equal(singletons[0], singletons[i])
		s.Equal(1, singletons[i].Id)
	}
}

func (s *PatternTestSuite) TestLongTimeWorking() {
	LongTimeWorking()
}

func (s *PatternTestSuite) TestObserver() {
	count := 32
	notifier := &ChatNotifier{
		observers: make(map[Observer]struct{}),
	}
	observers := make([]*ChatObserver, count)
	for i := 0; i < count; i++ {
		observers[i] = &ChatObserver{
			Name:      fmt.Sprintf("%v", i),
			EventList: list.New(),
		}
		notifier.Add(observers[i])
	}

	notifier.Remove(observers[31])

	notifier.Notify(Event{
		Msg: "first message",
	})
	notifier.Add(observers[31])
	notifier.Remove(observers[30])
	notifier.Notify(Event{
		Msg: "second message",
	})

	for i := 0; i < 30; i++ {
		s.Equal(2, observers[i].EventList.Len())

		iterator := observers[i].EventList.Front()
		s.Equal("first message", iterator.Value.(Event).Msg)
		iterator = iterator.Next()
		s.Equal("second message", iterator.Value.(Event).Msg)
	}

	s.Equal(31, len(notifier.observers))

	s.Equal(1, observers[30].EventList.Len())
	s.Equal(1, observers[31].EventList.Len())

	s.Equal("first message", observers[30].EventList.Front().Value.(Event).Msg)
	s.Equal("second message", observers[31].EventList.Front().Value.(Event).Msg)
}

func (s *PatternTestSuite) TestIntGenerator() {
	from := 5
	except := from
	for i := range IntGenerator(from, 10) {
		s.Equal(except, i)
		except += 1
	}
}

func (s *PatternTestSuite) TestFanIn1() {
	TestFanIn := func(f func(input1, input2 <-chan int) <-chan int) {
		input1 := make(chan int)
		input2 := make(chan int)
		fanOut := f(input1, input2)

		go func() {
			input1 <- 0
			input2 <- 1
			input1 <- 2
			input2 <- 3
		}()

		values := make([]int, 4)
		for i := 0; i < 4; i++ {
			values[i] = <-fanOut
		}
		sort.Ints(values)
		for i := 0; i < 4; i++ {
			s.Equal(i, values[i])
		}
	}

	TestFanIn(FanIn)
	TestFanIn(FanIn2)
}

func (s *PatternTestSuite) TestFanOut() {
	exitChan := make(chan int)
	input := make(chan int)
	output1 := make(chan int)
	output2 := make(chan int)
	FanOut(input, []chan<- int{output1, output2}, exitChan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 4; i++ {
			input <- i
			fmt.Printf("input <- %v\n", i)
		}
		fmt.Printf("write routine done\n")
	}()

	count := 0
	values := make([]int, 4)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer func() {
			wg2.Done()
			fmt.Printf("recv routine exit\n")
		}()
		for {
			select {
			case v := <-output1:
				fmt.Printf("%v <- output1\n", v)
				values[count] = v
				count++
			case v := <-output2:
				values[count] = v
				fmt.Printf("%v <- output2\n", v)
				count++
			}
			if count >= 4 {
				return
			}
		}
	}()

	wg.Wait()
	close(exitChan)

	wg2.Wait()
	sort.Ints(values)
	for i := 0; i < 4; i++ {
		s.Equal(i, values[i])
	}
}

func (s *PatternTestSuite) TestVisitor() {
	circle := Circle{R: 10}
	rectangle := Rectangle{W: 5, H: 20}

	circle.Accept(XmlVisitor)
	circle.Accept(JsonVisitor)
	rectangle.Accept(XmlVisitor)
	rectangle.Accept(JsonVisitor)
}

func TestPatternTestSuite(t *testing.T) {
	suite.Run(t, new(PatternTestSuite))
}
