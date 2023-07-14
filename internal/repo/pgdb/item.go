package pgdb

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/items"
	"github.com/antoniokichaev/hezzl-collector/pkg/postgres"
)

var ErrNotFoundItem = errors.New("errors.item.notFound")

type ItemRepo struct {
	*postgres.Postgres
}

func (i *ItemRepo) CreateItem(ctx context.Context, name string, campaignId int) (items.Item, error) {
	const fName = "CreateItem"
	sqlReq, args, _ := i.Builder.
		Insert(_tableItems).
		Columns("name", "campaign_id").
		Values(name, campaignId).
		Suffix("returning id,created_at,priority").
		ToSql() //todo: что делать если не campaign_id нет в бд

	item := items.Item{
		Name:       name,
		CampaignID: campaignId,
		Removed:    false,
	}
	err := i.Pool.QueryRow(ctx, sqlReq, args...).Scan(&item.ID, &item.CreatedAt, &item.Priority)
	if err != nil {
		return items.Item{}, fmt.Errorf("%s QueryRow %w", fName, err)
	}
	return item, nil

}

func (i *ItemRepo) UpdateItem(ctx context.Context, name, description string, id, campaignId int) (items.Item, error) {
	sqlReq, args, _ := i.Builder.
		Update(_tableItems).
		Set("name", name).
		Set("description", description).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"campaign_id": campaignId}}).
		Suffix("returning id,created_at,priority").
		ToSql() //todo: что делать если не campaign_id нет в бд создавать или бросать ошибку

	item := items.Item{
		Name:        name,
		CampaignID:  campaignId,
		Description: description,
		Removed:     false,
	}
	err := i.Pool.QueryRow(ctx, sqlReq, args...).Scan(&item.ID, &item.CreatedAt, &item.Priority)
	if err != nil {
		return items.Item{}, ErrNotFoundItem
	}
	return item, nil
}

func (i *ItemRepo) DeleteItem(ctx context.Context, id, campaignId int) (items.Item, error) {
	const fName = "UpdateCampaign"
	sqlReq, args, _ := i.Builder.
		Delete(_tableItems).
		Where(
			sq.And{
				sq.Eq{"id": id},
				sq.Eq{"campaign_id": campaignId},
			},
		).
		Suffix("returning *").
		ToSql()
	it := items.Item{}
	err := i.Pool.QueryRow(ctx, sqlReq, args...).Scan(&it.ID, &it.CampaignID, &it.Name, &it.Description, &it.Priority, &it.Removed, &it.CreatedAt)
	if err != nil {
		return items.Item{}, fmt.Errorf("%s %w", fName, ErrNotFoundItem)
	}

	it.Removed = true
	return it, nil
}

func (i *ItemRepo) GetItem(ctx context.Context, id, campaignId int) (items.Item, error) {
	const fName = "GetItem"
	item := items.Item{}
	sqlReq, args, _ := i.Builder.
		Select("*").
		From(_tableItems).
		Where(
			sq.And{
				sq.Eq{"id": id},
				sq.Eq{"campaign_id": campaignId},
			},
		).ToSql()
	err := i.Pool.QueryRow(ctx, sqlReq, args...).Scan(
		&item.ID,
		&item.CampaignID,
		&item.Name,
		&item.Description,
		&item.Priority,
		&item.Removed,
		&item.CreatedAt,
	)
	if err != nil {
		return items.Item{}, fmt.Errorf("%s QueryRow %w", fName, err)
	}

	return item, nil
}

func (i *ItemRepo) GetItems(ctx context.Context) ([]items.Item, error) {
	const fName = "GetCampaigns"
	itemList := make([]items.Item, 0)
	sqlReq, args, _ := i.Builder.Select("*").From(_tableItems).ToSql()
	rows, err := i.Pool.Query(ctx, sqlReq, args...)
	if err != nil {
		return nil, fmt.Errorf("%s Query %w", fName, err)
	}
	for rows.Next() {
		item := items.Item{}
		err = rows.Scan(&item.ID,
			&item.CampaignID,
			&item.Name,
			&item.Description,
			&item.Priority,
			&item.Removed,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s rows.Scan %w", fName, err)
		}
		itemList = append(itemList, item)
	}

	return itemList, nil
}

func NewItemRepo(db *postgres.Postgres) *ItemRepo {
	return &ItemRepo{Postgres: db}
}
