package service

import "gin_exercise/controller"

type LoginDTO struct {
	Name  string
	Token string
}

type ShowUsersDTO struct {
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
}

type ShowUsersRetDTO struct {
	Usercnt int64
	Users   []controller.User
}
