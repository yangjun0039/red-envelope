package configer

import (
	"testing"
	"fmt"
	"github.com/Unknwon/goconfig"
)

//var cfg *goconfig.ConfigFile

func init() {
	config, err := goconfig.LoadConfigFile("./app.conf")
	if err != nil {
		panic("get config file error:"+ err.Error())
	}
	cfg = config
}

func TestConf(t *testing.T) {
	mysql, err := cfg.GetSection("mysql")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mysql["production"])
}
