package inmemcache

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Cache struct {
	cleanUpTime time.Duration
	codes       map[uuid.UUID]int
	timers      map[uuid.UUID]*time.Timer // Храним таймеры для каждого чата
	mu          sync.RWMutex
}

func New(cleanUpTime time.Duration) *Cache {
	return &Cache{
		cleanUpTime: cleanUpTime,
		codes:       make(map[uuid.UUID]int),
		timers:      make(map[uuid.UUID]*time.Timer),
		mu:          sync.RWMutex{},
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
		// Удаляем запись и таймер по истечении времени.
		delete(s.codes, user)
		delete(s.timers, user)
		s.mu.Unlock()
	})

	s.timers[user] = timer
}

// Устанавливаем код для пользователя и запускаем таймер для его удаления.
func (s *Cache) Set(user uuid.UUID, code int) {
	s.mu.Lock()
	s.codes[user] = code
	s.mu.Unlock()
	s.cleanUpAfter(user, s.cleanUpTime)
}

// Получаем код для пользователя.
func (s *Cache) Get(user uuid.UUID) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.codes[user]

	return data, ok
}

// Удаляем код и таймер для пользователя.
func (s *Cache) Delete(user uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.codes, user)

	if timer, exists := s.timers[user]; exists {
		timer.Stop()
		delete(s.timers, user)
	}
}

// Принудительное удаление всех устаревших кодов.
func (s *Cache) CleanUpExpiredCodes() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for user, timer := range s.timers {
		timer.Stop()
		delete(s.codes, user)
		delete(s.timers, user)
	}
}
