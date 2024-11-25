package context

import (
	"crypto/sha256"
	"log/slog"
	"os"

	"kcpServer/message/pb"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/protobuf/proto"
)

func listen(ctx *Context) error {
	key := pbkdf2.Key([]byte(ctx.Pass), []byte(ctx.Salt), 1024, 32, sha256.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	listener, err := kcp.ListenWithOptions(ctx.Addr, block, 10, 3)
	if err != nil {
		os.Exit(1)
	}

	for {
		session, err := listener.AcceptKCP()
		if err != nil {
			slog.Error("Failed to accept KCP connection: %v", err)
			continue
		}

		session.SetNoDelay(1, 10, 2, 1)
		session.SetStreamMode(true)
		session.SetWindowSize(4096, 4096)
		session.SetReadBuffer(4 * 1024 * 1024)
		session.SetWriteBuffer(4 * 1024 * 1024)
		session.SetACKNoDelay(true)

		go readKcp(session, ctx)
	}
}

func readKcp(session *kcp.UDPSession, ctx *Context) {
	buffer := make([]byte, 1024)
	defer session.Close()

	for {
		_, err := session.Read(buffer)
		switch err {
		case nil:
		default:
			return
		}

		pack, err := UnpackMsg(buffer)
		if err != nil {
			slog.Warn("readKCP unpack msg err: ", err)
			continue
		}
		pack.Session = session

		handler, ok := ctx.Router.Handler[pb.MsgType(pack.Head.MsgType)]
		if !ok {
			continue
		}

		err = handler(pack, ctx)
		if err != nil {
			slog.Warn("handler err: ", err)
			continue
		}
	}
}

func (c *Context) Response(session *kcp.UDPSession, res pb.MsgType, msg proto.Message) {
	h := &Head{
		Len:     0,
		MsgType: uint16(res),
	}

	session.Write(PackMsg(h, msg))
}
