package red_packet

import (
	"fmt"
	"red-envelope/network"
)

type redPacketFailure int

const (
	defaultFail     redPacketFailure = iota + 1
	dataFail
	dbQueryQFail
	numberFail
	amountFail
	payMoneyFail
	totalAmountFail
	parameterFail
)

func (this redPacketFailure) Code() string {
	return "red-packet-" + fmt.Sprintf("%04d", this)
}

func (this redPacketFailure) ErrorMsg() string {
	switch this {
	case defaultFail:
		return "server error"
	case dataFail:
		return "data error"
	case dbQueryQFail:
		return "DB error"
	case numberFail:
		return "red envelop number error"
	case amountFail:
		return "red envelop amount error"
	case payMoneyFail:
		return "pay money error"
	case totalAmountFail:
		return "total amount error"
	case parameterFail:
		return "parameter error"
	default:
		return "fail"
	}
}

func (this redPacketFailure) DisplayedMsg() string {
	switch this {
	case amountFail:
		return "红包金额超过上限"
	case payMoneyFail:
		return "今日支付金额超过上限"
	default:
		return network.DiscardMsg
	}
}

type Success struct {
	Code int         `json:"code"`
	Info string      `json:"info"`
	Data interface{} `json:"data"`
}
