package main

import "sync"

type InMemoryPMStore struct {
	store map[string]int
	// synchronization of read/write operations
	lock sync.RWMutex
}

func NewInMemoryPMStore() *InMemoryPMStore {
	return &InMemoryPMStore{
		map[string]int{},
		sync.RWMutex{},
	}
}

func (i *InMemoryPMStore) Append(name string) {
	i.lock.Lock()
	i.store[name]++
	i.lock.Unlock()
}

func (i *InMemoryPMStore) GetDoneTasks(name string) int {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.store[name]
}

func (i *InMemoryPMStore) GetProjectInfo() []Worker {
	var projectInfo []Worker
	for name, doneTasks := range i.store {
		projectInfo = append(projectInfo, Worker{name, doneTasks})
	}
	return projectInfo
}
