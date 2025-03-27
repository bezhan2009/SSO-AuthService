package grpcapp

import (
	"errors"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strconv"

	authgrpc "SSO/internal/grpc/auth"
	pinggrpc "SSO/internal/grpc/ping"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.Auth, pingService pinggrpc.PingAPI, port int) *App {
	grpcServer := grpc.NewServer()

	authgrpc.Register(grpcServer, authService)
	pinggrpc.Register(grpcServer, pingService)

	return &App{
		log,
		grpcServer,
		port,
	}
}

func (a *App) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Start() error {
	const op = "grpc.Start"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", strconv.Itoa(a.port)))

	l, err := net.Listen("tcp", ":"+strconv.Itoa(a.port))
	if err != nil {
		return errors.New(op + ": " + err.Error())
	}

	log.Info("starting gRPC server on address: " + l.Addr().String())

	if err := a.gRPCServer.Serve(l); err != nil {
		return errors.New(op + ": " + err.Error())
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpc.Stop"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", strconv.Itoa(a.port)))

	log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
