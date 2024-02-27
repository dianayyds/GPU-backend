package api

import (
	"gin_exercise/config"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

func CpuinfoHandler(g *gin.Context) {
	client := config.SshConnect
	// 修改了原始命令，以提取用户空间使用(us)、系统空间使用(sy)和CPU空闲(id)的值
	command := "top -bn1 | grep '^%Cpu' | awk '{print $2, $4, $8}'"
	cpuUsage, err := runCommand(client, command)
	if err != nil {
		log.Printf("Failed to run cpu usage command: %s", err) // 使用Printf代替Fatalf，避免因错误而终止服务
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(), // 确保错误被适当地转换成字符串
		})
		return // 添加return，确保错误时不继续执行
	}

	// 假设cpuUsage是以空格分隔的字符串，例如："44.9 11.2 43.8"
	// 分割cpuUsage字符串以获取单独的值
	usageParts := strings.Fields(cpuUsage)
	if len(usageParts) < 3 {
		log.Printf("Unexpected format of cpu usage data: %s", cpuUsage)
		g.JSON(200, gin.H{
			"code":  1,
			"error": "unexpected format of cpu usage data",
		})
		return
	}

	// 将字符串值转换为浮点数，以便前端可以更灵活地处理这些值
	userUsage, err := strconv.ParseFloat(usageParts[0], 64)
	if err != nil {
		log.Printf("Error parsing user cpu usage to float: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": "error parsing user cpu usage",
		})
		return
	}
	systemUsage, err := strconv.ParseFloat(usageParts[1], 64)
	if err != nil {
		log.Printf("Error parsing system cpu usage to float: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": "error parsing system cpu usage",
		})
		return
	}
	idle, err := strconv.ParseFloat(usageParts[2], 64)
	if err != nil {
		log.Printf("Error parsing cpu idle to float: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": "error parsing cpu idle",
		})
		return
	}

	// 返回JSON数据，包括用户空间使用、系统空间使用和CPU空闲的值
	g.JSON(200, gin.H{
		"code":        0,
		"userUsage":   userUsage,
		"systemUsage": systemUsage,
		"idleUsage":   idle,
	})
}

func BaseinfoHandler(g *gin.Context) {
	client := config.SshConnect
	// 获取uname信息
	unameCommand := "uname -a"
	unameOutput, err := runCommand(client, unameCommand)
	if err != nil {
		log.Printf("Failed to run uname command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}
	// 获取lsb_release信息
	lsbReleaseCommand := "lsb_release -a"
	lsbReleaseOutput, err := runCommand(client, lsbReleaseCommand)
	if err != nil {
		log.Printf("Failed to run lsb_release command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}
	// 解析uname输出
	unameParts := strings.Fields(unameOutput)
	if len(unameParts) < 6 {
		log.Printf("Unexpected format of uname data: %s", unameOutput)
		g.JSON(200, gin.H{
			"code":  1,
			"error": "unexpected format of uname data",
		})
		return
	}
	// 解析lsb_release输出
	lsbReleaseLines := strings.Split(lsbReleaseOutput, "\n")
	lsbReleaseMap := make(map[string]string)
	for _, line := range lsbReleaseLines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				lsbReleaseMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
	// 构造并返回前端所需的数据
	g.JSON(200, gin.H{
		"code":            0,
		"operatingSystem": unameParts[0],                                              // 操作系统
		"hostname":        unameParts[1],                                              // 网络节点主机名
		"kernelVersion":   unameParts[2],                                              // 内核版本
		"cpuArchitecture": unameParts[len(unameParts)-2],                              // CPU架构
		"release":         lsbReleaseMap["Distributor ID"] + lsbReleaseMap["Release"], // 发行版版本
		"host":            config.GlobalConfig.Host,
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
