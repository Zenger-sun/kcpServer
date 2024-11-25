package context

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Context struct {
	*Conf
	*Router
}

func (ctx *Context) Server() {

	go listen(ctx)

	slog.Info("server started!")
}

func (ctx *Context) Setup() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		switch <-exit {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			ctx.Shutdown()
			os.Exit(0)
		default:
			os.Exit(1)
		}
	}
}

func (ctx *Context) Shutdown() {
	slog.Info("Shutdown...")
	slog.Info("Closed!")
}

func NewContext(conf *Conf) *Context {
	sync := &Context{
		Conf:   conf,
		Router: NewRouter(),
	}
	return sync
}
