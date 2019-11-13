package configer

import "sync"

type singleton struct {
	locker  sync.Mutex
	content map[string]string
}

func NewSingleton(once sync.Once) singleton {
	var instance singleton
	once.Do(func() {
		locker := sync.Mutex{}
		content := make(map[string]string)
		instance = singleton{locker, content}
	})

	return instance
}

func (self singleton) Set(key string, value string) {
	self.locker.Lock()
	self.content[key] = value
	self.locker.Unlock()
}

func (self singleton) Get(key string) string {
	return self.content[key]
}

func (self *singleton) ClearAll() {
	self.locker.Lock()
	self.content = make(map[string]string)
	self.locker.Unlock()
}
