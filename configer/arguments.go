package configer

import (
	"sync"
	"flag"
	"fmt"
	"os"
	"errors"
)

var once sync.Once
var configuration = NewSingleton(once)

func ParseArguments(args []string) {
	var rModeStr string
	flag.StringVar(&rModeStr, "RM", "Debug", "运行模式(RunningMode)，默认Debug模式")

	var dir string
	flag.StringVar(&dir, "dir", ".", "静态文件存储路径. 默认为当前文件夹")

	flag.Parse()

	if err := setRunningMode(RunningMode(rModeStr)); err != nil {
		fmt.Println("fail to parse arguments", err)
		os.Exit(-1)
	}

	setStaticFileDir(dir)

	fmt.Println("当前的运行模式为:", GetRunningMode().String())
}

type RunningMode string

const (
	Debug      RunningMode = "Debug"
	Staging    RunningMode = "Staging"
	Production RunningMode = "Production"
)

func setRunningMode(mode RunningMode) error {
	switch mode {
	case Debug, Staging, Production:
		configuration.Set("RunningMode", string(mode))
	default:
		return errors.New("-RM参数错误: 无效的运行模式")
	}
	return nil
}

func GetRunningMode() RunningMode {
	return RunningMode(configuration.Get("RunningMode"))
}

func (this RunningMode) String() string {
	switch this {
	case Debug:
		return "调试模式"
	case Staging:
		return "演练模式"
	case Production:
		return "生产模式"
	default:
		return "未定义的运行模式"
	}
}

func setStaticFileDir(dir string) {
	if GetRunningMode() == Debug {
		configuration.Set("dir", "./")
		return
	}
	configuration.Set("dir", dir)
}

func GetStaticFileDir() string {
	return configuration.Get("dir")
}

