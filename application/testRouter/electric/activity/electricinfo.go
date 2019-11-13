package activity

import (
	"red-envelope/network"
	"net/http"
	"red-envelope/application/testRouter/electric/model"
	"fmt"
)

func UserInfo(requester *network.Requester, w http.ResponseWriter, r *http.Request){
	fmt.Println("UserInfo")
	//panic("eeeeeee")
	data,err := model.QryUserInfo()
	if err != nil{
		network.NewFailure(http.StatusInternalServerError, dbQueryQFail).AppendErrorMsg(err.Error()).Response(w)
	}
	network.NewSuccess(http.StatusOK, data).Response(w)
}

func AccInfo(requester *network.Requester, w http.ResponseWriter, r *http.Request){
	fmt.Println("AccInfo")
	data,err := model.QryAccInfo()
	if err != nil{
		network.NewFailure(http.StatusInternalServerError, dbQueryQFail).AppendErrorMsg(err.Error()).Response(w)
	}
	network.NewSuccess(http.StatusOK, data).Response(w)
}