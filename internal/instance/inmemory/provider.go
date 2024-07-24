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

	p.cacheSettings.set(models.Settings{
		InstanceId:                        1,
		Wid:                               "11001234567@c.us",
		CountryInstance:                   "",
		TypeAccount:                       "",
		WebhookUrl:                        "https://mysite.com/webhook/green-api/",
		WebhookUrlToken:                   "",
		DelaySendMessagesMilliseconds:     5000,
		MarkIncomingMessagesReaded:        "no",
		MarkIncomingMessagesReadedOnReply: "no",
		SharedSession:                     "no",
		OutgoingWebhook:                   "yes",
		OutgoingMessageWebhook:            "yes",
		OutgoingAPIMessageWebhook:         "yes",
		IncomingWebhook:                   "yes",
		DeviceWebhook:                     "no", // Уведомление временно не работает
		StatusInstanceWebhook:             "no",
		StateWebhook:                      "no",
		EnableMessagesHistory:             "no",
		KeepOnlineStatus:                  "no",
		PollMessageWebhook:                "no",
		IncomingBlockWebhook:              "yes", // Уведомление временно не работает
		IncomingCallWebhook:               "yes",
	})

	return p
}

func (p *provider) GetInstance(_ context.Context, id int64) (models.Instance, error) {
	return p.cacheInstance.get(id)
}

func (p *provider) GetSettings(_ context.Context, id int64) (models.Settings, error) {
	return p.cacheSettings.get(id)
}
