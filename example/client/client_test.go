package client

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"kcpServer"
	"kcpServer/message/pb"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

const (
	SERVER_ADDR = "127.0.0.1:8081"
	SERVER_PASS = "JKwfwA051x"
	SERVER_SALT = "test"
)

func TestClinet_Echo(t *testing.T) {
	key := pbkdf2.Key([]byte(SERVER_PASS), []byte(SERVER_SALT), 1024, 32, sha256.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	sess, err := kcp.DialWithOptions(SERVER_ADDR, block, 10, 3)
	if err != nil {
		return
	}
	defer sess.Close()

	head := &context.Head{MsgType: uint16(pb.MsgType_MSG_ECHO_REQ)}
	msg := &pb.EchoReq{Data: "test echo content"}

	_, err = sess.Write(context.PackMsg(head, msg))
	if err != nil {
		return
	}

	pack := make([]byte, 1024)
	for {
		_, err := sess.Read(pack)
		if err != nil {
			continue
		}

		packet, err := context.UnpackMsg(pack)
		if err != nil {
			return
		}

		res := packet.Msg.(*pb.EchoRes)
		if res == nil {
			return
		}

		fmt.Println(head.MsgType, res.Data)

		break
	}
}

func BenchmarkClient(t *testing.B) {
	key := pbkdf2.Key([]byte(SERVER_PASS), []byte("test"), 1024, 32, sha256.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	sess, err := kcp.DialWithOptions(SERVER_ADDR, block, 10, 3)
	if err != nil {
		return
	}

	head := &context.Head{MsgType: uint16(pb.MsgType_MSG_ECHO_REQ)}
	msg := &pb.EchoReq{Data: "test echo content"}

	pack := make([]byte, 1024)
	for i := 0; i < t.N; i++ {

		sess.Write(context.PackMsg(head, msg))

		for {
			_, err := sess.Read(pack)
			if err != nil {
				continue
			}

			packet, err := context.UnpackMsg(pack)
			if err != nil {
				return
			}

			res := packet.Msg.(*pb.EchoRes)
			if res == nil {
				return
			}

			break
		}
	}
}