package redis

import (
	"testing"
	"fmt"
	"red-envelope/configer"
	"time"
)

func TestGetValue(t *testing.T) {
	delegateId := "2010002"
	configer.InitConfiger()
	New()
	redisValue, err := GetInstance().GetString(delegateId)
	if err != nil {
		fmt.Printf("redis error , get key=%v , err:%v ", delegateId, err.Error())
	}
	fmt.Println(redisValue)
}

func TestSetValue(t *testing.T) {
	delegateId := "2010002"
	code := "071CUt2s0LH6Nd1j9S4s0Kwr2s0CUt2L"
	configer.InitConfiger()
	New()
	redisValue, err := GetInstance().Set(delegateId, code, 20)
	if err != nil {
		fmt.Printf("redis error , get key=%v:code:%v , err:%v ", delegateId, code, err.Error())
	}
	fmt.Println(redisValue)
}

func TestSetSingleAmountValue(t *testing.T) {
	single := "single"
	amount := 200
	configer.InitConfiger()
	New()
	redisValue, err := GetInstance().Set(single, amount, -1)
	if err != nil {
		fmt.Printf("redis error , get key=%v:code:%v , err:%v ", single, amount, err.Error())
	}
	fmt.Println(redisValue)

}

func TestGetSingleAmountValue(t *testing.T) {
	single := "single"
	configer.InitConfiger()
	New()
	redisValue, err := GetInstance().GetInt(single)
	if err != nil {
		fmt.Printf("redis error , get key=%v , err:%v ", single, err.Error())
	}
	fmt.Println(redisValue)
}

func TestSetPay(t *testing.T) {
	uid := "1101"
	date := "2019-11-08"
	card_num := "11111111"
	configer.InitConfiger()
	New()

	today := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", today, time.Local)
	timediff := todayLastTime.Unix() - time.Now().Local().Unix()

	// 银行卡支付
	redisValue, err := GetInstance().Set(uid+":pay:"+date+":"+card_num, 1000.1, int(timediff))
	if err != nil {
		fmt.Println("redis error,", err.Error())
	}
	fmt.Println(redisValue)

	// 余额支付
	redisValue, err = GetInstance().Set(uid+":pay:"+date+":change", 2000, int(timediff))
	if err != nil {
		fmt.Println("redis error ,", err.Error())
	}
	fmt.Println(redisValue)

	// 总的支付
	redisValue, err = GetInstance().Set(uid+":pay:"+date+":total", 3000, int(timediff))
	if err != nil {
		fmt.Println("redis error ,", err.Error())
	}
	fmt.Println(redisValue)
}

func TestGetPay(t *testing.T) {

	uid := "1101"
	date := "2019-11-08"
	card_num := "11111111"
	configer.InitConfiger()
	New()

	var val float64

	err := GetInstance().GetObject(uid+":pay:"+date+":"+card_num, &val)
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			fmt.Println("nil return!!!")
		} else {
			fmt.Println("redis error ,  ", err.Error())
		}
	}
	fmt.Println(val)

}

func TestTime(t *testing.T){
	fmt.Println(time.Now())
}
