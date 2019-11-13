package login

import (
	"fmt"
	"red-envelope/network"
)

type loginFailure int

const (
	dbQueryQFail loginFailure = iota + 1
	accNotExitFail
	pwdErrorFail
	tokenGenerateFail
)

func (this loginFailure) Code() string {
	return "login-" + fmt.Sprintf("%04d", this)
}

func (this loginFailure) ErrorMsg() string {
	switch this {
	case dbQueryQFail:
		return "DB error"
	case accNotExitFail:
		return "account not find"
	case pwdErrorFail:
		return "pwd error"
	case tokenGenerateFail:
		return "pwd error"
	default:
		return "fail"
	}
}

func (this loginFailure) DisplayedMsg() string {
	switch this {
	case accNotExitFail:
		return "账户不存在"
	case pwdErrorFail:
		return "密码错误"
	default:
		return network.DiscardMsg
	}
}

