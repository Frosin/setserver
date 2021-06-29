package notification

import "github.com/Frosin/setserver/internal/api/gen"

type Notificator interface {
	RegisterConsumer(gen.Api_SubscribeServer) error
	Broadcast(operation, key, value string) error
}
