package gameserver

import (
	"cmgameserver/msgHandler"
	"errors"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
)

type ServerMessageHandler struct {
	gameServer *GameServer
	pool       *HandlerPool
}

func (handler ServerMessageHandler) Init() {
	handler.pool.Register(10001, &msgHandler.LoginToGameServerHandler{GameServer: handler.gameServer})
	handler.pool.Register(10003,&msgHandler.RoleRegisterGateHandler{GameServer:handler.gameServer})
}
func (handler ServerMessageHandler) MessageReceived(session network.SocketSessionInterface, message interface{}) error {
	if writeMsg, ok := message.(network.WriteMessage); !ok {
		return errors.New("不是WriteMessage类型")
	} else {
		log4g.Infof("收到消息%d", writeMsg.MsgId)
		msgHandler := handler.pool.GetHandler(int32(writeMsg.MsgId))
		if msgHandler == nil {
			log4g.Errorf("不存在该消息[%d]的处理器", writeMsg.MsgId)
		} else {
			msgHandler.Action(session, writeMsg.MsgData)
		}
	}
	return nil
}

func (handler ServerMessageHandler) MessageSent(session network.SocketSessionInterface, message interface{}) error {
	return nil
}

func (handler ServerMessageHandler) SessionOpened(session network.SocketSessionInterface) error {
	return nil
}

func (handler ServerMessageHandler) SessionClosed(session network.SocketSessionInterface) {

}

func (handler ServerMessageHandler) SessionPeriod(session network.SocketSessionInterface) {
	log4g.Info("Period")
}

func (handler ServerMessageHandler) ExceptionCaught(session network.SocketSessionInterface, err error) {

}
