package controller

type User struct {
	Username string
	Password string
}

type Mytoken struct {
	Token string
}

func (User) TableName() string {
	return "user_info" //自定义表名
}

type Database struct {
	Ip           string
	Port         string
	DatabaseName string
	Username     string
	Password     string
}
