package health

import (
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type App struct {
	pingers []func() error
	addr    string
	srv     *http.ServeMux
}

func New(pingers []func() error, addr, port string) *App {
	srv := http.NewServeMux()
	srv.HandleFunc("/healz", func(w http.ResponseWriter, r *http.Request) {
		errg := errgroup.Group{}
		for _, f := range pingers {
			errg.Go(f)
		}
		if errg.Wait() != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return &App{
		pingers: pingers,
		addr:    fmt.Sprint(addr, ":", port),
		srv:     srv,
	}
}

func (a *App) Start() {
	l := slog.Default()
	l.Debug("health starting")
	fmt.Println("health starting")
	go func() {
		if err := http.ListenAndServe(a.addr, a.srv); err != nil {
			l.Debug("err", slog.Any("err", err.Error()))
			fmt.Println(err.Error())
		}
	}()
	l.Debug("health started")
	fmt.Println("health started")
}
