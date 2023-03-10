package patterns

import "sync"

type Singleton struct {
	Data map[string]string
	Id   int
}

var (
	once      sync.Once
	singleton *Singleton
	idCount   int
)

func NewSingleton() *Singleton {
	once.Do(func() {
		idCount += 1
		singleton = &Singleton{
			Data: make(map[string]string),
			Id:   idCount,
		}
	})
	return singleton
}
