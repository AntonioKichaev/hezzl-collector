package repo

import (
	"context"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/campaign"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/event"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/items"
	clickhouse2 "github.com/antoniokichaev/hezzl-collector/internal/repo/clickhouse"
	"github.com/antoniokichaev/hezzl-collector/internal/repo/mnats"
	"github.com/antoniokichaev/hezzl-collector/internal/repo/pgdb"
	redisRepo "github.com/antoniokichaev/hezzl-collector/internal/repo/redis"
	"github.com/antoniokichaev/hezzl-collector/pkg/clickhouse"
	"github.com/antoniokichaev/hezzl-collector/pkg/postgres"
	"github.com/nats-io/nats.go"
	redis "github.com/redis/go-redis/v9"
)

//go:generate mockery --name Campaign
type Campaign interface {
	CreateCampaign(ctx context.Context, name string) (campaign.Campaign, error)
	UpdateCampaign(ctx context.Context, name string, id int) (campaign.Campaign, error)
	DeleteCampaign(ctx context.Context, id int) (campaign.Campaign, error)
	GetCampaign(ctx context.Context, id int) (campaign.Campaign, error)
	GetCampaigns(ctx context.Context) ([]campaign.Campaign, error)
}

//go:generate mockery --name Item
type Item interface {
	CreateItem(ctx context.Context, name string, campaignId int) (items.Item, error)
	UpdateItem(ctx context.Context, name, description string, id, campaignId int) (items.Item, error)
	DeleteItem(ctx context.Context, id, campaignId int) (items.Item, error)
	GetItem(ctx context.Context, id, campaignId int) (items.Item, error)
	GetItems(ctx context.Context) ([]items.Item, error)
}
type Event interface {
	CreateEvent(ctx context.Context, event []event.ClickhouseEvent) error
}

type Repositories struct {
	Campaign
	Item
	Event
}

func NewRepositories(db *postgres.Postgres, clDb *clickhouse.Clickhouse, client *redis.Client, js nats.JetStreamContext) *Repositories {
	itemRepo := pgdb.NewItemRepo(db)
	campaignRepo := pgdb.NewCampaignRepo(db)

	redisItemRepo := redisRepo.NewItemRepo(itemRepo, client)
	redisCampaignRepo := redisRepo.NewCampaignRepo(campaignRepo, client)

	natsItemRepo, _ := mnats.NewItemRepo(redisItemRepo, js)
	eventRepo := clickhouse2.NewEventRepo(clDb)

	return &Repositories{
		Campaign: redisCampaignRepo,
		Item:     natsItemRepo,
		Event:    eventRepo,
	}
}
