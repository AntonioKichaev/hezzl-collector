package items

import "time"

type Item struct {
	ID          int       `db:"id"`
	CampaignID  int       `db:"campaign_id"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}
