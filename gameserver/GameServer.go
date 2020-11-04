package gameserver

import (
	"cmgameserver/battleManager"
	"cmgameserver/face"
	"cmgameserver/message"
	"cmgameserver/roleManager"
	"cmgameserver/roomManager"
	"fmt"
	"github.com/bianchengxiaobei/cmgo/db"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"cmgameserver/bean"
	"cmgameserver/matchManager"
	"time"
	"sync"
)

type GameServer struct {
	gameConfig               *GameBaseConfig
	IsRunning                bool
	NeedClose 				bool//是否延时关闭
	sync.RWMutex
	GameClientServer          map[int]*network.TcpClient
	gameBaseConfigPath        string
	gameSessionConfigPath     string
	gameAchieveConfigPath     string
	gameTaskConfigPath        string
	gameBoxItemInfoPath       string
	gameHeroEquipInfoPath     string
	gameSoldierEquipInfoPath  string
	gameEmailConfigPath string
	gameMaterialConfigPath string
	gameHeroConfigPath		string
	gameShopEquipConfigPath     string
	gameGiftConfigPath			string
	gameCardInfoConfigPath    string
	gameSkillConfigPath        string
	gameLevelUpgradeConfigPath string
	gameChatConfigPath 			string
	gameHeroCardConfigPath		string
	TaskConfig                *bean.TaskConfig
	BoxItemInfoConfig         *bean.BoxItemInfoConfig
	AchieveConfig             *bean.AchieveConfig
	HeroQualityEquipConfig    *bean.HeroQualityMapEquipConfig
	HeroItemIdEquipConfig     *bean.HeroItemIdMapEquipConfig
	SoldierQualityEquipConfig *bean.SoldierQualityEquipConfig
	SoldierItemEquipConfig *bean.SoldierItemIdEquipConfig
	EmailConfig					*bean.ServerEmailConfig
	MaterialConfig         *bean.ServerMaterialConfig
	ShopEquipConfig        *bean.ServerShopEquipConfig
	HeroConfig				*bean.HeroConfig
	GiftConfig 				*bean.GiftConfig
	CardInfoConfig          *bean.CardInfoConfig
	SkillConfig             *bean.ServerSkillConfig
	LevelUpgradeConfig 		*bean.LevelUpgradeConfig
	ChatConfig 				*bean.ChatConfig
	ShopHeroCardConfig      *bean.ShopHeroCardConfig
	DBManager                *db.MongoBDManager
	RoleManager              face.IRoleManager
	RoomManager              face.IRoomManager
	BattleManager            face.IBattleManager
	MatchManager				face.IMatchManager
	GameVersion              string
	Release    bool
	RankLevelRankList		[30]bean.RankListItem//段位榜单
	RoleLevelRankList		[30]bean.RankListItem//等级榜单
	RoleHeroRankList		[30]bean.RankListItem//英雄榜单
	//聊天
	ChatCache              [10]message.ChatInfo
	RankListUpdateTime		int64
	AllOnlineDiam			int32
}
type GameBaseConfig struct {
	Version              string
	Name                 string
	Id                   int32
	Release      bool
	DBAddress            string
	GateConnectConfigMap map[int]string
}

var (
	serverCodec   ServerProtocol
	serverHandler ServerMessageHandler
)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

//创建游戏服务器
func NewGameServer() *GameServer {
	server := &GameServer{
		IsRunning:        false,
		NeedClose:			false,
		GameClientServer: make(map[int]*network.TcpClient),
	}
	return server
}
func (server *GameServer) Init(gameBaseConfig string, gameSessionConfig string, gameAchieveConfig string,
	gameTaskConfig string,gameboxItemConfig string,gameHeroEquipConfig string,gameSoldierEquipConfig string,
		gameEmailConfig string,gameMaterialConfig string,gameHeroConfig string,gameShopEquipConfig string,
			gameCardInfoConfig string,gameSkillConfig string,gameGiftConfig string,gameLevelUpgradeConfig string,
				gameChatConfig string,gameShopHeroCardConfig string) {
	var (
		gameConfig network.SocketSessionConfig
	)
	rootPath, _ := os.Getwd()
	server.gameBaseConfigPath = filepath.Join(rootPath, gameBaseConfig)
	server.gameSessionConfigPath = filepath.Join(rootPath, gameSessionConfig)
	server.gameAchieveConfigPath = filepath.Join(rootPath, gameAchieveConfig)
	server.gameTaskConfigPath = filepath.Join(rootPath, gameTaskConfig)
	server.gameBoxItemInfoPath = filepath.Join(rootPath,gameboxItemConfig)
	server.gameHeroEquipInfoPath = filepath.Join(rootPath,gameHeroEquipConfig)
	server.gameSoldierEquipInfoPath = filepath.Join(rootPath,gameSoldierEquipConfig)
	server.gameEmailConfigPath = filepath.Join(rootPath,gameEmailConfig)
	server.gameMaterialConfigPath = filepath.Join(rootPath,gameMaterialConfig)
	server.gameHeroConfigPath = filepath.Join(rootPath,gameHeroConfig)
	server.gameShopEquipConfigPath = filepath.Join(rootPath,gameShopEquipConfig)
	server.gameCardInfoConfigPath = filepath.Join(rootPath,gameCardInfoConfig)
	server.gameSkillConfigPath = filepath.Join(rootPath,gameSkillConfig)
	server.gameGiftConfigPath = filepath.Join(rootPath,gameGiftConfig)
	server.gameLevelUpgradeConfigPath = filepath.Join(rootPath,gameLevelUpgradeConfig)
	server.gameChatConfigPath = filepath.Join(rootPath,gameChatConfig)
	server.gameHeroCardConfigPath = filepath.Join(rootPath,gameShopHeroCardConfig)
	LoadSessionConfig(server.gameSessionConfigPath, &gameConfig)
	server.LoadBaseConfig(server.gameBaseConfigPath)
	server.AllOnlineDiam = 736200
	//加载任务的其他配置
	//设置编解码
	serverCodec = ServerProtocol{
		pool: &ProtoMessagePool{
			messages: make(map[int32]reflect.Type),
		},
	}
	serverCodec.Init()
	//设置事件处理器
	serverHandler = ServerMessageHandler{
		gameServer: server,
		pool: &HandlerPool{
			handlers: make(map[int32]HandlerBase),
		},
		pingMsg: &message.M2C2M_GamePing{},
	}
	serverHandler.Init()
	for id, _ := range server.gameConfig.GateConnectConfigMap {
		server.GameClientServer[id] = network.NewTcpClient("tcp", &gameConfig)
		server.GameClientServer[id].SetProtocolCodec(serverCodec)
		server.GameClientServer[id].SetMessageHandler(serverHandler)
	}
	server.DBManager = db.NewMongoBD(server.gameConfig.DBAddress, 5,server.Release)
	server.RoleManager = roleManager.NewRoleManager(server)
	server.RoomManager = roomManager.NewRoomManager(server)
	server.BattleManager = battleManager.NewBattleManager(server)
	server.MatchManager = matchManager.NewMatchManager(server)
}
func (server *GameServer) Run() {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	if server.IsRunning == false {
		//开始连接网关服务器
		if server.GameClientServer != nil && len(server.GameClientServer) > 0 {
			for id, client := range server.GameClientServer {
				addr := server.gameConfig.GateConnectConfigMap[id]
				err := client.Connect(addr)
				if err != nil {
					log4g.Error(err.Error())
					return
				}
				server.RegisterGate(id)
				log4g.Infof("连接网关[%d]地址:[%s]!", id, addr)
			}

		}
		//定时月卡
		go func() {
			for {
				server.RoleManager.SendCardDiamToAllRole()
				now := time.Now().Add(time.Second * 5)
				// 计算下一个零点
				next := now.Add(time.Hour * 24)
				next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
				t := time.NewTimer(next.Sub(now))
				<-t.C
			}
		}()
		//定时赛季结算
		go func() {
			for {
				now := time.Now()
				// 计算下一个零点
				last := now.AddDate(0,0,-now.Day() + 1).AddDate(0,1,0)
				zero := time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, last.Location())
				t := time.NewTimer(zero.Sub(now))
				<-t.C
				//结算赛季奖励，前几名获得奖励
				server.RoleManager.SendSaiJiAward()
				server.AllOnlineDiam = 736200
				t2 := time.NewTimer(zero.AddDate(0,0,1).Sub(zero))
				<-t2.C
			}
		}()
		log4g.Infof("%s[%d]开始运行!", server.gameConfig.Name, server.gameConfig.Id)
		server.IsRunning = true
	}
}
func (server *GameServer)SetNeedClose(){
	server.Lock()
	defer server.Unlock()
	server.NeedClose = true
}
func (server *GameServer)SubAllOnlineDiam(value int32)(bool,int32){
	if server.AllOnlineDiam <= 0{
		return false,0
	}
	server.AllOnlineDiam -= value
	if server.AllOnlineDiam <= 0{
		return false,0
	}
	return true,server.AllOnlineDiam
}
func (server *GameServer)GetNeedClose()bool{
	server.RLock()
	defer server.RUnlock()
	return server.NeedClose
}
func (server *GameServer) Close() {
	if server.IsRunning == true {
		for _, gameClient := range server.GameClientServer {
			gameClient.Close()
		}
		//保存数据
		server.RoleManager.AllRoleQuit()
		//保存聊天
		server.SaveChatDataConfig()
		server.IsRunning = false
	}
}

//改变游戏版本号
func (server *GameServer) ChangeGameVersion(ver string) {
	server.GameVersion = ver
}

//取得游戏版本号
func (server *GameServer) GetGameVersion() string {
	return server.GameVersion
}
func (server *GameServer)GetBoxItemInfoConfig() bean.BoxItemInfoConfig{
	return *server.BoxItemInfoConfig
}
func (server *GameServer)GetHeroQualityEquipConfig() bean.HeroQualityMapEquipConfig {
	return *server.HeroQualityEquipConfig
}
func (server *GameServer)GetSoldierQualityEquipConfig() bean.SoldierQualityEquipConfig {
	return *server.SoldierQualityEquipConfig
}
func (server *GameServer)GetHeroItemIdEquipConfig() bean.HeroItemIdMapEquipConfig {
	return *server.HeroItemIdEquipConfig
}
func (server *GameServer)GetShopHeroCardConfig() bean.ShopHeroCardConfig{
	return *server.ShopHeroCardConfig
}
func (server *GameServer)GetSoldierItemIdEquipConfig() bean.SoldierItemIdEquipConfig {
	return *server.SoldierItemEquipConfig
}
func (server *GameServer)GetEmailConfig() bean.ServerEmailConfig {
	return *server.EmailConfig
}
func (server *GameServer)GetMaterialConfig() bean.ServerMaterialConfig {
	return *server.MaterialConfig
}
func (server *GameServer)GetHeroConfig() bean.HeroConfig {
	return *server.HeroConfig
}
func (server *GameServer)GetAchieveConfig() bean.AchieveConfig{
	return *server.AchieveConfig
}
func (server *GameServer)GetTaskConfig() bean.TaskConfig{
	return *server.TaskConfig
}
func (server *GameServer)GetShopEquipConfig() bean.ServerShopEquipConfig{
	return *server.ShopEquipConfig
}
func (server *GameServer)GetCardInfoConfig() bean.CardInfoConfig{
	return *server.CardInfoConfig
}
func (server *GameServer)GetSkillConfig() bean.ServerSkillConfig{
	return *server.SkillConfig
}
func (server *GameServer)GetGiftConfig() bean.GiftConfig{
	return *server.GiftConfig
}
func (server *GameServer)GetLevelUpgradeConfig() bean.LevelUpgradeConfig{
	return *server.LevelUpgradeConfig
}
func (server *GameServer) LoadNormalConfig() {
	server.LoadTaskConfig()
	server.LoadAchieveConfig()
	server.LoadBoxItemConfig()
	server.LoadHeroEquipConfig()
	server.LoadSoldierEquipConfig()
	server.LoadEmailConfig()
	server.LoadMaterialConfig()
	server.LoadHeroConfig()
	server.LoadShopEquipConfig()
	server.LoadCardInfoConfig()
	server.LoadSkillConfig()
	server.LoadGiftConfig()
	server.LoadShopHeroCardConfig()
	server.LoadLevelUpgradeConfig()
	server.LoadChatConfig()
}
func (server *GameServer)LoadChatConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameChatConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameChatConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.ChatConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.ChatConfig = config
	//转换成cache
	if len(server.ChatConfig.Datas) > 0{
		for _,v := range server.ChatConfig.Datas{
			server.AddChatInfoData(v)
		}
	}
}
//保存聊天记录
func (server *GameServer)SaveChatDataConfig(){
	if len(server.ChatCache) > 0{
		//if server.ChatConfig == nil{
		//	server.ChatConfig = &bean.ChatConfig{
		//		Datas:make([]bean.ChatData,1),
		//	}
		//}else{
		//	//清除
		//}
		server.ChatConfig = &bean.ChatConfig{

		}
		for k,v := range server.ChatCache{
			//不保存房间组队聊天信息
			if v.Name != "" && v.RoomId <= 0{
				data := bean.ChatData{
					AvatarId:v.AvatarId,
					Time:v.Time,
					Chat:v.Chat,
					Name:v.Name,
					Level:v.Level,
					RankScore:v.RankScore,
					Sex:v.Sex,
				}
				server.ChatConfig.Datas[k] = data
			}
		}
		//保存到json
		var (
			err    error
			file   *os.File
			data   []byte
		)
		defer func() {
			file.Close()
			if err := recover(); err != nil {
				fmt.Println(err)
				return
			}
		}()
		_, err = os.Stat(server.gameChatConfigPath)
		if err != nil {
			//不存在，新建
			if file, err = os.Create(server.gameChatConfigPath); err != nil {
				log4g.Info(err.Error())
			}
		}
		if file == nil{
			if file,err = os.OpenFile(server.gameChatConfigPath,os.O_RDWR,0);err != nil{
				log4g.Info(err.Error())
			}
		}
		data, err = json.Marshal(server.ChatConfig)
		if _, err = file.Write(data); err != nil {
			log4g.Info(err.Error())
		}
	}
}
func (server *GameServer)LoadLevelUpgradeConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameLevelUpgradeConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameLevelUpgradeConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.LevelUpgradeConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.LevelUpgradeConfig = config
}
func (server *GameServer)LoadShopHeroCardConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameHeroCardConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameHeroCardConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.ShopHeroCardConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.ShopHeroCardConfig = config
}
func (server *GameServer)LoadAchieveConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameAchieveConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameAchieveConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.AchieveConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.AchieveConfig = config
}
func (server *GameServer)LoadTaskConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameTaskConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameTaskConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.TaskConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.TaskConfig = config
}
func (server *GameServer)LoadBoxItemConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameBoxItemInfoPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameBoxItemInfoPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.BoxItemInfoConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.BoxItemInfoConfig = config
}
func (server *GameServer)LoadHeroConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameHeroConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameHeroConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.HeroConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.HeroConfig = config
}
func (server *GameServer)LoadSkillConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameSkillConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameSkillConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.ServerSkillConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.SkillConfig = config
}
func (server *GameServer)LoadGiftConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameGiftConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameGiftConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.GiftConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.GiftConfig = config
}
func (server *GameServer)LoadShopEquipConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameShopEquipConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameShopEquipConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.ServerShopEquipConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.ShopEquipConfig = config
}
func (server *GameServer)LoadCardInfoConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameCardInfoConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameCardInfoConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.CardInfoConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.CardInfoConfig = config
}
func (server *GameServer)LoadHeroEquipConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameHeroEquipInfoPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameHeroEquipInfoPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.HeroQualityMapEquipConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.HeroQualityEquipConfig = config
	server.HeroItemIdEquipConfig = &bean.HeroItemIdMapEquipConfig{
		Data:make(map[int32]bean.ServerEquipData),
	}
	for _,v := range server.HeroQualityEquipConfig.Data{
		if len(v) > 0{
			for _,data := range v{
				server.HeroItemIdEquipConfig.Data[data.ItemId] = data
			}
		}
	}
}
func (server *GameServer)LoadSoldierEquipConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameSoldierEquipInfoPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameSoldierEquipInfoPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.SoldierQualityEquipConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.SoldierQualityEquipConfig = config
	server.SoldierItemEquipConfig = &bean.SoldierItemIdEquipConfig{
		Data:make(map[int32]bean.ServerEquipData),
	}
	for _,v := range server.SoldierQualityEquipConfig.Data{
		if len(v) > 0{
			for _,data := range v{
				server.SoldierItemEquipConfig.Data[data.ItemId] = data
			}
		}
	}
}
func (server *GameServer)LoadEmailConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameEmailConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameEmailConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.ServerEmailConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.EmailConfig = config
}
func (server *GameServer)LoadMaterialConfig(){
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(server.gameMaterialConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(server.gameMaterialConfigPath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	config := new(bean.ServerMaterialConfig)
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	server.MaterialConfig = config
}
//加载json配置
func LoadSessionConfig(filePath string, sessionConfig *network.SocketSessionConfig) {
	var (
		err  error
		file *os.File
		data []byte
	)
	defer file.Close()
	_, err = os.Stat(filePath)
	if err != nil {
		//不存在，新建
		if file, err = os.Create(filePath); err != nil {
			fmt.Println(err)
		}
		config := network.SocketSessionConfig{
			TcpNoDelay:         true,
			TcpKeepAlive:       true,
			TcpKeepAlivePeriod: 3e9,
			TcpReadBuffSize:    1024,
			TcpWriteBuffSize:   1024,
			ReadChanLen:        1024,
			WriteChanLen:       1024,
			PeriodTime:         5e9,
		}
		data, err = json.Marshal(config)
		if _, err = file.Write(data); err != nil {
			fmt.Println(err)
		}
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(filePath)
		if err != nil {
			panic(err)

		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	err = json.Unmarshal(data, sessionConfig)
	if err != nil {
		panic(err)
	}
}
func (server *GameServer) LoadBaseConfig(filePath string) {
	var (
		err    error
		file   *os.File
		data   []byte
		config *GameBaseConfig
	)
	defer func() {
		file.Close()
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	_, err = os.Stat(filePath)
	if err != nil {
		//不存在，新建
		if file, err = os.Create(filePath); err != nil {
			fmt.Println(err)
		}
		config = &GameBaseConfig{
			Version:              "1.0.12",
			Name:                 "游戏服务器",
			Id:                   1,
			Release:false,
			DBAddress:"127.0.0.1",
			GateConnectConfigMap: make(map[int]string),
		}
		config.GateConnectConfigMap[1] = "127.0.0.1:8001"
		data, err = json.Marshal(config)
		if _, err = file.Write(data); err != nil {
			fmt.Println(err)
		}
	}
	if data == nil || len(data) == 0 {
		file, err = os.Open(filePath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	if config == nil {
		config = new(GameBaseConfig)
		err = json.Unmarshal(data, config)
		if err != nil {
			panic(err)
		}
	}
	server.gameConfig = config
	server.GameVersion = config.Version
	server.Release = config.Release
	server.LoadNormalConfig()
}
func (server *GameServer) RegisterGate(gateId int) {
	//发送给网关服务器注册自己
	message := &message.M2G_RegisterGate{
		Id: server.gameConfig.Id,
	}
	innerMsg := network.InnerWriteMessage{
		//RoleId:make([]int64,0),
		MsgData: message,
	}
	server.GameClientServer[gateId].Session.WriteMsg(10000, innerMsg)
}
func (server *GameServer) GetId() int32 {
	return server.gameConfig.Id
}
func (server *GameServer) GetDBManager() *db.MongoBDManager {
	return server.DBManager
}
func (server *GameServer) GetRoleManager() face.IRoleManager {
	return server.RoleManager
}
func (server *GameServer) GetRoomManager() face.IRoomManager {
	return server.RoomManager
}
func (server *GameServer) GetBattleManager() face.IBattleManager {
	return server.BattleManager
}
func (server *GameServer) GetMatchManager() face.IMatchManager {
	return server.MatchManager
}
func (server *GameServer)GetRankListTime() int64{
	return server.RankListUpdateTime
}
func (server *GameServer)SetRankListTime(time int64){
	server.RankListUpdateTime = time
}
func (server *GameServer)GetRankLevelRankList () *[30]bean.RankListItem{
	return &server.RankLevelRankList
}
func (server *GameServer)GetRoleHeroCountRankList() *[30]bean.RankListItem{
	return &server.RoleHeroRankList
}
func (server *GameServer)GetRoleLevelRankList() *[30]bean.RankListItem{
	return &server.RoleLevelRankList
}
func (server *GameServer)GetChatInfo() *[10]message.ChatInfo{
	return &server.ChatCache
}
func (server *GameServer)AddChatInfoData(data bean.ChatData){
	server.AddChatInfo(data.AvatarId,data.Name,data.Chat,data.Time,
		data.Level,data.RankScore,data.Sex,data.RoomId)
}
func (server *GameServer)AddChatInfo(avatarId int32,name string,content string,time int64,level int32,
	rankScore int32,sex int32,roomId int32){
	for i:=0;i<10;i++{
		chat := server.ChatCache[i]
		if chat.AvatarId == 0{
			//说明没有放入
			chat.AvatarId = avatarId
			chat.Time = time
			chat.Chat = content
			chat.Name = name
			chat.Level = level
			chat.RoomId = roomId
			chat.Sex = sex
			chat.RankScore = rankScore
			server.ChatCache[i] = chat
			return
		}
	}
	//说明已经超出了
	//超出就根据队列重排
	for i:=0;i<10;i++{
		if i < 9{
			server.ChatCache[i] = server.ChatCache[i+1]
		}else{
			chat := server.ChatCache[i]
			chat.AvatarId = avatarId
			chat.Time = time
			chat.Chat = content
			chat.Name = name
			chat.Level = level
			chat.RoomId = roomId
			chat.Sex = sex
			chat.RankScore = rankScore
			server.ChatCache[i] = chat
		}
	}
}
func (server *GameServer) WriteInnerMsg(session network.SocketSessionInterface, roleId int64, msgId int, msg proto.Message) {
	innerMsg := network.InnerWriteMessage{
		MsgData: msg,
	}
	innerMsg.RoleId = roleId
	if err := session.WriteMsg(msgId, innerMsg); err != nil {
		log4g.Info(err.Error())
	}
}