package main

import (
	"red-envelope/configer"
	"red-envelope/network"
	"net/http"
	"red-envelope/databases/mysql"
	"red-envelope/application/testRouter"
	"red-envelope/databases/redis"
	"red-envelope/application/login"
	"red-envelope/application/kexin-wallet"
)

func Init(){
	configer.InitConfiger()
	configer.InitLogger()
	configer.InitZone()

	mysql.New()
	redis.New()

	network.InitNetwork(configer.GetLogger(),nil, new(network.ReRequestRecorder))
}

func serverInfoHandler(w http.ResponseWriter, r *http.Request){
	info := map[string]string{
		"name": "red-envelope",
		"version":"1.0",
	}
	panic("err")
	network.NewSuccess(http.StatusOK, info).Response(w)
}

var routes = network.Routes{
	network.Route{http.MethodGet, "/", serverInfoHandler},
}

func main(){

	// 初始化配置
	Init()

	// 创建路由
	r := network.NewRouter(routes)

	// 路由前缀
	//r.Multiplexer = r.Multiplexer.PathPrefix("/v1").Subrouter()

	// 路由挂载
	testRouter.MountSubrouterOn(r)
	login.MountSubrouterOn(r)
	kexin_wallet.MountSubrouterOn(r)

	r.Startup(network.HTTP, 8888)
}



