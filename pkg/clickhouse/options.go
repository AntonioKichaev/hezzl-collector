package clickhouse

type Option func(postgres *Clickhouse)

func WithDatabase(name string) Option {
	return func(cl *Clickhouse) {
		cl.database = name
	}
}

func WithUserName(username string) Option {
	return func(cl *Clickhouse) {
		cl.username = username
	}
}

func WithPassword(password string) Option {
	return func(cl *Clickhouse) {
		cl.password = password
	}
}
