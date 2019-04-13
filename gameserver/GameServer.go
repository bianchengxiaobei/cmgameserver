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
)

type GameServer struct {
	gameConfig               *GameBaseConfig
	IsRunning                bool
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
	TaskConfig                *bean.TaskConfig
	BoxItemInfoConfig         *bean.BoxItemInfoConfig
	AchieveConfig             *bean.AchieveConfig
	HeroQualityEquipConfig    *bean.HeroQualityMapEquipConfig
	HeroItemIdEquipConfig     *bean.HeroItemIdMapEquipConfig
	SoldierQualityEquipConfig *bean.SoldierQualityEquipConfig
	SoldierItemEquipConfig *bean.SoldierItemIdEquipConfig
	EmailConfig					*bean.ServerEmailConfig
	MaterialConfig         *bean.ServerMaterialConfig
	DBManager                *db.MongoBDManager
	RoleManager              face.IRoleManager
	RoomManager              face.IRoomManager
	BattleManager            face.IBattleManager
	GameVersion              string
	RankLevelRankList		[30]bean.RankListItem//段位榜单
	RoleLevelRankList		[30]bean.RankListItem//等级榜单
	RoleHeroRankList		[30]bean.RankListItem//英雄榜单
	RankListUpdateTime		int64
}
type GameBaseConfig struct {
	Version              string
	Name                 string
	Id                   int32
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
		GameClientServer: make(map[int]*network.TcpClient),
	}
	return server
}
func (server *GameServer) Init(gameBaseConfig string, gameSessionConfig string, gameAchieveConfig string,
	gameTaskConfig string,gameboxItemConfig string,gameHeroEquipConfig string,gameSoldierEquipConfig string,
		gameEmailConfig string,gameMaterialConfig string) {
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
	LoadSessionConfig(server.gameSessionConfigPath, &gameConfig)
	server.LoadBaseConfig(server.gameBaseConfigPath)
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
	server.DBManager = db.NewMongoBD("127.0.0.1", 5)
	server.RoleManager = roleManager.NewRoleManager(server)
	server.RoomManager = roomManager.NewRoomManager(server)
	server.BattleManager = battleManager.NewBattleManager(server)
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
		log4g.Infof("%s[%d]开始运行!", server.gameConfig.Name, server.gameConfig.Id)
		server.IsRunning = true
	}
}
func (server *GameServer) Close() {
	if server.IsRunning == true {
		for _, gameClient := range server.GameClientServer {
			gameClient.Close()
		}
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
func (server *GameServer)GetSoldierItemIdEquipConfig() bean.SoldierItemIdEquipConfig {
	return *server.SoldierItemEquipConfig
}
func (server *GameServer)GetEmailConfig() bean.ServerEmailConfig {
	return *server.EmailConfig
}
func (server *GameServer)GetMaterialConfig() bean.ServerMaterialConfig {
	return *server.MaterialConfig
}
func (server *GameServer) LoadNormalConfig() {
	server.LoadTaskConfig()
	server.LoadAchieveConfig()
	server.LoadBoxItemConfig()
	server.LoadHeroEquipConfig()
	server.LoadSoldierEquipConfig()
	server.LoadEmailConfig()
	server.LoadMaterialConfig()
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
		json = nil
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

func (server *GameServer) WriteInnerMsg(session network.SocketSessionInterface, roleId int64, msgId int, msg proto.Message) {
	innerMsg := network.InnerWriteMessage{
		MsgData: msg,
	}
	innerMsg.RoleId = roleId
	if err := session.WriteMsg(msgId, innerMsg); err != nil {
		log4g.Info(err.Error())
	}
}
