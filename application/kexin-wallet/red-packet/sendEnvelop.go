package red_packet

import (
	"net/http"
	"red-envelope/network"
	"encoding/json"
	"red-envelope/application/kexin-wallet/model"
	"strconv"
	"fmt"
	"io/ioutil"
	"time"
	"red-envelope/databases/redis"
	"math/rand"
)

// 发红包
func SendRedEnvelop(requester *network.Requester, w http.ResponseWriter, r *http.Request) {
	//解析红包数据
	var redEnve model.RedEnvelopSend
	var cerr *CustomError
	var amounts []float64

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		network.NewFailure(http.StatusInternalServerError, defaultFail).AppendErrorMsg("数据读取出错: ").Response(w)
		return
	}
	err = json.Unmarshal(body, &redEnve)
	if err != nil {
		network.NewFailure(http.StatusForbidden, dataFail).AppendErrorMsg("数据解析出错: ").Response(w)
		return
	}
	config, err := model.QryBasicConfig()
	if err != nil {
		network.NewFailure(http.StatusInternalServerError, dbQueryQFail).AppendErrorMsg("查询配置信息出错: ").Response(w)
		return
	}

	// 判断红包数目金额是否符合
	if redEnve.EnvelopType == "single" {
		cerr = singleChangeCheck(&redEnve, config)
		if cerr != nil {
			network.NewFailure(cerr.ErrorCode, cerr.Failurable).AppendErrorMsg(cerr.ErrInfo).Response(w)
			return
		}
	} else if redEnve.EnvelopType == "group"  || redEnve.EnvelopType == "luck" {
		cerr = groupChangeCheck(&redEnve, config)
		if cerr != nil {
			network.NewFailure(cerr.ErrorCode, cerr.Failurable).AppendErrorMsg(cerr.ErrInfo).Response(w)
			return
		}
		amounts = make([]float64, redEnve.Number)
		var totalAmount float64

		if redEnve.EnvelopType == "group" {
			for i := 0; i < len(amounts); i++ {
				amounts[i] = demical(redEnve.TotalAmount / float64(redEnve.Number))
				totalAmount += amounts[i]
			}
			if totalAmount != redEnve.TotalAmount {
				network.NewFailure(http.StatusForbidden, totalAmountFail).AppendErrorMsg("红包总金额错误: ").Response(w)
			}

		} else if redEnve.EnvelopType == "luck" {
			var totalAmount float64
			send, total := randSeed(redEnve.Number)
			for i, v := range send {
				if i == len(send)-1 {
					break
				}
				if v == 0 {
					amounts[i] = 0.01
				} else {
					amounts[i] = demical(redEnve.TotalAmount * (float64(send[i] / total)))
				}
				totalAmount += amounts[i]
			}
			amounts[len(amounts)-1] = demical(redEnve.TotalAmount - totalAmount)
		} else {
			network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("红包金额类型不对: ").Response(w)
			return
		}

	} else {
		network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("红包数目类型不对: ").Response(w)
		return
	}

	if redEnve.PayType == "change" { // 零钱支付
		// 判断今日总的支付上限是否达到
		cerr = changePayCheck(&redEnve, config)
		if cerr != nil {
			network.NewFailure(cerr.ErrorCode, cerr.Failurable).AppendErrorMsg(cerr.ErrInfo).Response(w)
			return
		}

	} else if redEnve.PayType == "bank_card" { // 银行卡支付
		// 判断该银行卡的支付上限是否达到
		// todo
	} else {
		network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("支付方式错误: ").Response(w)
		return
	}

	network.NewSuccess(http.StatusOK, "发送成功").Response(w)
}

func demical(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func randSeed(len int) ([]int, int) {
	seed := make([]int, len)
	var totalSeed int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		seed[i] = rand.Intn(1000)
		totalSeed += seed[i]
	}
	return seed, totalSeed
}

func getTodaySecond() int {
	today := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", today, time.Local)
	todaySecond := todayLastTime.Unix() - time.Now().Local().Unix()
	return int(todaySecond)
}

// 单个红包检查
func singleChangeCheck(redEnve *model.RedEnvelopSend, config map[string]string) *CustomError {
	if redEnve.Number != 1 {
		return &CustomError{http.StatusForbidden, "红包数目不正确", numberFail}
	}
	// 单个红包金额上限
	singleEnveCeiling, err := strconv.ParseFloat(config["single_envelop_ceiling"], 64)
	if err != nil {
		return &CustomError{http.StatusInternalServerError, err.Error(), defaultFail}
	}
	if redEnve.TotalAmount > singleEnveCeiling {
		//panic("单个红包金额超过上限")
		return &CustomError{http.StatusForbidden, "单个红包金额超过上限", amountFail}
	}
	return nil
}

// 群红包检查
func groupChangeCheck(redEnve *model.RedEnvelopSend, config map[string]string) *CustomError {
	// 群红包数目和金额上限
	groupEnveCeiling, err := strconv.ParseFloat(config["group_envelop_ceiling"], 64)
	if err != nil {
		return &CustomError{http.StatusInternalServerError, err.Error(), defaultFail}
	}
	if redEnve.TotalAmount > groupEnveCeiling {
		return &CustomError{http.StatusForbidden, "群红包金额超过上限", amountFail}
	}
	groupEnveNumCeiling, err := strconv.Atoi(config["group_envelop_num_ceiling"])
	if err != nil {
		return &CustomError{http.StatusInternalServerError, err.Error(), defaultFail}
	}
	if redEnve.Number > groupEnveNumCeiling {
		return &CustomError{http.StatusForbidden, "群红包个数超过上限", numberFail}
	}
	return nil
}

// 零钱支付检查
func changePayCheck(redEnve *model.RedEnvelopSend, config map[string]string) *CustomError {
	var alreadyPay float64
	key := redEnve.UserId + ":pay:" + time.Now().Format("2006-01-02") + ":change"
	err := redis.GetInstance().GetObject(key, &alreadyPay)
	if err != nil && err.Error() != "redigo: nil returned" {
		return &CustomError{http.StatusInternalServerError, err.Error(), defaultFail}
	}
	dayPayCeiling, err := strconv.ParseFloat(config["day_pay_ceiling"], 64)
	if err != nil {
		return &CustomError{http.StatusInternalServerError, err.Error(), defaultFail}
	}
	if alreadyPay+redEnve.TotalAmount > dayPayCeiling {
		return &CustomError{http.StatusForbidden, "支付金额已超上限", payMoneyFail}
	} else {
		_, err := redis.GetInstance().Set(key, alreadyPay+redEnve.TotalAmount, getTodaySecond())
		if err != nil {
			return &CustomError{http.StatusInternalServerError, err.Error(), defaultFail}
		}
	}
	return nil
}





