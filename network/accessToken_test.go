package network

import (
	"testing"
	"fmt"
)

func TestAccessToken(t *testing.T) {
	at := AccessToken{
		"yangjun",
		"03050039",
	}
	secretKey := "123456789"

	token,err := at.SignedString(secretKey)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(token)

    fmt.Println("-------------------------")

	err = at.ParseWith(token, secretKey)
	if err != nil{
		fmt.Println(err)
	}
}

func TestGetToken(t *testing.T){
	at := AccessToken{
		"yangjun",
		"03050039",
	}
	secretKey := "123456789"

	token,err := at.SignedString(secretKey)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(token)
}


func TestVerfyToken(t *testing.T){
	secretKey := "123456789"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzI1Mzc1OTksImp0aSI6InlhbmdqdW4iLCJpc3MiOiJ6aWh1YW4iLCJuYmYiOjE1NzExMDgyOTAsInN1YiI6IkFjY2Vzc1Rva2VuIiwicGFzc3dvcmQiOiIwMzA1MDAzOSJ9.bdnWPvraoqUAFkGPWM_1tNhyp9uaT_2teNOIqv1yC-U"
    at := AccessToken{}
	fmt.Println(at)
	err := at.ParseWith(token, secretKey)
	if err != nil{
		fmt.Println(err)
	} else {
		fmt.Println("token is ok")
		fmt.Println(at)
	}
}
