package clickhouse

import (
	"context"
	"fmt"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/event"
	"github.com/antoniokichaev/hezzl-collector/pkg/clickhouse"
)

type EventRepo struct {
	*clickhouse.Clickhouse
}

func (er *EventRepo) CreateEvent(ctx context.Context, clickhouseEvents []event.ClickhouseEvent) error {
	insertBuilder := er.Builder.Insert("items").
		Columns(
			"Id",
			"CampaignId",
			"Name",
			"Description",
			"Priority",
			"Removed",
			"EventTime",
		)

	for _, ce := range clickhouseEvents {
		insertBuilder = insertBuilder.Values(ce.Id, ce.CampaignId, ce.Name, ce.Description, ce.Priority, ce.Removed, ce.EventTime)
	}
	sqlReq, args, _ := insertBuilder.ToSql()

	err := er.Clickhouse.DB.Exec(ctx, sqlReq, args...)
	if err != nil {
		return fmt.Errorf("clickhouse.CreateEvent Exec %w", err)
	}
	return nil
}

func NewEventRepo(cl *clickhouse.Clickhouse) *EventRepo {
	return &EventRepo{Clickhouse: cl}

}
