package redis

import (
	"context"
	"encoding/json"
	"fmt"
	campaignEntity "github.com/antoniokichaev/hezzl-collector/internal/entity/campaign"
	"github.com/antoniokichaev/hezzl-collector/internal/repo/pgdb"
	"github.com/redis/go-redis/v9"
)

type CampaignRepo struct {
	*pgdb.CampaignRepo
	cache *redis.Client
}

func (c *CampaignRepo) GetCampaigns(ctx context.Context) ([]campaignEntity.Campaign, error) {
	redisKey := "GetCampaigns"
	data, err := c.cache.Get(ctx, redisKey).Bytes()
	if err != nil {
		campList, err := c.CampaignRepo.GetCampaigns(ctx)
		if err != nil {
			return nil, err
		}

		data, err = json.Marshal(campList)
		if err != nil {
			return campList, err
		}

		c.cache.SetNX(ctx, redisKey, data, _defaultExpiration)
		return campList, nil
	}

	campList := make([]campaignEntity.Campaign, 0)
	err = json.Unmarshal(data, &campList)
	if err != nil {
		return nil, err
	}

	return campList, nil
}
func (c *CampaignRepo) GetCampaign(ctx context.Context, id int) (campaignEntity.Campaign, error) {
	redisKey := fmt.Sprintf("GetCampaign-%d", id)
	data, err := c.cache.Get(ctx, redisKey).Bytes()
	if err != nil {
		camp, err := c.CampaignRepo.GetCampaign(ctx, id)

		if err != nil {
			return campaignEntity.Campaign{}, err
		}
		data, err := json.Marshal(camp)
		if err != nil {
			return campaignEntity.Campaign{}, err
		}
		c.cache.SetNX(ctx, redisKey, data, _defaultExpiration)
		return camp, nil
	}
	camp := campaignEntity.Campaign{}
	err = json.Unmarshal(data, &camp)
	if err != nil {
		return campaignEntity.Campaign{}, err
	}
	return camp, nil
}
func (c *CampaignRepo) UpdateCampaign(ctx context.Context, name string, id int) (campaignEntity.Campaign, error) {
	c.deleteKey(ctx, id)
	return c.CampaignRepo.UpdateCampaign(ctx, name, id)
}
func (c *CampaignRepo) DeleteCampaign(ctx context.Context, id int) (campaignEntity.Campaign, error) {
	c.deleteKey(ctx, id)
	return c.CampaignRepo.DeleteCampaign(ctx, id)
}
func (i *CampaignRepo) deleteKey(ctx context.Context, id int) {
	redisKey := fmt.Sprintf("GetCampaign-%d", id)
	i.cache.Del(ctx, redisKey)
	i.cache.Del(ctx, "GetCampaigns")
}

func NewCampaignRepo(repo *pgdb.CampaignRepo, client *redis.Client) *CampaignRepo {
	return &CampaignRepo{
		CampaignRepo: repo,
		cache:        client,
	}
}
