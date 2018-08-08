package sync

import (
	"sync"
)

type memoryMutexer struct {
	sync.Map
}

func InitMemory() {
	DefaultMutexer = NewMemoryMutexer()
}

func NewMemoryMutexer() Mutexer {
	return &memoryMutexer{}
}

func (m *memoryMutexer) NewMutex(key string) Mutex {
	v, ok := m.Load(key)
	if ok {
		mu, ok := v.(*sync.Mutex)
		if ok && mu != nil {
			return mu
		}
	}

	mu := &sync.Mutex{}
	m.Store(key, mu)
	return mu
}
