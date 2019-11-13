package red_packet

import (
	"net/http"
	"red-envelope/network"
	"red-envelope/application/kexin-wallet/model"
	"io/ioutil"
	"encoding/json"
	"time"
)

// 发红包
func SendJRMFRedEnvelop(requester *network.Requester, w http.ResponseWriter, r *http.Request) {
	var redEnve model.RedEnvelopSend
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		network.NewFailure(http.StatusInternalServerError, defaultFail).AppendErrorMsg("数据读取出错: " + err.Error()).Response(w)
		return
	}
	err = json.Unmarshal(body, &redEnve)
	if err != nil {
		network.NewFailure(http.StatusForbidden, dataFail).AppendErrorMsg("数据解析出错: " + err.Error()).Response(w)
		return
	}
	if redEnve.Id == "" {
		network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("红包id不能为空: ").Response(w)
		return
	}
	if redEnve.UserId == "" {
		network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("用户id不能为空: ").Response(w)
		return
	}
	redEnve.SendTime = time.Now()
	redEnve.CreateTime = time.Now()
	redEnve.SendStatus = "success"

	err = model.InsertRedEnvelopV2(&redEnve)
	if err != nil {
		network.NewFailure(http.StatusInternalServerError, dbQueryQFail).AppendErrorMsg("插入发送红包数据出错: " + err.Error()).Response(w)
		return
	}
	succ := Success{200, "发送成功", nil}
	network.NewSuccess(http.StatusOK, succ).Response(w)
}

// 收红包
func ReciveJRMFRedEnvelip(requester *network.Requester, w http.ResponseWriter, r *http.Request) {
	var rer model.RedEnvelopReceive

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		network.NewFailure(http.StatusInternalServerError, defaultFail).AppendErrorMsg("数据读取出错: " + err.Error()).Response(w)
		return
	}
	err = json.Unmarshal(body, &rer)
	if err != nil {
		network.NewFailure(http.StatusForbidden, dataFail).AppendErrorMsg("数据解析出错: " + err.Error()).Response(w)
		return
	}
	if rer.RedEnvelopId == "" {
		network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("红包id不能为空: ").Response(w)
		return
	}
	if rer.UserId == "" {
		network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("用户id不能为空: ").Response(w)
		return
	}
	rer.ReceiveTime = time.Now()
	rer.CreateTime = time.Now()
	err = model.InsertEnvelopReceiveV2(&rer)
	if err != nil {
		network.NewFailure(http.StatusInternalServerError, dbQueryQFail).AppendErrorMsg("插入接收红包数据出错: " + err.Error()).Response(w)
		return
	}
	succ := Success{200, "发送成功", nil}
	network.NewSuccess(http.StatusOK, succ).Response(w)
}

// 展示红包
func ShowJRMFRedEnvelop(requester *network.Requester, w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	if userId == "" {
		network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("红包id不能为空: ").Response(w)
		return
	}
	data, err := model.QryUserEnvelopV2(userId)
	if err != nil {
		network.NewFailure(http.StatusInternalServerError, dbQueryQFail).AppendErrorMsg("查询用户红包数据错误: " + err.Error()).Response(w)
		return
	}
	respon := make([]string, len(data))
	for i, d := range data {
		respon[i] = d["red_envelop_id"]
	}
	succ := Success{200, "请求成功", respon}
	network.NewSuccess(http.StatusOK, succ).Response(w)
}

// 控制红包展示
func ShowJRMF(requester *network.Requester, w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		network.NewFailure(http.StatusForbidden, parameterFail).AppendErrorMsg("id不能为空: ").Response(w)
		return
	}
	data,err := model.QryFuncShowV2(id)
	if err != nil{
		network.NewFailure(http.StatusInternalServerError, dbQueryQFail).AppendErrorMsg("查询: " + err.Error()).Response(w)
		return
	}
	succ := Success{200, "请求成功", data}
	network.NewSuccess(http.StatusOK, succ).Response(w)
}
