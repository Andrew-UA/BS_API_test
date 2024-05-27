package app

import (
	"context"
	"github.com/Andrew-UA/BS_API_test/internal/api/http"
	"github.com/Andrew-UA/BS_API_test/internal/server"
	"github.com/Andrew-UA/BS_API_test/internal/service"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func Run() {

	deviceService := &service.DeviceServiceMock{}
	taskManagerService := service.NewTaskManagerService(deviceService)
	taskHandler := http.NewTaskHandler(taskManagerService)
	router := http.NewRouter(taskHandler)
	srv := server.NewServer(router)

	// Create dummy devices for mock device service.
	createDevices(taskManagerService)

	go func() {
		if err := srv.Run(); err != nil {
			slog.Error("SERVER ERROR ", "err", err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		slog.Error("Failed to stop HTTP server: %s", "err", err.Error())
	}
}

func createDevices(tmManager *service.TaskManagerService) {

	for i := 0; i <= 4; i++ {
		_ = tmManager.AddDevice("localhost:600"+strconv.Itoa(i), "login", "password")
	}
}
