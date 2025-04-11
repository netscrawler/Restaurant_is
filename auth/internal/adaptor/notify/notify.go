package notify

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Notify struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Notify {
	return &Notify{
		log: log,
	}
}

func (n *Notify) Send(ctx context.Context, to string, msg string) {
	n.log.Info(fmt.Sprintf("To %s, msg: %s", to, msg))
}
