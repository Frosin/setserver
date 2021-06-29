package main

import (
	"os"
	"strconv"
	"time"

	"github.com/Frosin/setserver/internal/api"
	"github.com/Frosin/setserver/internal/notification"
	"github.com/Frosin/setserver/internal/storage"
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
	port := os.Getenv("SERVER_PORT")

	return api.Config{
		Timeout: time.Duration(timeout) * time.Second,
		Port:    port,
	}
}

func main() {
	setServer := api.NewServer(
		getConfig(),
		storage.NewMapStorage(),
		notification.NewBasicNotificator(),
	)
	setServer.RunServer()
}
