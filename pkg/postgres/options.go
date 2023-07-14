package postgres

import "time"

type Option func(postgres *Postgres)

func MaxPoolSize(maxSize int) Option {
	return func(postgres *Postgres) {
		postgres.maxPoolSize = maxSize
	}
}

func ConnAttempts(attempts int) Option {
	return func(postgres *Postgres) {
		postgres.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(postgres *Postgres) {
		postgres.ConnTimeout = timeout
	}
}
