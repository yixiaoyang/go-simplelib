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

func (s *LruTestSuite) TestAddGet() {
	for i := 0; i < 3; i++ {
		s.lru.Add(i, fmt.Sprintf("I'm %v", i))
		fmt.Printf("test %v", i)
		value, ok := s.lru.Get(i)
		s.True(ok)
		s.Equal(fmt.Sprintf("I'm %v", i), value)
	}

	s.Equal(3, s.lru.Len())
}

func TestLruTestSuite(t *testing.T) {
	suite.Run(t, new(LruTestSuite))
}
