package messages

import (
	"context"

	"github.com/porebric/green-api-test/internal/messages/models"
)

type Provider interface {
	GetMessage(context.Context, int64) (models.Message, error)
	SaveMessage(context.Context, models.Message) error
}
