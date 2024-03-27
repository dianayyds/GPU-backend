package config

import (
	"encoding/json"
	"fmt"
	"gin_exercise/controller"
	"gin_exercise/help"
	"gin_exercise/mydb"
	"os"

	"github.com/cihub/seelog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	DBName   string `json:"dbname"`
}

var GlobalConfig *Config

type Config struct {
	ListenPort int
	User       string
	Password   string
	Host       string
	Port       string
}

func init() {
	GlobalConfig = GetConfig()
}

func GetConfig() *Config {
	return &Config{
		ListenPort: 8080,
		//ssh服务器用户名，密码，ip地址，端口号
		User:     "ycx",
		Password: "20231105",
		Host:     "211.71.76.205",
		Port:     "22",
	}
}

func Initlog() {
	help.SetupLogger()
	seelog.Info(fmt.Sprintf("Begin Seelog"))
}

func LoadConfig() (map[string]interface{}, error) {
	content, err := os.ReadFile("Usersconfig.json")
	if err != nil {
		seelog.Error(err)
		// fmt.Println("read file error")
		return nil, err
	}
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		seelog.Error(err)
		return nil, err
	}
	return payload, nil
}

func InitUserdatabase() {
	config, _ := LoadConfig()
	username := config["username"].(string)
	password := config["password"].(string)
	address := config["address"].(string)
	dbname := config["dbname"].(string)
	dsn := username + ":" + password + "@tcp(" + address + ")/" + dbname + "?charset=utf8&parseTime=True&loc=Local"
	userDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		seelog.Error(err)
	} else {
		mydb.UserDB = userDB
	}
	p := mydb.UserDB.Migrator().HasTable(&controller.User{}) //检测是否存在users表单
	if !p {
		seelog.Info(fmt.Sprintf("create table users"))
		mydb.UserDB.Migrator().CreateTable(&controller.User{}) //不存在则创建
	} else {
		seelog.Info(fmt.Sprintf("table users already exists"))
	}
	p = mydb.UserDB.Migrator().HasTable(&controller.Userssh{}) //检测是否存在userssh表单
	if !p {
		seelog.Info(fmt.Sprintf("create table userssh"))
		mydb.UserDB.Migrator().CreateTable(&controller.Userssh{}) //不存在则创建
	} else {
		seelog.Info(fmt.Sprintf("table userssh already exists"))
	}
}
