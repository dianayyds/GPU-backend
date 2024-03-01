package api

import (
	"gin_exercise/config"
	"log"
	"strconv"
	"strings"

	"github.com/cihub/seelog"
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

func GpuinfoHandler(g *gin.Context) {
	client := config.SshConnect
	// 使用nvidia-smi命令获取温度(Temp)，功率使用(Pwr:Usage)和GPU利用率(GPU-Util)的信息
	// 这个命令的输出需要根据实际输出进行适配调整
	command := `nvidia-smi --query-gpu=temperature.gpu,utilization.gpu,power.draw --format=csv,noheader,nounits`
	gpuInfo, err := runCommand(client, command)
	if err != nil {
		seelog.Error("Failed to run GPU info command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}
	// 假设gpuInfo的格式为"70, 30, 160\n65, 40, 150"，每行代表一个GPU的温度，使用率，功率使用量
	// 分割gpuInfo字符串以获取单独的GPU信息
	infoLines := strings.Split(strings.TrimSpace(gpuInfo), "\n")
	gpus := make([]gin.H, len(infoLines))
	for i, line := range infoLines {
		parts := strings.Split(line, ", ")
		if len(parts) < 3 {
			seelog.Error("Unexpected format of GPU info data: %s", line)
			continue
		}
		// 将字符串值转换为浮点数
		temp, _ := strconv.ParseFloat(parts[0], 64)
		utilization, _ := strconv.ParseFloat(parts[1], 64)
		powerDraw, _ := strconv.ParseFloat(parts[2], 64)
		gpus[i] = gin.H{
			"temperature": temp,
			"utilization": utilization,
			"powerDraw":   powerDraw,
		}
	}
	// 返回JSON数据，包括每个GPU的温度、使用率和功率使用量
	g.JSON(200, gin.H{
		"code": 0,
		"gpus": gpus,
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
