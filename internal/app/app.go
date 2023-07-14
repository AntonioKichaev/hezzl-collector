package app

import (
	"context"
	"github.com/antoniokichaev/hezzl-collector/config"
	v1 "github.com/antoniokichaev/hezzl-collector/internal/controller/http/v1"
	"github.com/antoniokichaev/hezzl-collector/internal/repo"
	"github.com/antoniokichaev/hezzl-collector/internal/service"
	"github.com/antoniokichaev/hezzl-collector/pkg/clickhouse"
	"github.com/antoniokichaev/hezzl-collector/pkg/httpserver"
	"github.com/antoniokichaev/hezzl-collector/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	"syscall"
)

func Run(configPath string) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	//init logger

	//init db
	pg, err := postgres.New(
		cfg.Postgres.URL,
		postgres.ConnAttempts(cfg.Postgres.ConnAttempts),
		postgres.MaxPoolSize(cfg.Postgres.MaxPoolSize),
	)
	if err != nil {
		panic(err)
	}
	defer pg.Close()

	//init clickhouse
	clDb, err := clickhouse.New(
		cfg.Clickhouse.Addr,
		cfg.Clickhouse.NativePort,
		clickhouse.WithDatabase(cfg.Clickhouse.DB),
		clickhouse.WithUserName(cfg.Clickhouse.Username),
		clickhouse.WithPassword(cfg.Clickhouse.Password),
	)
	if err != nil {
		panic(err)
	}
	_ = clDb.DB.Ping(ctx)

	//init natsClient
	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		panic(err)
	}
	defer func() { _ = nc.Drain() }()
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	//init RedisClient
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	//init services
	repos := repo.NewRepositories(pg, clDb, redisClient, js)

	sd := service.SDependencies{Repos: repos}
	services := service.NewServices(&sd, js)
	go services.EventSaver.Start(ctx)

	ginEngine := gin.New()

	// register routers
	v1.New(ginEngine, services)
	httpServer := httpserver.New(ginEngine)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		//log

	case <-httpServer.Notify():
		//log
	}
	cancel()

	//shutdown
	_ = httpServer.Shutdown()

}
