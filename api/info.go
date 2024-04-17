package api

import (
	"gin_exercise/config"
	"gin_exercise/controller"
	"gin_exercise/dao"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

var SshConnect *ssh.Client

func CpuinfoHandler(g *gin.Context) {
	client := SshConnect
	// 修改了原始命令，以提取用户空间使用(us)、系统空间使用(sy)和CPU空闲(id)的值
	command := "top -bn1 | grep '^%Cpu' | awk '{print $2, $4, $8}'"
	cpuUsage, err := dao.RunCommand(client, command)
	if err != nil {
		seelog.Info("Failed to run cpu usage command: %s", err) // 使用Printf代替Fatalf，避免因错误而终止服务
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
		seelog.Info("Unexpected format of cpu usage data: %s", cpuUsage)
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
	name, _ := g.Get("name")

	// 返回JSON数据，包括用户空间使用、系统空间使用和CPU空闲的值
	g.JSON(200, gin.H{
		"name":        name,
		"code":        0,
		"userUsage":   userUsage,
		"systemUsage": systemUsage,
		"idleUsage":   idle,
	})
}

func MemoryinfoHandler(g *gin.Context) {
	client := SshConnect
	// 运行命令以获取内存总量和已使用内存的值
	command := "free -m | awk 'NR==2{print $2,$3}'"
	memoryUsage, err := dao.RunCommand(client, command)
	if err != nil {
		seelog.Info("Failed to run memory usage command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}

	// 分割memoryUsage字符串以获取单独的值
	usageParts := strings.Fields(memoryUsage)
	if len(usageParts) < 2 {
		seelog.Info("Unexpected format of memory usage data: %s", memoryUsage)
		g.JSON(200, gin.H{
			"code":  1,
			"error": "unexpected format of memory usage data",
		})
		return
	}

	// 将字符串值转换为浮点数
	totalMemory, err := strconv.ParseFloat(usageParts[0], 64)
	if err != nil {
		log.Printf("Error parsing total memory to float: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": "error parsing total memory",
		})
		return
	}
	usedMemory, err := strconv.ParseFloat(usageParts[1], 64)
	if err != nil {
		log.Printf("Error parsing used memory to float: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": "error parsing used memory",
		})
		return
	}

	// 计算内存使用率
	memoryUsageRate := (usedMemory / totalMemory) * 100

	// 返回JSON数据，包括内存总量和内存使用率的值
	g.JSON(200, gin.H{
		"code":        0,
		"totalMemory": totalMemory,
		"memoryUsage": math.Round(memoryUsageRate*100) / 100,
	})
}

func DetailedGPUInfoHandler(g *gin.Context) {
	client := SshConnect
	// 修改命令以获取型号名称，唯一标识符，总内存大小，NVIDIA 驱动版本，电源使用限制
	command := `nvidia-smi --query-gpu=name,uuid,memory.total,driver_version,power.limit --format=csv,noheader,nounits`
	gpuInfo, err := dao.RunCommand(client, command)
	if err != nil {
		seelog.Error("Failed to run detailed GPU info command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}
	infoLines := strings.Split(strings.TrimSpace(gpuInfo), "\n")
	gpus := make([]gin.H, len(infoLines))
	for i, line := range infoLines {
		parts := strings.Split(line, ", ")
		if len(parts) < 5 {
			seelog.Error("Unexpected format of detailed GPU info data: %s", line)
			continue
		}
		// 由于内存和电源限制已经是以正确的单位返回，这里不需要转换
		memoryTotal, _ := strconv.ParseInt(parts[2], 10, 64)
		powerLimit, _ := strconv.ParseFloat(parts[4], 64)
		gpus[i] = gin.H{
			"modelName":     parts[0],
			"uuid":          parts[1],
			"memoryTotalMB": memoryTotal,
			"driverVersion": parts[3],
			"powerLimitW":   powerLimit,
		}
	}
	// 返回JSON数据，包括每个GPU的型号名称，唯一标识符，总内存大小，驱动版本和电源使用限制
	g.JSON(200, gin.H{
		"code": 0,
		"gpus": gpus,
	})
}

func GpuinfoHandler(g *gin.Context) {
	client := SshConnect
	// 使用nvidia-smi命令获取温度(Temp)，功率使用(Pwr:Usage)和GPU利用率(GPU-Util)的信息
	// 这个命令的输出需要根据实际输出进行适配调整
	command := `nvidia-smi --query-gpu=temperature.gpu,utilization.gpu,power.draw --format=csv,noheader,nounits`
	gpuInfo, err := dao.RunCommand(client, command)
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
	client := SshConnect
	// 获取uname信息
	unameCommand := "uname -a"
	unameOutput, err := dao.RunCommand(client, unameCommand)
	if err != nil {
		seelog.Error("Failed to run uname command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}
	// 获取lsb_release信息
	lsbReleaseCommand := "lsb_release -a"
	lsbReleaseOutput, err := dao.RunCommand(client, lsbReleaseCommand)
	if err != nil {
		seelog.Error("Failed to run lsb_release command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}
	// 解析uname输出
	unameParts := strings.Fields(unameOutput)
	if len(unameParts) < 6 {
		seelog.Error("Unexpected format of uname data: ", unameOutput)
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

func Whoinfohandler(g *gin.Context) {
	client := SshConnect
	// 获取当前登录用户的信息
	command := "who"
	whoOutput, err := dao.RunCommand(client, command)
	if err != nil {
		log.Printf("Failed to run who command: %s", err)
		g.JSON(200, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}
	// 解析who输出
	whoLines := strings.Split(whoOutput, "\n")
	users := make([]gin.H, 0)
	for _, line := range whoLines {
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		users = append(users, gin.H{
			"user":     parts[0],
			"terminal": parts[1],
			"date":     parts[2],
			"time":     parts[3],
		})
	}
	// 返回当前登录用户的信息
	g.JSON(200, gin.H{
		"code":  0,
		"users": users,
	})
}

func WinfoHandler(g *gin.Context) {

	client := SshConnect
	command := "cat /proc/cpuinfo| grep 'physical id'| sort| uniq| wc -l"
	cpunum, _ := dao.RunCommand(client, command)
	cpunum = strings.TrimSpace(cpunum)
	cpunums, _ := strconv.Atoi(cpunum)

	command = "grep -m 1 'cpu cores' /proc/cpuinfo | awk '{print $4}'"
	cpucorenum, _ := dao.RunCommand(client, command)
	cpucorenum = strings.TrimSpace(cpucorenum)
	cpucorenums, _ := strconv.Atoi(cpucorenum)
	idealload := float64(cpunums*cpucorenums) * 0.7

	command = "cat /proc/loadavg"
	avgload, _ := dao.RunCommand(client, command)
	avgloads := strings.Fields(avgload)
	load1min := avgloads[0]
	load5min := avgloads[1]
	load15min := avgloads[2]

	command = "w"
	wOutput, _ := dao.RunCommand(client, command)
	lines := strings.Split(wOutput, "\n")

	users := make([]gin.H, 0)
	for _, line := range lines[2:] { // 从第三行开始是用户信息
		parts := strings.Fields(line)
		if len(parts) < 8 {
			continue
		}
		userInfo := gin.H{
			"user": parts[0],
			"from": parts[2],
			"pcpu": parts[6],
			"what": strings.Join(parts[7:], " "), // 用户正在运行的命令
		}
		users = append(users, userInfo)
	}

	g.JSON(200, gin.H{
		"load1min":   load1min,
		"load5min":   load5min,
		"load15min":  load15min,
		"cpunum":     cpunums,
		"cpucorenum": cpucorenums,
		"idealload":  idealload,
		"users":      users,
		"code":       0,
	})
}

func LstmHandler(g *gin.Context) {
	requestData :=controller.RequestData{}
	g.Bind(&requestData)
	seelog.Info("data:", requestData)
	g.JSON(200, gin.H{
		"code":   0,
		"output": requestData.Data,
	})
	
	// client := SshConnect
	// command := "cd project && /home/ycx/anaconda3/envs/pytorch/bin/python /home/ycx/project/remote.py"
	// session, _ := client.NewSession()

	// defer session.Close()

	// // 创建一个stdin管道用来发送数据到远程Python脚本
	// stdin, _ := session.StdinPipe()

	// // 从远程Python脚本接收输出
	// var b bytes.Buffer
	// session.Stdout = &b

	// // 运行远程Python脚本
	// session.Start(command)

	// // 向Python脚本发送数据
	// fmt.Fprintln(stdin, data)
	// stdin.Close()

	// // 等待脚本执行结束
	// session.Wait()

	

}
