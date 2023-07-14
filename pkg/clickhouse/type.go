package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	sq "github.com/Masterminds/squirrel"
)

type Clickhouse struct {
	DB       driver.Conn
	database string
	username string
	password string
	Builder  sq.StatementBuilderType
}

func New(host string, port int, opts ...Option) (*Clickhouse, error) {
	cl := &Clickhouse{Builder: sq.StatementBuilder.PlaceholderFormat(sq.Question)}

	for _, opt := range opts {
		opt(cl)
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: clickhouse.Auth{
			Database: cl.database,
			Username: cl.username,
			Password: cl.password,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("init clickhouse err %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	cl.DB = conn

	return cl, nil
}

func (c *Clickhouse) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}
