package network

import (
	"sync"
	"red-envelope/foundation"
	//"red-envelope/configer"
	"strings"
)

type NetworkManager struct {
	// 日志记录器
	Logger *foundation.Logger

	// 互斥元
	mutex *sync.Mutex

	// 跨域白名单
	whiteOrigins map[string]bool

	// JWT密码箱
	JWTKeyBox JWTSecretKeyBox

	// Token生成工厂
	TokenFactory map[string]func() Tokenable

	// RSA密码箱
	RSAKeyBox RSAKeyBox

	// 请求记录器
	RequestRecorder RequestRecorder
}

func (m *NetworkManager) RegisterToken(tokenGenerator func() Tokenable) {
	m.TokenFactory[tokenGenerator().Subject()] = tokenGenerator
}

func (m *NetworkManager) GetTokenGenerator(subject string) func() Tokenable {
	return m.TokenFactory[subject]
}

var SharedManager *NetworkManager
var onceForInitializing sync.Once

type JWTSecretKeyBox interface {
	NewSecretKey(key string) string
	GetSecretKey(key string) (string, error)
}

type RSAKeyBox interface {
	GetPublickeyPEM(key string) ([]byte, error)
}

type RequestRecorder interface {
	SetLastTimestamp(key string, value int64) error
	GetLastTimestamp(key string) (int64, error)
}

type WhiteIPManager struct {
	IPs map[string]bool
}

// 初始化network包
func InitNetwork(Logger *foundation.Logger, keyBox JWTSecretKeyBox,reqRecorder RequestRecorder) {
	onceForInitializing.Do(func() {
		mutex := &sync.Mutex{}
		SharedManager = &NetworkManager{mutex: mutex}
		SharedManager.whiteOrigins = map[string]bool{}

		// 白名单功能
		//origins := configer.NetworkSettings("white_origins")
		//switch origins.(type) {
		//case []interface{}:
		//	for _, origin := range origins.([]interface{}) {
		//		switch origin.(type) {
		//		case string:
		//			SharedManager.whiteOrigins[origin.(string)] = true
		//		}
		//	}
		//}
		SharedManager.Logger = Logger
		SharedManager.JWTKeyBox = keyBox
		SharedManager.RequestRecorder = reqRecorder
		SharedManager.TokenFactory = make(map[string]func() Tokenable)

		// 注册不同类型的token
		SharedManager.RegisterToken(func() Tokenable {
			return new(AccessToken)
		})
	})
}

func (m *WhiteIPManager) isWhite(ip string) bool {
	value, ok := m.IPs[ip]
	if value && ok {
		return true //直接匹配
	}

	var craIP = strings.Split(ip, ".")
	for val, ok := range m.IPs {
		var crackingIP = strings.Split(val, ".")
		intervalIP := strings.Split(crackingIP[3], "/")
		if len(intervalIP) == 1 {
			continue //非区间段IP且不匹配
		}
		if craIP[0] != crackingIP[0] || craIP[1] != crackingIP[1] || craIP[2] != crackingIP[2] {
			continue //区间段IP，但不匹配
		}
		if (craIP[3] >= intervalIP[0] && craIP[3] <= intervalIP[1] && ok) {
			return true //区间段IP，匹配
		}
		//区间段IP，不在范围内
	}
	return false
}
