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

// @title           Cars API
// @version         1.0
// @description     This is a example API service.

// @contact.name   Dmitry Gridnev
// @contact.url    https://github.com/gibiw/api-example
// @contact.email  gibiw1@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

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
