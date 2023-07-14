package mnats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/event"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/items"
	"github.com/nats-io/nats.go"
	"time"
)

type Item interface {
	CreateItem(ctx context.Context, name string, campaignId int) (items.Item, error)
	UpdateItem(ctx context.Context, name, description string, id, campaignId int) (items.Item, error)
	DeleteItem(ctx context.Context, id, campaignId int) (items.Item, error)
	GetItem(ctx context.Context, id, campaignId int) (items.Item, error)
	GetItems(ctx context.Context) ([]items.Item, error)
}

type ItemRepo struct {
	Item
	stream nats.JetStreamContext
}

func (i *ItemRepo) UpdateItem(ctx context.Context, name, description string, id, campaignId int) (items.Item, error) {
	ce := &event.ClickhouseEvent{}
	it, err := i.Item.UpdateItem(ctx, name, description, id, campaignId)

	if err != nil {
		return it, err
	}

	ce.Id = it.ID
	ce.CampaignId = it.CampaignID
	ce.Name = it.Name
	ce.Description = it.Description
	ce.Priority = it.Priority
	ce.Removed = it.Removed
	ce.EventTime = time.Now()
	ce.Description = it.Description
	i.sendEvent(ce)
	return it, err
}
func (i *ItemRepo) DeleteItem(ctx context.Context, id, campaignId int) (items.Item, error) {
	it, err := i.Item.DeleteItem(ctx, id, campaignId)

	ce := &event.ClickhouseEvent{
		Id:          id,
		CampaignId:  campaignId,
		EventTime:   time.Now(),
		Name:        it.Name,
		Description: it.Description,
		Priority:    it.Priority,
		Removed:     it.Removed,
	}
	i.sendEvent(ce)

	return it, err
}

func (i *ItemRepo) sendEvent(event *event.ClickhouseEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	s, err := i.stream.Publish("events.item", data)
	fmt.Println(s.Domain, s.Stream, s.Sequence, err, event)

}

func NewItemRepo(repo Item, js nats.JetStreamContext) (*ItemRepo, error) {
	cfg := &nats.StreamConfig{
		Name:      "EVENTS",
		Subjects:  []string{"events.>"},
		Retention: nats.WorkQueuePolicy,
	}

	_, err := js.AddStream(cfg)
	if err != nil {
		panic(err)
	}
	return &ItemRepo{
		Item:   repo,
		stream: js,
	}, nil
}
