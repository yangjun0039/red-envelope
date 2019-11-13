package model

import
(
	"red-envelope/databases/mysql"
)

func QryUserInfo() ([]map[string]string, error){
	sql := `select * from user;`
	data, err := mysql.FetchRowsD(sql)
	if err != nil{
		return nil,err
	}
	return data, nil
}

func QryAccInfo() ([]map[string]string, error){
	sql := `select * from account;`
	data, err := mysql.FetchRowsD(sql)
	if err != nil{
		return nil,err
	}
	return data, nil
}
