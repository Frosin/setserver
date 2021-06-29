package api

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Frosin/setserver/internal/api/gen"
	"github.com/Frosin/setserver/internal/notification"
	"github.com/Frosin/setserver/internal/storage"
	"google.golang.org/grpc"
)

const (
	operationAdd    = "add"
	operationDelete = "delete"
	waitingTimeout  = 1
)

type Config struct {
	Timeout time.Duration
	Port    string
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

func (s *Server) RunServer() {
	grpcServer := grpc.NewServer()
	gen.RegisterApiServer(grpcServer, s)

	l, err := net.Listen("tcp", ":"+s.config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}
	log.Printf("the server is running, config: %#v", s.config)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
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
		log.Printf("failed to add: %s", err.Error())
		return &gen.Result{Result: false}, nil
	}
	//send notifications
	go func() {
		s.sendNotifications(operationAdd, in.Name, in.Value)
	}()
	return &gen.Result{Result: true}, nil
}

func (s *Server) Delete(ctx context.Context, in *gen.Name) (*gen.Result, error) {
	tctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()
	value, err := s.storage.Delete(tctx, in.Name)
	if err != nil {
		log.Printf("failed to delete: %s", err.Error())
		return &gen.Result{Result: false}, nil
	}
	//send notifications
	go func() {
		s.sendNotifications(operationDelete, in.Name, value)
	}()
	return &gen.Result{Result: true}, nil
}

func (s *Server) Subscribe(in *gen.Empty, sv gen.Api_SubscribeServer) error {
	err := s.notificator.RegisterConsumer(sv)
	if err != nil {
		return err
	}
	for {
		time.Sleep(time.Second * waitingTimeout)
	}
}
