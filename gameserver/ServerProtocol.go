package gameserver
import (
	"github.com/bianchengxiaobei/cmgo/network"
	"bytes"
	"encoding/binary"
	"unsafe"
	"reflect"
	"github.com/golang/protobuf/proto"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ServerProtocol struct {
	pool 	*ProtoMessagePool
}
type MessageHeader struct {
	MessageId  	int32
	OrderId    	int32
	MsgBodyLen 	int32
}
var messageHeaderLen = (int32)(unsafe.Sizeof(MessageHeader{}))
func (protocol ServerProtocol) Init(){
	//注册消息
	protocol.pool.Register(10000,reflect.TypeOf(message.M2G_RegisterGate{}))
	protocol.pool.Register(10001,reflect.TypeOf(message.G2M_LoginToGameServer{}))
}
func (protocol ServerProtocol) Decode(session network.SocketSessionInterface, data []byte)(interface{},int,error){
	var (
		err 		error
		ioBuffer 	*bytes.Buffer
		msgHeader	MessageHeader
		chanMsg		network.WriteMessage
	)
	msgHeader = MessageHeader{}
	ioBuffer = bytes.NewBuffer(data)
	if int32(ioBuffer.Len()) < messageHeaderLen{
		return nil,0,nil
	}
	err = binary.Read(ioBuffer,binary.LittleEndian,&msgHeader)
	if err != nil{
		return nil,0,err
	}
	allLen := msgHeader.MsgBodyLen + messageHeaderLen
	if int32(ioBuffer.Len()) < allLen{
		return nil,0,nil
	}
	var perOrder = session.GetAttribute(network.PREORDERID)
	if perOrder == nil{
		session.SetAttribute(network.PREORDERID,msgHeader.OrderId+1)
		//if msgHeader.OrderId == 0{
		//	fmt.Println("用户客户端发送消息序列成功")
		//}
	}else{
		if msgHeader.OrderId == perOrder{
			session.SetAttribute(network.PREORDERID,msgHeader.OrderId+1)
		}else {
			log4g.Error("发送消息序列出错")
			return nil, 0, nil
		}
	}
	var msgType = protocol.pool.GetMessageType(msgHeader.MessageId)
	msg := reflect.New(msgType.Elem()).Interface()
	proto.Unmarshal(ioBuffer.Bytes(),msg.(proto.Message))
	chanMsg = network.WriteMessage{
		MsgId:int(msgHeader.MessageId),
		MsgData:msg,
	}
	return chanMsg,int(allLen),nil
}
func (protocol ServerProtocol) Encode(session network.SocketSessionInterface,writeMsg interface{}) error{
	var (
		err 			error
		ioBuffer 		*bytes.Buffer
		msgHeader	MessageHeader
		ok 			bool
		msg			network.WriteMessage
		protoMsg 	proto.Message
		data 		[]byte
	)
	msg,ok = writeMsg.(network.WriteMessage)
	if ok == false{
		panic("Message != WriteMsg")
	}
	msgHeader.MessageId = int32(msg.MsgId)

	msgHeader.OrderId = 0
	protoMsg,ok = msg.MsgData.(proto.Message)
	if ok == false{
		panic("Msg != ProtoMessage")
	}
	data,err = proto.Marshal(protoMsg)
	if err != nil{
		panic("ProtoMessage Marshal Error")
	}
	msgHeader.MsgBodyLen = int32(len(data))

	ioBuffer = &bytes.Buffer{}
	err = binary.Write(ioBuffer,binary.LittleEndian,msgHeader)
	if err != nil{
		log4g.Error(err.Error())
		return err
	}
	ioBuffer.Write(data)
	if err = session.WriteBytes(ioBuffer.Bytes());err != nil{
		log4g.Error("WriteBytes Error")
		return err
	}
	return nil
}