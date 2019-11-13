package model

import
(
	"testing"
	"fmt"
	"red-envelope/configer"
	"red-envelope/databases/mysql"
	"time"
)

func init(){
	configer.InitConfiger()
	mysql.New()
}

func TestQryUserPermission(t *testing.T){
	data,err := QryUserPermission("4")
	if err != nil{
		fmt.Println(err)
	}
	for _,v := range(data){
		fmt.Println(v)
	}
}

func TestDayTime(t *testing.T){
	todayTime := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayTime, time.Local)
	timediff := todayLastTime.Unix()-time.Now().Local().Unix()

	today := time.Now().Format("2006-01-02")

	fmt.Println(today)
	fmt.Println(timediff)

}

type aa struct{
	i int
}
func TestAa( *testing.T){

}

