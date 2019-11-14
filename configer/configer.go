package configer

import (
	"os"
	"path/filepath"
	"fmt"
	"time"
	"github.com/Unknwon/goconfig"
)

var ServerMap map[string]string

type mysqlConfig struct {
	DataSourceName string
}

type redisConfig struct {
	Addr     string
	Db       string
	Password string
}

// 配置路径
func configPath() string {
	var filePath string
	if GetRunningMode() == Production {
		filePath = "/root/"
	} else {
		GOPATH := os.Getenv("GOPATH")
		if GOPATH != "" {
			filePath = filepath.Join(GOPATH, "src/red-envelope/configer/")
		} else {
			filePath = "./configer/"
		}
	}
	fmt.Println(filePath)
	return filePath
}

// 初始化时区
func InitZone() {
	local, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = local
}

var Cfg *goconfig.ConfigFile

// 初始化配置文件
func InitConfiger() {
	config, err := goconfig.LoadConfigFile(configPath()+"/app.conf")
	if err != nil {
		panic("init config file error:"+ err.Error())
	}
	Cfg = config
	ServerMap = ServerConf()
}

// 数据库配置
func MySqlConfig() mysqlConfig {
	model := "production"
	mysql, err := Cfg.GetSection("mysql")
	if err != nil {
		panic("mysql配置读取出错: "+err.Error())
	}
	if data := mysql[model]; data != "" {
		return mysqlConfig{
			data,
		}
	}
	panic("mysqlConfig is nil")
}

// redis配置
func RedisConfig() redisConfig {
	redis, err := Cfg.GetSection("redis")
	if err != nil {
		panic("redis配置读取出错: "+err.Error())
	}
	var (
		password string
		ok       bool
	)
	if password, ok = redis["password"]; !ok {
		password = ""
	}
	if redis["addr"] == "" || redis["db"] == "" {
		panic("redis is nil")
	}
	return redisConfig{
		Addr:     redis["addr"],
		Db:       redis["db"],
		Password: password,
	}
}

// server
func ServerConf() map[string]string{
	server, err := Cfg.GetSection("server")
	if err != nil {
		panic("读取server配置文件出错: "+err.Error())
	}
	return server
}
