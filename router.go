package context

import (
	"kcpServer/message/pb"
)

type Handler func(packet *Packet, ctx *Context) error

type Router struct {
	Handler map[pb.MsgType]Handler
}

func (r *Router) Register(msgType pb.MsgType, handler Handler) {
	r.Handler[msgType] = handler
}

func NewRouter() *Router {
	return &Router{
		Handler: make(map[pb.MsgType]Handler),
	}
}
