package item

import (
	"context"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/items"
	"github.com/antoniokichaev/hezzl-collector/internal/repo"
)

type UseCase struct {
	repo repo.Item
}

func (i *UseCase) CreateItem(ctx context.Context, name string, campaignId int) (items.Item, error) {
	if name == "" {
		return items.Item{}, ErrIncorrectName
	}
	return i.repo.CreateItem(ctx, name, campaignId)
}

func (i *UseCase) UpdateItem(ctx context.Context, name, description string, id, campaignId int) (items.Item, error) {
	if name == "" {
		return items.Item{}, ErrIncorrectName
	}
	return i.repo.UpdateItem(ctx, name, description, id, campaignId)
}

func (i *UseCase) DeleteItem(ctx context.Context, id, campaignId int) (items.Item, error) {
	return i.repo.DeleteItem(ctx, id, campaignId)
}

func (i *UseCase) GetItem(ctx context.Context, id, campaignId int) (items.Item, error) {
	return i.repo.GetItem(ctx, id, campaignId)
}

func (i *UseCase) GetItems(ctx context.Context) ([]items.Item, error) {
	return i.repo.GetItems(ctx)
}

func NewItemUseCase(repo repo.Item) *UseCase {
	return &UseCase{repo: repo}
}
