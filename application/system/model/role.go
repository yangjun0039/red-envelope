package model

import "time"

type Role struct {
	Id          int
	Name        string
	//Permissions []int
	CreatedAt   time.Time
	UpdateAt    time.Time
}


