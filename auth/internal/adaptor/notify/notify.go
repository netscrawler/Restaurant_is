package notify

import (
	"context"
	"sync"
	"time"

	notifyclient "github.com/netscrawler/Restaurant_is/auth/internal/app/notifyclient"
	notify "github.com/netscrawler/RispProtos/proto/gen/go/v1/notify"
	"go.uber.org/zap"
)

type tasks struct {
	to  string
	msg string
}

type Notify struct {
	log   *zap.Logger
	ntf   *notifyclient.Client
	tasks chan tasks
	wg    sync.WaitGroup
}

func (n *Notify) worker(id int) {
	defer n.wg.Done()

	n.log.Info("notify worker started", zap.Int("worker_id", id))

	for task := range n.tasks {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		_, err := n.ntf.Notify.Send(ctx, &notify.SendRequest{
			Phone: task.to,
			Data:  task.msg,
		})

		cancel()

		if err != nil {
			n.log.Warn("failed to send notification",
				zap.Int("worker_id", id),
				zap.String("phone", task.to),
				zap.Error(err),
			)
		}
	}

	n.log.Info("notify worker stopped", zap.Int("worker_id", id))
}

func New(log *zap.Logger, notify *notifyclient.Client) *Notify {
	n := &Notify{
		log:   log,
		ntf:   notify,
		tasks: make(chan tasks, 1000),
		wg:    sync.WaitGroup{},
	}

	for i := range 10 {
		n.wg.Add(1)
		go n.worker(i)
	}

	return n
}

func (n *Notify) Send(ctx context.Context, to string, msg string) {
	select {
	case n.tasks <- tasks{to: to, msg: msg}:
	default:
		n.log.Error("notify queue is full", zap.String("phone", to))
	}
}

func (n *Notify) Shutdown() {
	close(n.tasks)
	n.wg.Wait()
}
