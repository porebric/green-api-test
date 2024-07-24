package instance

import (
	"context"

	"github.com/porebric/green-api-test/internal/instance/models"
)

type Provider interface {
	GetInstance(context.Context, int64) (models.Instance, error)
	GetSettings(context.Context, int64) (models.Settings, error)
}
