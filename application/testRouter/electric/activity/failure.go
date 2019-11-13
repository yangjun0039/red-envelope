package activity

import (
	"fmt"
	"red-envelope/network"
)

type electricFailure int

const (
	dbQueryQFail electricFailure = iota + 1
	userNotExitFail
)

func (this electricFailure) Code() string {
	return "electric-" + fmt.Sprintf("%04d", this)
}

func (this electricFailure) ErrorMsg() string {
	switch this {
	case dbQueryQFail:
		return "DB 查询失败"
	case userNotExitFail:
		return "用户不存在"
	default:
		return "fail"
	}
}

func (this electricFailure) DisplayedMsg() string {
	switch this {
	case userNotExitFail:
		return "用户不存在"
	default:
		return network.DiscardMsg
	}
}
