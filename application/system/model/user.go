package model

import "time"

type User struct {
	Id        int
	UserName  string
	Password  string
	//Roles     []int
	CreatedAt time.Time
	UpdatedAt time.Time
}

//func QryUserPermis(uid int) []Permission {
//	var roles, permisId []int
//	var permissions []Permission
//	for _, u := range (Users) {
//		if uid == u.Id {
//			roles = u.Roles
//			break
//		}
//	}
//	for _, role := range (Roles) {
//		for _, r := range (roles) {
//			if r == role.Id {
//				permisId = append(permisId, role.Permissions...)
//			}
//		}
//	}
//
//	for
//
//	return permissions
//}
