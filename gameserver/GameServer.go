package gameserver

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"os"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"reflect"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/json-iterator/go"
)

type GameServer struct{
	gameConfig 				*GameBaseConfig
	IsRunning 				bool
	GameClientServer		map[int]*network.TcpClient
	gameBaseConfigPath		string
	gameSessionConfigPath	string
}
type GameBaseConfig struct {
	Name                 string
	Id                   int32
	GateConnectConfigMap map[int]string
}

var (
	serverCodec 	ServerProtocol
	serverHandler	ServerMessageHandler
)
var json = jsoniter.ConfigCompatibleWithStandardLibrary
//创建游戏服务器
func NewGameServer() *GameServer{
	server := &GameServer{
		IsRunning:false,
		GameClientServer:make(map[int]*network.TcpClient),
	}
	return server
}
func (server *GameServer)Init(gameBaseConfig string,gameSessionConfig string){
	var (
		gameConfig 	network.SocketSessionConfig
	)
	rootPath,_ := os.Getwd()
	server.gameBaseConfigPath = filepath.Join(rootPath, gameBaseConfig)
	server.gameSessionConfigPath = filepath.Join(rootPath,gameSessionConfig)
	LoadSessionConfig(server.gameSessionConfigPath,&gameConfig)
	server.LoadBaseConfig(server.gameBaseConfigPath)

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
	}
	for id,_:=range server.gameConfig.GateConnectConfigMap {
		server.GameClientServer[id] = network.NewTcpClient("tcp",&gameConfig)
		server.GameClientServer[id].SetProtocolCodec(serverCodec)
		server.GameClientServer[id].SetMessageHandler(serverHandler)
	}
}
func (server *GameServer)Run(){
	defer func() {
		if err := recover(); err != nil {
			//log4g.Error("网关服务器监听出错!")
			fmt.Println(err)
			return
		}
	}()
	if server.IsRunning == false {
		//开始连接网关服务器
		if server.GameClientServer != nil && len(server.GameClientServer) > 0{
			for id,client := range server.GameClientServer{
				addr := server.gameConfig.GateConnectConfigMap[id]
				client.Connect(addr)
				server.RegisterGate(id)
				log4g.Infof("连接网关[%d]地址:[%s]!",id,addr)
			}

		}
		log4g.Infof("%s[%d]开始运行!",server.gameConfig.Name,server.gameConfig.Id)
		server.IsRunning = true
	}
}
func (server *GameServer) Close() {
	if server.IsRunning == true {
		for _,gameClient := range server.GameClientServer{
			gameClient.Close()
		}
		server.IsRunning = false
	}
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
			ReadChanLen:        1,
			WriteChanLen:       1,
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
			Name:                 "游戏服务器",
			Id:                   1,
			GateConnectConfigMap: make(map[int]string),
		}
		config.GateConnectConfigMap[1] = "127.0.0.1:8000"
		fmt.Println(len(config.GateConnectConfigMap))
		data, err = json.Marshal(config)
		if _, err = file.Write(data); err != nil {
			fmt.Println(err)
		}
	}
	if data == nil || len(data) == 0{
		file, err = os.Open(filePath)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	if config == nil{
		config = new(GameBaseConfig)
		err = json.Unmarshal(data, config)
		if err != nil {
			panic(err)
		}
	}
	server.gameConfig = config
}
func (server *GameServer)RegisterGate(gateId int)  {
	//发送给网关服务器注册自己
	message := &message2.M2G_RegisterGate{
		Id:server.gameConfig.Id,
	}
	server.GameClientServer[gateId].Session.WriteMsg(10000,message)
}