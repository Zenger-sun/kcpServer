package context

import (
	"crypto/sha256"
	"kcpServer/message/pb"
	"log/slog"
	"os"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/protobuf/proto"
)

const (
	BUFFER_LEN = 4 * 1024 * 1024
)

func listen(ctx *Context) error {
	key := pbkdf2.Key([]byte(ctx.Pass), []byte(ctx.Salt), 1024, 32, sha256.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	listener, err := kcp.ListenWithOptions(ctx.Addr, block, 10, 3)
	if err != nil {
		slog.Error("Failed to listen addr[%v], err: %v", ctx.Addr, err)
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
		session.SetReadBuffer(BUFFER_LEN)
		session.SetWriteBuffer(BUFFER_LEN)
		session.SetACKNoDelay(true)

		go readKcp(session, ctx)
	}
}

func readKcp(session *kcp.UDPSession, ctx *Context) {
	buffer := make([]byte, BUFFER_LEN)
	defer session.Close()

	for {
		_, err := session.Read(buffer)
		if err != nil {
			return
		}

		pack, err := UnpackMsg(buffer)
		if err != nil {
			slog.Warn("readKCP unpack msg err: ", err)
			continue
		}

		pack.Session = session
		ctx.Send(ctx.Pid(), pack)
	}
}

func (c *Context) Response(session *kcp.UDPSession, res pb.MsgType, msg proto.Message) {
	h := &Head{
		Len:     0,
		MsgType: uint16(res),
	}

	session.Write(PackMsg(h, msg))
}
