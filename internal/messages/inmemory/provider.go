package inmemory

import (
	"context"

	"github.com/porebric/green-api-test/internal/messages/models"
)

type provider struct {
	cache *cache
}

func NewProvider() *provider {
	return &provider{
		cache: newCache(),
	}
}

func (p *provider) GetMessage(_ context.Context, id int64) (models.Message, error) {
	return p.cache.get(id)
}

func (p *provider) SaveMessage(_ context.Context, msg models.Message) error {
	return p.cache.set(msg)
}
