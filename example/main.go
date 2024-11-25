package main

import (
	"log/slog"

	"kcpServer"
	"kcpServer/message/pb"
)

const (
	SERVER_ADDR = "127.0.0.1:8080"
	SERVER_PASS = "JKwfwA051x"
	SERVER_SALT = "test"
)

type Echo struct {
	*context.Context
}

func (e *Echo) Echo(packet *context.Packet, ctx *context.Context) error {
	req := packet.Msg.(*pb.EchoReq)
	if req == nil {
		return nil
	}

	slog.Debug("receive new echo msg: %v", req.GetData())

	res := &pb.EchoRes{Data: req.GetData()}
	ctx.Response(packet.Session, pb.MsgType_MSG_ECHO_RES, res)

	return nil
}

func main() {
	conf := &context.Conf{
		Addr: SERVER_ADDR,
		Pass: SERVER_PASS,
		Salt: SERVER_SALT,
	}

	ctx := context.NewContext(conf)
	ctx.Server()

	echo := &Echo{ctx}
	ctx.Register(pb.MsgType_MSG_ECHO_REQ, echo.Echo)

	ctx.Setup()
}
