package context

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"kcpServer/message"
	"kcpServer/message/pb"

	"github.com/xtaci/kcp-go"
	"google.golang.org/protobuf/proto"
)

const (
	PACK_MAX_LEN = 10240
	HEAD_LEN     = 6
)

type Head struct {
	Len     uint32
	MsgType uint16
}

type Packet struct {
	Head    *Head
	Msg     proto.Message
	Session *kcp.UDPSession
}

func UnpackMsg(data []byte) (pack *Packet, err error) {
	var head Head

	head.Len = uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	head.MsgType = uint16(data[4]) | uint16(data[5])<<8
	if head.Len > (PACK_MAX_LEN + HEAD_LEN) {
		return nil, errors.New(fmt.Sprintf("msg len error: %d", head.Len))
	}

	msg := message.GetMsgStruct(pb.MsgType(head.MsgType))
	if msg == nil {
		return nil, errors.New(fmt.Sprintf("no msg struct code %d", head.MsgType))
	}

	err = proto.Unmarshal(data[HEAD_LEN:head.Len], msg)
	if err != nil {
		return nil, err
	}

	packet := &Packet{
		Head: &head,
		Msg:  msg,
	}

	return packet, nil
}

func PackMsg(head *Head, message proto.Message) []byte {
	msg, _ := proto.Marshal(message)
	head.Len = HEAD_LEN + uint32(len(msg))

	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, head.Len)
	binary.Write(buff, binary.LittleEndian, head.MsgType)
	buff.Write(msg)

	return buff.Bytes()
}
