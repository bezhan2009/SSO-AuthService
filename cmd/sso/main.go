package main

import (
	"SSO/internal/app"
	"SSO/internal/config"
	"SSO/internal/lib/logger"
	kafkaProducer "SSO/pkg/brokers/kafka"
	"SSO/pkg/db"
	"SSO/pkg/db/redis"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load(".env") // Два уровня вверх от ./cmd/sso
	if err != nil {
		err = godotenv.Load("example.env")
		if err != nil {
			panic(errors.New(fmt.Sprintf("error loading .env file. Error is %s", err)))
		}
	}

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.AppParams.Env)

	log.Info("Starting application")

	err = db.ConnectToDB(cfg)
	if err != nil {
		panic(err)
	}

	err = redis.InitializeRedis(cfg.RedisParams)
	if err != nil {
		panic(err)
	}

	err = kafkaProducer.CreateProducer(cfg.KafkaParams)
	if err != nil {
		panic(err)
	}

	application := app.New(log, cfg)

	go application.GRPCServer.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sing := <-stop

	log.Info("Stopping application: ", slog.String("signal", sing.String()))

	err = db.CloseDB(cfg.AppParams.DBSM)
	if err != nil {
		log.Error("Failed to close DB connection", slog.String("error", err.Error()))
	}

	application.GRPCServer.Stop()
	log.Info("Application stopped successfully")
}
