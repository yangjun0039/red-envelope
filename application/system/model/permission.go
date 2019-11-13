package model

import "time"


type Permission struct {
	Id         int
	Name       string
	HttpMethod string
	HttpPath   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}


