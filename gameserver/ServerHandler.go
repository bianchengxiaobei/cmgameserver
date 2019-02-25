package gameserver

import (
	"cmgameserver/msgHandler"
	"errors"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"time"
)

type ServerMessageHandler struct {
	gameServer *GameServer
	pool       *HandlerPool
	pingMsg 	*message.M2C2M_GamePing
}

func (handler ServerMessageHandler) Init() {
	handler.pool.Register(10001, &msgHandler.LoginToGameServerHandler{GameServer: handler.gameServer})
	handler.pool.Register(10003,&msgHandler.RoleRegisterGateHandler{GameServer:handler.gameServer})
	handler.pool.Register(10004,&msgHandler.RoleQuitHandler{GameServer:handler.gameServer})

	handler.pool.Register(5001,&msgHandler.ReqRefreshRoomListHandle{GameServer:handler.gameServer})
	handler.pool.Register(5003,&msgHandler.CreateRoomHandler{GameServer:handler.gameServer})
	handler.pool.Register(5005,&msgHandler.ReqJoinRoomHandler{GameServer:handler.gameServer})
	handler.pool.Register(5006,&msgHandler.ReadyHandler{GameServer:handler.gameServer})
	handler.pool.Register(5008,&msgHandler.ReqStartBattleHandler{GameServer:handler.gameServer})
	handler.pool.Register(5010,&msgHandler.BattleLoadFinishedHandler{GameServer:handler.gameServer})
	handler.pool.Register(5013,&msgHandler.ReqCommandHandler{GameServer:handler.gameServer})
	handler.pool.Register(5014,&msgHandler.GamePingHandler{GameServer:handler.gameServer})
	handler.pool.Register(5017,&msgHandler.WinBattleHandler{GameServer:handler.gameServer})
	handler.pool.Register(5018,&msgHandler.FailedBattleHandler{GameServer:handler.gameServer})
	handler.pool.Register(5020,&msgHandler.WatchAdsHandler{GameServer:handler.gameServer})
	handler.pool.Register(5022,&msgHandler.ChangeNickNameHandler{GameServer:handler.gameServer})
	handler.pool.Register(5023,&msgHandler.ChangeAvatarIdHandler{GameServer:handler.gameServer})
	handler.pool.Register(5026,&msgHandler.ChangeEquipItemPosHandler{GameServer:handler.gameServer})
	handler.pool.Register(5027,&msgHandler.RoleQuitRoomHandler{GameServer:handler.gameServer})
	handler.pool.Register(5030,&msgHandler.BuyHeroHandler{GameServer:handler.gameServer})
	handler.pool.Register(5032,&msgHandler.RoomChatHandler{GameServer:handler.gameServer})
	handler.pool.Register(5033,&msgHandler.LearnSkillHandler{GameServer:handler.gameServer})
	handler.pool.Register(5035,&msgHandler.ChangeSkillHandler{GameServer:handler.gameServer})
	handler.pool.Register(5036,&msgHandler.GetAchieveHandler{GameServer:handler.gameServer})
	handler.pool.Register(5037,&msgHandler.GetTaskAwardHandler{GameServer:handler.gameServer})
	handler.pool.Register(5038,&msgHandler.GetSingHandler{GameServer:handler.gameServer})
	handler.pool.Register(5039,&msgHandler.ChangeFreeSoldierDataHandler{GameServer:handler.gameServer})
	handler.pool.Register(5040,&msgHandler.ChangeSexHandler{GameServer:handler.gameServer})
	handler.pool.Register(5043,&msgHandler.ChangeSignHandler{GameServer:handler.gameServer})
}
func (handler ServerMessageHandler) MessageReceived(session network.SocketSessionInterface, message interface{}) error {
	if writeMsg, ok := message.(network.WriteMessage); !ok {
		return errors.New("不是WriteMessage类型")
	} else {
		//log4g.Infof("收到消息%d", writeMsg.MsgId)
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
	//log4g.Info("Period")
	roleMap := handler.gameServer.RoleManager.GetAllOnlineRole(session)
	for k,v := range roleMap{
		connect := v.IsConnected()
		if connect{
			es := time.Now().Sub(v.GetPingTime())
			esTime := es.Seconds()
			if esTime > 15{
				log4g.Infof("超时关闭Session[%f]",esTime)
				message := &message.M2G_CloseSession{}
				message.RoleId = v.GetRoleId()
				handler.gameServer.WriteInnerMsg(v.GetGateSession(),v.GetRoleId(),10006,message)
				continue
			}
			handler.gameServer.WriteInnerMsg(v.GetGateSession(),k,5014,handler.pingMsg)
		}
	}
}

func (handler ServerMessageHandler) ExceptionCaught(session network.SocketSessionInterface, err error) {

}
