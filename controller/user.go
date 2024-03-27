package controller

type User struct {
	Username string
	Password string
}

type Userssh struct {
	Username string
	Host     string
	Port     string
	User     string
	Password string
}

type Mytoken struct {
	Token string
}

func (User) TableName() string {
	return "user_info" //自定义表名
}


func (Userssh) TableName() string {
	return "user_ssh_info" //自定义表名
}