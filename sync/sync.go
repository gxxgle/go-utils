package sync

var (
	DefaultMutexer = NewMemoryMutexer()
)

type Mutex interface {
	Lock()
	Unlock()
}

type Mutexer interface {
	NewMutex(string) Mutex
}

func NewMutex(key string) Mutex {
	return DefaultMutexer.NewMutex(key)
}
