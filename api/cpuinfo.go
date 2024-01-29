package api

import (
	"gin_exercise/config"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

func CpuinfoHandler(g *gin.Context) {
	client, err := connect()
	if err != nil {
		g.JSON(200, gin.H{
			"code":  2,
			"error": err,
		})
		return 
	}
	defer client.Close()
	cpuUsage, err := runCommand(client, "top -bn1 | grep '^%Cpu' | awk '{print $2}'")
	if err != nil {
		log.Fatalf("Failed to run cpu usage command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err,
		})
	}
	g.JSON(200, gin.H{
		"code":     0,
		"cpuUsage": cpuUsage,
	})
}

func connect() (*ssh.Client, error) {
	//从config获取信息
	var host = config.GlobalConfig.Host
	var port = config.GlobalConfig.Port
	var user = config.GlobalConfig.User
	var password = config.GlobalConfig.Password
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
		return nil, err
	} else {
		return client, nil
	}

}

func runCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
