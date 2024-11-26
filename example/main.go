package main

import (
	"log/slog"

	"kcpServer"
	"kcpServer/message/pb"

	"github.com/asynkron/protoactor-go/actor"
)

const (
	SERVER_ADDR = "127.0.0.1:8081"
	SERVER_PASS = "JKwfwA051x"
	SERVER_SALT = "test"
)

type Echo struct {
	ctx *context.Context
	pid *actor.PID
}

func (e *Echo) Receive(actorCtx actor.Context) {
	switch msg := actorCtx.Message().(type) {
	case *actor.Started:
		slog.Info("service echo started!")

	case *context.Packet:
		switch pb.MsgType(msg.Head.MsgType) {
		case pb.MsgType_MSG_ECHO_REQ:
			e.handleEcho(msg)
		}
	}
}

func (e *Echo) handleEcho(msg *context.Packet) {
	req := msg.Msg.(*pb.EchoReq)
	if req == nil {
		return
	}

	slog.Info("receive new echo", "req", req.GetData())

	res := &pb.EchoRes{Data: req.GetData()}
	e.ctx.Response(msg.Session, pb.MsgType_MSG_ECHO_RES, res)

	return
}

func (e *Echo) Pid() *actor.PID {
	return e.pid
}

func NewEcho(ctx *context.Context) *Echo {
	echo := &Echo{ctx: ctx}
	echo.pid = ctx.Spawn(actor.PropsFromProducer(func() actor.Actor { return echo }))
	return echo
}

func main() {
	conf := &context.Conf{
		Addr: SERVER_ADDR,
		Pass: SERVER_PASS,
		Salt: SERVER_SALT,
	}

	ctx := context.NewContext(conf)
	ctx.Server()

	echo := NewEcho(ctx)
	ctx.Send(ctx.Pid(), &context.SvcRegister{Svc: echo, Msg: pb.MsgType_MSG_ECHO_REQ})

	ctx.Setup()
}
