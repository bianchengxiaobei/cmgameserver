package gameserver
import (
	"github.com/bianchengxiaobei/cmgo/network"
)

type ServerMessageHandler struct {
	gameServer *GameServer
}

func (handler ServerMessageHandler) MessageReceived(session network.SocketSessionInterface, message interface{}) error {
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

}

func (handler ServerMessageHandler) ExceptionCaught(session network.SocketSessionInterface, err error) {

}