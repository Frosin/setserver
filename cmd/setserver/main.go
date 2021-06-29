package main

import (
	"os"
	"strconv"
	"time"

	"github.com/Frosin/setserver/internal/api"
	"github.com/Frosin/setserver/internal/api/gen"
	"github.com/Frosin/setserver/internal/notification"
	"github.com/Frosin/setserver/internal/storage"

	"google.golang.org/grpc"
)

const (
	defaultTimeout = time.Second * 3
)

func getConfig() api.Config {
	timeoutStr := os.Getenv("STORAGE_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = int(defaultTimeout)
	}
	return api.Config{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func main() {
	grpcServer := grpc.NewServer()
	setServer := api.NewServer(
		getConfig(),
		storage.NewMapStorage(),
		notification.NewBasicNotificator(),
	)
	gen.RegisterApiServer(grpcServer, &setServer)
}
