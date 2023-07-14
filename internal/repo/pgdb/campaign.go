package pgdb

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	campaignEntity "github.com/antoniokichaev/hezzl-collector/internal/entity/campaign"
	"github.com/antoniokichaev/hezzl-collector/pkg/postgres"
)

var ErrNotFound = errors.New("campaign.not.found")

type CampaignRepo struct {
	*postgres.Postgres
}

func (c *CampaignRepo) CreateCampaign(ctx context.Context, name string) (campaignEntity.Campaign, error) {
	const fName = "CreateCampaign"
	sqlReq, args, err := c.Builder.Insert(_tableCampaigns).Columns("name").Values(name).Suffix("returning id").ToSql()
	if err != nil {
		return campaignEntity.Campaign{}, fmt.Errorf("%s builder %w", fName, err)
	}
	var id int
	err = c.Pool.QueryRow(ctx, sqlReq, args...).Scan(&id)
	if err != nil {
		return campaignEntity.Campaign{}, fmt.Errorf("%s QueryRow %w", fName, err)
	}
	return campaignEntity.Campaign{ID: id, Name: name}, err
}

func (c *CampaignRepo) UpdateCampaign(ctx context.Context, name string, id int) (cmp campaignEntity.Campaign, err error) {
	const fName = "UpdateCampaign"
	cmp.Name = name
	cmp.ID = id
	tx, err := c.Pool.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("%s Begin %w", fName, err)
		return
	}
	defer func() {
		if err != nil {
			if errRb := tx.Rollback(ctx); errRb != nil {
				err = fmt.Errorf("%s Rollback %w", fName, err)
			}
		}
	}()

	sqlReq, args, _ := c.Builder.
		Update(_tableCampaigns).
		Set("name", name).
		Where(sq.Eq{"id": id}).
		ToSql()

	rowsAffected, err := tx.Exec(ctx, sqlReq, args...)
	if err != nil {
		err = fmt.Errorf("%s Exec %w", fName, err)
		return
	}
	if rowsAffected.RowsAffected() == 0 {
		err = ErrNotFound
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("%s Commit %w", fName, err)
		return
	}
	return

}

func (c *CampaignRepo) DeleteCampaign(ctx context.Context, id int) (campaignEntity.Campaign, error) {
	const fName = "DeleteCampaign"
	sqlReq, args, _ := c.Builder.
		Delete(_tableCampaigns).
		Where(sq.Eq{"id": id}).
		Suffix("returning name").ToSql()
	var name string
	err := c.Pool.QueryRow(ctx, sqlReq, args...).Scan(&name)
	if err != nil {
		return campaignEntity.Campaign{}, fmt.Errorf("%s %w", fName, ErrNotFound)
	}

	return campaignEntity.Campaign{Name: name, ID: id}, nil
}

func (c *CampaignRepo) GetCampaign(ctx context.Context, id int) (campaignEntity.Campaign, error) {
	const fName = "GetCampaign"
	camp := campaignEntity.Campaign{}
	sqlReq, args, _ := c.Builder.Select("*").From(_tableCampaigns).Where(sq.Eq{"id": id}).ToSql()
	err := c.Pool.QueryRow(ctx, sqlReq, args...).Scan(&camp.ID, &camp.Name)
	if err != nil {
		return campaignEntity.Campaign{}, fmt.Errorf("%s QueryRow %w", fName, err)
	}

	return camp, nil
}

func (c *CampaignRepo) GetCampaigns(ctx context.Context) ([]campaignEntity.Campaign, error) {
	const fName = "GetCampaigns"
	camps := make([]campaignEntity.Campaign, 0)
	sqlReq, args, _ := c.Builder.Select("*").From(_tableCampaigns).ToSql()
	rows, err := c.Pool.Query(ctx, sqlReq, args...)
	if err != nil {
		return nil, fmt.Errorf("%s Query %w", fName, err)
	}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("%s rows.Scan %w", fName, err)
		}
		camps = append(camps, campaignEntity.Campaign{ID: id, Name: name})
	}

	return camps, nil
}

func NewCampaignRepo(db *postgres.Postgres) *CampaignRepo {
	return &CampaignRepo{
		Postgres: db,
	}
}
