package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	ConnTimeout  time.Duration
	Pool         *pgxpool.Pool
	Builder      sq.StatementBuilderType
}

func New(dataSource string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		ConnTimeout:  _defaultConnTimeout,
	}

	pg.Builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if dataSource == "" {
		return pg, nil
	}
	for _, opt := range opts {
		opt(pg)
	}

	pgxConfig, err := pgxpool.ParseConfig(dataSource)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}
	pgxConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), pgxConfig)
		if err == nil {
			break
		}

		pg.connAttempts--
		time.Sleep(pg.ConnTimeout)
	}

	if err != nil {
		return nil, fmt.Errorf("pgdb - New - pgxpool.NewWithConfig: %w", err)
	}

	return pg, nil

}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
