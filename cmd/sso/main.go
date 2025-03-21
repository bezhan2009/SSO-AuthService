package main

import (
	"SSO/internal/app"
	"SSO/internal/config"
	"SSO/internal/lib/logger"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// Создаем команду ls
	cmd := exec.Command("ls")

	// Получаем стандартный вывод команды
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Преобразуем байтовый срез в строку и выводим на консоль
	fmt.Println(strings.TrimSpace(string(output)))

	err = godotenv.Load(".env") // Два уровня вверх от ./cmd/sso
	if err != nil {
		err = godotenv.Load("example.env")
		if err != nil {
			panic(errors.New(fmt.Sprintf("error loading .env file. Error is %s", err)))
		}
	}

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting application")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCServer.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sing := <-stop

	log.Info("stopping application: ", slog.String("signal", sing.String()))

	application.GRPCServer.Stop()
	log.Info("Application stopped successfully")
}
