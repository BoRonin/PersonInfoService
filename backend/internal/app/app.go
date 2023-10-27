package app

import (
	"emtest/internal/config"
	"emtest/internal/fetcher"
	"emtest/internal/repository"
	"emtest/internal/service"
	"emtest/internal/storage/postgres"
	"emtest/internal/storage/redis"
	"emtest/pkg/logger"
	"fmt"
	"net/http"
	"os"

	"log/slog"
)

type App struct {
	Config config.Config
	Serv   http.Server
	PS     *service.PersonService
	Log    *slog.Logger
}

func NewApp(conf config.Config, log *slog.Logger) *App {
	//Инициализируем DB и Redis
	pg, err := postgres.NewDB(conf.DSN)
	if err != nil {
		log.Error("falied to init storage", logger.Err(err))
		os.Exit(1)
	}
	redis := redis.NewRedis(conf.RedisAddress)

	//Инициализируем сервисы
	r := repository.NewRedis(redis)
	bd := repository.NewPostgresService(pg)
	fetcher := fetcher.NewFetcher("https://api.agify.io/", "https://api.genderize.io/", "https://api.nationalize.io/")
	ps := service.NewPersonService(r, bd, fetcher, log)
	app := &App{
		Log:    log,
		Config: conf,
		Serv: http.Server{
			Addr: fmt.Sprintf(":%s", conf.Port),
		},
		PS: ps,
	}
	//Инициализируем пути
	app.Serv.Handler = app.NewRouter()
	return app
}
