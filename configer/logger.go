package configer

import (
	"red-envelope/foundation"
	"os"
	"path/filepath"
	"fmt"
)

var Logger *foundation.Logger

// 启动日志管理
func InitLogger() {
	GOPATH := os.Getenv("GPPATH")
	var path = "./logs/"
	if GOPATH != "" {
		path = filepath.Join(GOPATH, "src", "red-envelope", "logs") + "/"
	}
	fmt.Println("path", path)
	Logger = foundation.NewLogger("red-envelope", path)
}

// 获取日志
func GetLogger() *foundation.Logger {
	if Logger != nil {
		return Logger
	}
	panic("init logger fail")
}
