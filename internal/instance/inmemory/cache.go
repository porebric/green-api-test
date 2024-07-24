package inmemory

import (
	"sync"

	"github.com/porebric/green-api-test/internal/instance/models"
)

type instanceModel interface {
	models.Instance | models.Settings

	GetInstanceId() int64
}

type cache[M instanceModel] struct {
	cache map[int64]M
	mu    sync.RWMutex
}

func newCache[M instanceModel]() *cache[M] {
	return &cache[M]{
		cache: make(map[int64]M),
	}
}

func (mc *cache[M]) get(id int64) (M, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if msg, found := mc.cache[id]; found {
		return msg, nil
	}

	var m M
	return m, nil
}

func (mc *cache[M]) set(msg M) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cache[msg.GetInstanceId()] = msg
	return nil
}
