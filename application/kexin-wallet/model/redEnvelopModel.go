package model

import (
	"time"
	"red-envelope/databases/mysql"
)

type RedEnvelopSend struct {
	Id          string    `json:"id"`
	UserId      string    `json:"user_id"`
	EnvelopName string    `json:"envelop_name"`
	EnvelopWish string    `json:"envelop_wish"`
	EnvelopType string    `json:"envelop_type"`
	Number      int       `json:"number"`
	TotalAmount float64   `json:"total_amount"`
	SendTime    time.Time `json:"send_time"`
	PayType     string    `json:"pay_type"`
	CardNum     string    `json:"card_num"`
	SendStatus  string    `json:"send_status"`
	CreateTime  time.Time `json:"create_time"`
}

type RedEnvelopReceive struct {
	RedEnvelopId  string    `json:"red_envelop_id"`
	UserId        string    `json:"user_id"`
	ReceiveAmount float64   `json:"receive_amount"`
	ReceiveTime   time.Time `json:"receive_time"`
	CreateTime    time.Time `json:"create_time"`
}

//type FuncShow struct {
//	Id         int       `json:"id"`
//	Name       string    `json:"name"`
//	IsShow     bool      `json:"is_show"`
//	CreateTime time.Time `json:"create_time"`
//	UpdateTime time.Time `json:"update_time"`
//}

func QryBasicConfig() (map[string]string, error) {
	sql := `select * from basic_config where is_effect = true;`
	data, err := mysql.FetchRowD(sql)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func InsertRedEnvelopV2(res *RedEnvelopSend) error {
	sql := `
	insert into red_envelop_send (id, user_id, envelop_name, envelop_wish, envelop_type, number,
	total_amount, send_time, send_status, create_time)
	value (?,?,?,?,?,?,?,?,?,?)
    `
	_, err := mysql.Insert(sql, res.Id, res.UserId, res.EnvelopName, res.EnvelopWish, res.EnvelopType, res.Number, res.TotalAmount, res.SendTime, res.SendStatus, res.CreateTime)
	if err != nil {
		return err
	}
	return nil
}

func InsertEnvelopReceiveV2(rer *RedEnvelopReceive) error {
	sql := `
	insert into red_envelop_receive (red_envelop_id, user_id, receive_time, create_time)
	value (?,?,?,?)
	`
	_, err := mysql.Insert(sql, rer.RedEnvelopId, rer.UserId, rer.ReceiveTime, rer.CreateTime)
	if err != nil {
		return err
	}
	return nil
}

func QryUserEnvelopV2(userId string) ([]map[string]string, error) {
	sql := `
	select red_envelop_id from red_envelop_receive where user_id = ?
	`
	data, err := mysql.FetchRowsD(sql, userId)
	if err != nil {
		return data, err
	}
	return data, nil
}

func QryFuncShowV2(id string) (map[string]string, error) {
	sql := `
	select * from func_show where id = ?
	`
	data, err := mysql.FetchRowD(sql, id)
	if err != nil {
		return data, err
	}
	return data, nil
}
