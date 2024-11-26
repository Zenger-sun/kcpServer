package context

import (
	"kcpServer/message/pb"

	"github.com/asynkron/protoactor-go/actor"
)

type Service interface {
	Receive(ctx actor.Context)
	Pid() *actor.PID
}

type SvcRegister struct {
	Svc Service
	Msg pb.MsgType
}

type Router struct {
	*actor.RootContext
	pid    *actor.PID
	svcMap map[pb.MsgType]Service
}

func (r *Router) Receive(actorCtx actor.Context) {
	switch msg := actorCtx.Message().(type) {
	case *SvcRegister:
		r.svcMap[msg.Msg] = msg.Svc

	case *Packet:
		if svc, ok := r.svcMap[pb.MsgType(msg.Head.MsgType)]; ok {
			r.Send(svc.Pid(), msg)
		}
	}
}

func (r *Router) Pid() *actor.PID {
	return r.pid
}

func NewRouter() *Router {
	router := &Router{svcMap: make(map[pb.MsgType]Service)}
	router.RootContext = actor.NewActorSystem().Root
	router.pid = router.RootContext.Spawn(actor.PropsFromProducer(func() actor.Actor { return router }))
	return router
}
