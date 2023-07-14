package service

import (
	"context"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/campaign"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/event"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/items"
	"github.com/antoniokichaev/hezzl-collector/internal/repo"
	serviceCompaign "github.com/antoniokichaev/hezzl-collector/internal/service/campaign"
	"github.com/antoniokichaev/hezzl-collector/internal/service/eventSaver"
	"github.com/antoniokichaev/hezzl-collector/internal/service/item"
	"github.com/nats-io/nats.go"
)

type Campaign interface {
	CreateCampaign(ctx context.Context, name string) (campaign.Campaign, error)
	UpdateCampaign(ctx context.Context, name string, id int) (campaign.Campaign, error)
	DeleteCampaign(ctx context.Context, id int) (campaign.Campaign, error)
	GetCampaign(ctx context.Context)
	GetCampaigns(ctx context.Context) ([]campaign.Campaign, error)
}

type Item interface {
	CreateItem(ctx context.Context, name string, campaignId int) (items.Item, error)
	UpdateItem(ctx context.Context, name, description string, id, campaignId int) (items.Item, error)
	DeleteItem(ctx context.Context, id, campaignId int) (items.Item, error)
	GetItem(ctx context.Context, id, campaignId int) (items.Item, error)
	GetItems(ctx context.Context) ([]items.Item, error)
}
type EventSaver interface {
	Start(ctx context.Context)
}

type Event interface {
	CreateEvent(ctx context.Context, event []event.ClickhouseEvent) error
}

type Services struct {
	Item       Item
	Campaign   Campaign
	EventSaver EventSaver
}

type SDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps *SDependencies, js nats.JetStreamContext) *Services {
	return &Services{
		Campaign:   serviceCompaign.NewCampaingUseCase(deps.Repos),
		Item:       item.NewItemUseCase(deps.Repos),
		EventSaver: eventSaver.New(deps.Repos.Event, js),
	}
}
