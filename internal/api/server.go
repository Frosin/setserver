package api

import (
	"context"
	"log"
	"time"

	"github.com/Frosin/setserver/internal/api/gen"
	"github.com/Frosin/setserver/internal/notification"
	"github.com/Frosin/setserver/internal/storage"
)

const (
	operationAdd    = "add"
	operationDelete = "delete"
)

type Config struct {
	Timeout time.Duration
}

type Server struct {
	gen.UnsafeApiServer
	storage     storage.Storage
	notificator notification.Notificator
	config      Config
}

func NewServer(config Config, storage storage.Storage, notificator notification.Notificator) Server {
	return Server{
		config:      config,
		storage:     storage,
		notificator: notificator,
	}
}

func (s *Server) sendNotifications(operation, key, value string) {
	if err := s.notificator.Broadcast(operation, key, value); err != nil {
		log.Printf("broadcast error: %s", err.Error())
	}
}

func (s *Server) Add(ctx context.Context, in *gen.Set) (*gen.Result, error) {
	tctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()
	if err := s.storage.Add(tctx, in.Name, in.Value); err != nil {
		return &gen.Result{Result: false}, err
	}
	//send notifications
	go func() {
		s.sendNotifications(operationAdd, in.Name, in.Value)
	}()
	return &gen.Result{Result: true}, nil
}

func (s *Server) Delete(ctx context.Context, in *gen.Set) (*gen.Result, error) {
	tctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()
	value, err := s.storage.Delete(tctx, in.Name)
	if err != nil {
		return &gen.Result{Result: false}, err
	}
	//send notifications
	go func() {
		s.sendNotifications(operationDelete, in.Name, value)
	}()
	return &gen.Result{Result: true}, nil
}

func (s *Server) Subscribe(in *gen.Empty, sv gen.Api_SubscribeServer) error {
	return s.notificator.RegisterConsumer(sv)
}
