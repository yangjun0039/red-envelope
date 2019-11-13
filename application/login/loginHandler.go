package login

import (
	"net/http"
	"red-envelope/network"
	"fmt"
)

func Login(requester *network.Requester, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login")
	accInfo, err := QryUserInfo(requester.Token.Identifier())
	if err != nil {
		network.NewFailure(http.StatusInternalServerError, dbQueryQFail).AppendErrorMsg(err.Error()).Response(w)
		return
	}
	if len(accInfo) == 0 {
		network.NewFailure(http.StatusBadRequest, accNotExitFail).Response(w)
		return

	}
	if accInfo["password"] != requester.Token.Validation() {
		network.NewFailure(http.StatusBadRequest, pwdErrorFail).Response(w)
		return
	}
	token ,err := generate(requester)
	if err != nil{
		network.NewFailure(http.StatusInternalServerError, tokenGenerateFail)
	}
	respInfo := map[string]string{
		"token": token,
	}
	network.NewSuccess(http.StatusOK, respInfo).Response(w)
}

func generate(requester *network.Requester) (string, error) {
	key := fmt.Sprintf("%v:%v", requester.Token.Subject(), requester.Token.Identifier())
	token, err := requester.Token.SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}
