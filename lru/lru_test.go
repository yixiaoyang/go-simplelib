package lru

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LruTestSuite struct {
	suite.Suite
	lru *Lru[int, string]
}

func (s *LruTestSuite) SetupTest() {
	s.lru = NewLru[int, string]()
	s.NotNil(s.lru)
}

func lru_print(key int, value string) bool {
	fmt.Printf("%v=%v\n", key, value)
	return false
}

func (s *LruTestSuite) TestAddGet() {
	for i := 0; i < 3; i++ {
		s.lru.Add(i, fmt.Sprintf("I'm %v", i))
		value, ok := s.lru.Get(i)
		s.True(ok)
		s.Equal(fmt.Sprintf("I'm %v", i), value)
	}

	s.Equal(3, s.lru.Len())

	fmt.Println("---")
	s.lru.Iterate(lru_print)
	fmt.Println("---")
	s.lru.Get(1)
	s.lru.Iterate(lru_print)
}

func TestLruTestSuite(t *testing.T) {
	suite.Run(t, new(LruTestSuite))
}
