package campaign

type Campaign struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
