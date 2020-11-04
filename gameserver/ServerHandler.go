package gameserver

import (
	"cmgameserver/msgHandler"
	"errors"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"time"
	"fmt"
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

	handler.pool.Register(100001,&msgHandler.GMCommandHandler{GameServer:handler.gameServer})

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
	handler.pool.Register(5042,&msgHandler.ChangeSignHandler{GameServer:handler.gameServer})
	handler.pool.Register(5044,&msgHandler.GetBoxAwardHandler{GameServer:handler.gameServer})
	handler.pool.Register(5046,&msgHandler.SellItemHandler{GameServer:handler.gameServer})
	handler.pool.Register(5048,&msgHandler.GetEmailAwardHandler{GameServer:handler.gameServer})
	handler.pool.Register(5050,&msgHandler.DeleteEmailHandler{GameServer:handler.gameServer})
	handler.pool.Register(5052,&msgHandler.UseItemHandler{GameServer:handler.gameServer})
	handler.pool.Register(5054,&msgHandler.ReqPauseBattleHandler{GameServer:handler.gameServer})
	handler.pool.Register(5055,&msgHandler.AgreePauseBattleHandler{GameServer:handler.gameServer})
	handler.pool.Register(5058,&msgHandler.ReqRankListHandler{GameServer:handler.gameServer})
	handler.pool.Register(5061,&msgHandler.BindZhangHaoHandler{GameServer:handler.gameServer})
	handler.pool.Register(5062,&msgHandler.AutoMatchHandler{GameServer:handler.gameServer})
	handler.pool.Register(5066,&msgHandler.RoleQuitMatchTeamHandler{GameServer:handler.gameServer})
	handler.pool.Register(5068,&msgHandler.TeamStartMatchHandler{GameServer:handler.gameServer})
	handler.pool.Register(5069,&msgHandler.MatchRoomPrepareHandler{GameServer:handler.gameServer})
	handler.pool.Register(5072,&msgHandler.PaiWeiLoadFinishedHandler{GameServer:handler.gameServer})
	handler.pool.Register(5073,&msgHandler.CancelStartMatchHandler{GameServer:handler.gameServer})
	handler.pool.Register(5076,&msgHandler.UpdateAchieveDataHandler{GameServer:handler.gameServer})
	handler.pool.Register(5077,&msgHandler.BuyShopEquipHandler{GameServer:handler.gameServer})
	handler.pool.Register(5080,&msgHandler.BuyShopDiamHandler{GameServer:handler.gameServer})
	handler.pool.Register(5082,&msgHandler.CheckBuyDiamHandler{GameServer:handler.gameServer})
	handler.pool.Register(5083,&msgHandler.BuyShopCardHandler{GameServer:handler.gameServer})
	handler.pool.Register(5084,&msgHandler.CheckBuyCardHandler{GameServer:handler.gameServer})
	handler.pool.Register(5085,&msgHandler.CompleteGuideHandler{GameServer:handler.gameServer})
	handler.pool.Register(5086,&msgHandler.BuyShopBoxHandler{GameServer:handler.gameServer})
	handler.pool.Register(5088,&msgHandler.BuyShopHeroCardHandler{GameServer:handler.gameServer})
	handler.pool.Register(5090,&msgHandler.UpgradeHeroLevelHandler{GameServer:handler.gameServer})
	handler.pool.Register(5092,&msgHandler.QuickGetAllNoReadEmailHandler{GameServer:handler.gameServer})
	handler.pool.Register(5094,&msgHandler.DeleteAllReadEmailHandler{GameServer:handler.gameServer})
	handler.pool.Register(5096,&msgHandler.DeleteBagItemHandler{GameServer:handler.gameServer})
	handler.pool.Register(5097,&msgHandler.GetGiftHandler{GameServer:handler.gameServer})
	handler.pool.Register(5099,&msgHandler.ChangeRoomCityIdHandler{GameServer:handler.gameServer})
	handler.pool.Register(5100,&msgHandler.ReStartPauseBattleHandler{GameServer:handler.gameServer})
	handler.pool.Register(5101,&msgHandler.ReqBattleBugHandler{GameServer:handler.gameServer})
	handler.pool.Register(5103,&msgHandler.EnterBattleStateHandler{GameServer:handler.gameServer})
	handler.pool.Register(5104,&msgHandler.InviteRoomHandler{GameServer:handler.gameServer})
	handler.pool.Register(5106,&msgHandler.OnlinePlayerHandler{GameServer:handler.gameServer})
	handler.pool.Register(5109,&msgHandler.ChangePasswordHandler{GameServer:handler.gameServer})
	handler.pool.Register(5111,&msgHandler.CheckOnlineGetDiamHandler{GameServer:handler.gameServer})
}
func (handler ServerMessageHandler) MessageReceived(session network.SocketSessionInterface, message interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
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

func (handler ServerMessageHandler) SessionClosed(session network.SocketSessionInterface,err error) {
	if err != nil{
		log4g.Info(err.Error())
	}
}

func (handler ServerMessageHandler) SessionPeriod(session network.SocketSessionInterface) {
	//log4g.Info("Period")
	roleMap := handler.gameServer.RoleManager.GetAllOnlineRole(session)
	for k,v := range roleMap{
		connect := v.IsConnected()
		if connect{
			es := time.Now().Sub(v.GetPingTime())
			esTime := es.Seconds()
			if esTime > 35{
				//log4g.Infof("[%d]超时关闭Session[%f]",k,esTime)
				v.SetConnected(false)
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
	log4g.Info(err.Error())
}
