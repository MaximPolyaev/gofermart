package mutex

import (
	"sync"
)

const lockPortions = 10

type Mutex struct {
	locks map[lockPortion]mutexMap
}

type mutexMap struct {
	vals map[int]*sync.Mutex
	lock *sync.Mutex
}

type lockPortion int8

func New() *Mutex {
	locks := make(map[lockPortion]mutexMap, lockPortions)

	for i := 0; i < lockPortions; i++ {
		locks[lockPortion(i)] = mutexMap{
			vals: make(map[int]*sync.Mutex),
			lock: &sync.Mutex{},
		}
	}

	return &Mutex{locks: locks}
}

func (m *Mutex) Lock(val int) {
	locker := m.getLocker(val)
	locker.Lock()
}

func (m *Mutex) Unlock(val int) {
	locker := m.getLocker(val)
	locker.Unlock()
}

func (m *Mutex) getLocker(val int) *sync.Mutex {
	portion := m.getLockPotion(val)

	mutexMap := m.locks[portion]
	mutexMap.lock.Lock()
	defer mutexMap.lock.Unlock()

	mutex, ok := mutexMap.vals[val]
	if ok {
		return mutex
	}

	mutexMap.vals[val] = &sync.Mutex{}

	return mutexMap.vals[val]
}

func (m *Mutex) getLockPotion(val int) lockPortion {
	return lockPortion(val % lockPortions)
}
