package model

import
(
	"red-envelope/databases/mysql"
)


func QryUserPermission(uid string) ([]map[string]string, error){
	sql := `
    select u.id,p.*
    from red_user u
    left join red_user_role ur on u.id = ur.user_id
    left join red_role r on ur.role_id = r.id
    left join red_role_permission rp on r.id = rp.role_id
    left join red_permission p on rp.permission_id = p.id
    where u.user_name = ?
    `
	data, err := mysql.FetchRowsD(sql, uid)
	if err != nil{
		return nil,err
	}
	return data, nil
}
