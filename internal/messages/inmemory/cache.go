package inmemory

import (
	"sync"

	"github.com/porebric/green-api-test/internal/messages/models"
)

type cache struct {
	cache map[int64]models.Message
	mu    sync.RWMutex
}

func newCache() *cache {
	return &cache{
		cache: make(map[int64]models.Message),
	}
}

func (mc *cache) get(id int64) (models.Message, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if msg, found := mc.cache[id]; found {
		return msg, nil
	}
	return models.Message{}, nil
}

func (mc *cache) set(msg models.Message) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cache[msg.Id] = msg
	return nil
}
