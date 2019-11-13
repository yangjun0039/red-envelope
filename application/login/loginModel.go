package login

import (
	"red-envelope/databases/mysql"
)

func QryUserInfo(uid string) (map[string]string, error) {
	sql := `select * from red_user where user_name = ?;`
	data, err := mysql.FetchRowD(sql, uid)
	if err != nil {
		return nil, err
	}
	return data, nil
}
