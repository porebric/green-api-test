package inmemory

import (
	"context"

	"github.com/porebric/green-api-test/internal/instance/models"
)

type provider struct {
	cacheInstance *cache[models.Instance]
	cacheSettings *cache[models.Settings]
}

func NewProvider() *provider {
	p := &provider{
		cacheInstance: newCache[models.Instance](),
		cacheSettings: newCache[models.Settings](),
	}

	p.cacheInstance.set(models.Instance{
		Id:    1,
		Name:  "test1",
		Token: "test_token_1",
	})
	p.cacheInstance.set(models.Instance{
		Id:    2,
		Name:  "test2",
		Token: "test_token_2",
	})
	p.cacheSettings.set(models.Settings{
		InstanceId: 1,
		Settings: []models.Setting{
			{
				Id:    1,
				Name:  "setting_1",
				Value: "val1",
			},
			{
				Id:    2,
				Name:  "setting_2",
				Value: "val2",
			},
		},
	})
	p.cacheSettings.set(models.Settings{
		InstanceId: 1,
		Settings: []models.Setting{
			{
				Id:    3,
				Name:  "setting_1",
				Value: "val1",
			},
			{
				Id:    4,
				Name:  "setting_2",
				Value: "val2",
			},
		},
	})

	return p
}

func (p *provider) GetInstance(_ context.Context, id int64) (models.Instance, error) {
	return p.cacheInstance.get(id)
}

func (p *provider) GetSettings(_ context.Context, id int64) (models.Settings, error) {
	return p.cacheSettings.get(id)
}
