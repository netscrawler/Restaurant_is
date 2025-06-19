package notify

import (
	"context"
	"log/slog"
	"sync"
	"time"

	notifyclient "github.com/netscrawler/Restaurant_is/auth/internal/app/notifyclient"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
	notify "github.com/netscrawler/RispProtos/proto/gen/go/v1/notify"
)

type tasks struct {
	ctx context.Context
	to  string
	msg string
}

type Notify struct {
	log   *slog.Logger
	ntf   *notifyclient.Client
	tasks chan tasks
	wg    sync.WaitGroup
}

func (n *Notify) worker(id int) {
	defer n.wg.Done()

	log := utils.LoggerWithTrace(context.Background(), n.log)
	log.Info("notify worker started", slog.Int("worker_id", id))

	for task := range n.tasks {
		log := utils.LoggerWithTrace(task.ctx, n.log)
		ctx, cancel := context.WithTimeout(task.ctx, 2*time.Second)

		_, err := n.ntf.Notify.Send(ctx, &notify.SendRequest{
			Phone: task.to,
			Data:  task.msg,
		})

		cancel()

		if err != nil {
			log.Warn("failed to send notification",
				slog.Int("worker_id", id),
				slog.String("phone", task.to),
				slog.Any("error", err),
			)
		}
	}

	log.Info("notify worker stopped", slog.Int("worker_id", id))
}

func New(log *slog.Logger, notify *notifyclient.Client) *Notify {
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
	case n.tasks <- tasks{ctx: ctx, to: to, msg: msg}:
	default:
		n.log.Error("notify queue is full", slog.String("phone", to))
	}
}

func (n *Notify) Shutdown() {
	close(n.tasks)
	n.wg.Wait()
}
