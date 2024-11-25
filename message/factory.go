package message

import (
	"kcpServer/message/pb"

	"google.golang.org/protobuf/proto"
)

func GetMsgStruct(msgCode pb.MsgType) proto.Message {
	switch msgCode {
	case pb.MsgType_MSG_ECHO_REQ:
		return &pb.EchoReq{}
	case pb.MsgType_MSG_ECHO_RES:
		return &pb.EchoRes{}
	}
	return nil
}
