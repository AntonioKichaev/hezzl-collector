package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/items"
	"github.com/antoniokichaev/hezzl-collector/internal/repo/pgdb"
	"github.com/redis/go-redis/v9"
	"time"
)

const _defaultExpiration = time.Minute

type ItemRepo struct {
	*pgdb.ItemRepo
	cache *redis.Client
}

func (i *ItemRepo) GetItems(ctx context.Context) ([]items.Item, error) {
	val, err := i.cache.Get(ctx, "GetItems").Bytes()
	if err != nil {
		itemsList, err := i.ItemRepo.GetItems(ctx)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(itemsList)
		if err != nil {
			return itemsList, nil
		}
		i.cache.SetNX(ctx, "GetItems", data, _defaultExpiration)
		return itemsList, nil
	}

	itemsList := make([]items.Item, 0)
	err = json.Unmarshal(val, &itemsList)
	if err != nil {
		return nil, err
	}
	return itemsList, err
}
func (i *ItemRepo) GetItem(ctx context.Context, id, campaignId int) (items.Item, error) {
	itemKey := fmt.Sprintf("GetItem-%d-%d", id, campaignId)
	val, err := i.cache.Get(ctx, itemKey).Bytes()

	if err != nil {
		it, err := i.ItemRepo.GetItem(ctx, id, campaignId)
		if err != nil {
			return items.Item{}, err
		}

		data, err := json.Marshal(it)
		if err != nil {
			return it, nil
		}
		i.cache.SetNX(ctx, itemKey, data, _defaultExpiration)
		return it, nil
	}

	it := items.Item{}
	err = json.Unmarshal(val, &it)
	if err != nil {
		return it, err
	}
	return it, err
}
func (i *ItemRepo) UpdateItem(ctx context.Context, name, description string, id, campaignId int) (items.Item, error) {
	i.deleteKey(ctx, id, campaignId)
	return i.ItemRepo.UpdateItem(ctx, name, description, id, campaignId)
}
func (i *ItemRepo) DeleteItem(ctx context.Context, id, campaignId int) (items.Item, error) {
	i.deleteKey(ctx, id, campaignId)
	return i.ItemRepo.DeleteItem(ctx, id, campaignId)
}
func (i *ItemRepo) deleteKey(ctx context.Context, id, campaignId int) {
	itemKey := fmt.Sprintf("GetItem-%d-%d", id, campaignId)
	i.cache.Del(ctx, itemKey)
	i.cache.Del(ctx, "GetItems")
}

func NewItemRepo(repo *pgdb.ItemRepo, client *redis.Client) *ItemRepo {
	return &ItemRepo{
		ItemRepo: repo,
		cache:    client,
	}
}
