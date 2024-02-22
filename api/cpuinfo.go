package api

import (
	"gin_exercise/config"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

func CpuinfoHandler(g *gin.Context) {
	client:= config.SshConnect
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
