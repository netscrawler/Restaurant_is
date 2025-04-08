package inmemcache

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

const cleanUpTime = 2 * time.Minute

type Cache struct {
	codes  map[uuid.UUID]int
	timers map[uuid.UUID]*time.Timer // Храним таймеры для каждого чата
	mu     sync.RWMutex
}

func New() *Cache {
	return &Cache{
		codes:  make(map[uuid.UUID]int),
		timers: make(map[uuid.UUID]*time.Timer),
		mu:     sync.RWMutex{},
	}
}

// Очистка с отменой старого таймера.
func (s *Cache) cleanUpAfter(user uuid.UUID, dur time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if oldTimer, exists := s.timers[user]; exists {
		oldTimer.Stop()
	}

	timer := time.AfterFunc(dur, func() {
		s.mu.Lock()
		delete(s.codes, user)
		delete(s.timers, user)
		s.mu.Unlock()
	})

	// Сохраняем таймер в map
	s.timers[user] = timer
}

func (s *Cache) Set(user uuid.UUID, code int) {
	s.mu.Lock()
	s.codes[user] = code
	s.mu.Unlock()
	s.cleanUpAfter(user, cleanUpTime)
}

func (s *Cache) Get(user uuid.UUID) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.codes[user]

	return data, ok
}

func (s *Cache) Delete(user uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.codes, user)

	if timer, exists := s.timers[user]; exists {
		timer.Stop()
		delete(s.timers, user)
	}
}
