package main

import (
	"log"

	"gihub.com/gibiw/api-example/internal/config"
	"gihub.com/gibiw/api-example/internal/repository"
	"gihub.com/gibiw/api-example/internal/transport/httpserver"
	"gihub.com/gibiw/api-example/internal/usecases"
	"gihub.com/gibiw/api-example/pkg/database"
	"github.com/gibiw/cache"
	"github.com/gookit/slog"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	var cfg config.Config
	err := cleanenv.ReadConfig("config/config.yml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.SetFormatter(slog.NewJSONFormatter())
	slog.SetLogLevel(slog.LevelByName(cfg.LoggerCfg.Level))

	db, err := database.Initialize(cfg.DBCfg)
	if err != nil {
		slog.Fatal("can not initialize database", err)
	}

	repo := repository.New(db)

	cache := cache.New()
	ucs := usecases.New(repo)
	srv := httpserver.New(cfg.ServiceCfg, ucs, cache)

	err = srv.Run()
	if err != nil {
		slog.Fatal("can not start server", err)
	}
}
