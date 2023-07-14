package campaign

import (
	"context"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/campaign"
	"github.com/antoniokichaev/hezzl-collector/internal/repo"
)

type UseCase struct {
	repo repo.Campaign
}

func (u *UseCase) CreateCampaign(ctx context.Context, name string) (campaign.Campaign, error) {
	if name == "" {
		return campaign.Campaign{}, ErrIncorrectName
	}
	return u.CreateCampaign(ctx, name)
}

func (u *UseCase) UpdateCampaign(ctx context.Context, name string, campaignId int) (campaign.Campaign, error) {
	if name == "" {
		return campaign.Campaign{}, ErrIncorrectName
	}
	return u.UpdateCampaign(ctx, name, campaignId)
}

func (u *UseCase) DeleteCampaign(ctx context.Context, id int) (campaign.Campaign, error) {
	return u.repo.DeleteCampaign(ctx, id)
}

func (u *UseCase) GetCampaign(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UseCase) GetCampaigns(ctx context.Context) ([]campaign.Campaign, error) {
	return u.repo.GetCampaigns(ctx)
}

func NewCampaingUseCase(repo repo.Campaign) *UseCase {
	return &UseCase{repo: repo}
}
