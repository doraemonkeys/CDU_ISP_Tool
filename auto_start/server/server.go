package server

import (
	"bufio"
	"fmt"
	"io"
	"ispTool_auto_start/utils"
	"log"
	"os"
	"strings"
	"time"
)

//启动打卡程序主体
func StartNewProgram() error {
	path, err := utils.GetCurrentPath()
	if err != nil {
		log.Println("获取当前文件目录失败！", err)
		return err
	}
	log.Println("获取到当前程序路径:", path)
	lastindex := strings.LastIndex(path, "\\")
	path = path[:lastindex]
	lastindex = strings.LastIndex(path, "\\")
	path = path[:lastindex]
	//path = strings.Replace(path, `\`, `\\`, -1)
	_, err = utils.Cmd_NoWait(path, []string{"powershell", "/c", "start", ".\\*.exe"})
	if err != nil {
		log.Println("启动打卡程序失败！", err)
		return err
	}
	return nil
}

//检查用户信息配置文件是否存在
func ConfigFileExist() bool {
	config, err := os.Open("./配置文件.config")
	//检查文件是否为空
	if err == nil {
		defer config.Close()
		temp := make([]byte, 20)
		n, err := config.Read(temp)
		if err != nil && err != io.EOF {
			log.Println("预读取配置文件失败，Error:", err)
			fmt.Println("预读取配置文件失败，Error:", err)
			return false
		}
		if n < 10 {
			log.Println("配置文件为空!")
			fmt.Println("配置文件为空!")
			return false
		}
		return true
	}
	log.Println("配置文件不存在或打开失败", err)
	fmt.Println("配置文件不存在或打开失败")
	return false
}

//是否已经设置为自启动
func CheckAutoStart() bool {
	autoStart, err := os.Open("./auto_start.config")
	if err != nil {
		log.Println("读取自启动信息失败！将默认设置为不自启动。", err)
		autoStart, err := os.OpenFile("./auto_start.config", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		if err != nil {
			log.Println("初始化自启动信息失败！", err)
			return false
		}
		defer autoStart.Close()
		//默认设置为不自启动
		autoStart.WriteString("false\n")
		return false
	}
	defer autoStart.Close()
	reader := bufio.NewReader(autoStart)
	auto, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			autoStart.WriteString("false\n")
		}
		log.Println("读取自动打卡信息失败！", err)
		return false
	}
	if strings.TrimSpace(auto) == "true" {
		return true
	}
	if strings.TrimSpace(auto) == "false" {
		return false
	}
	return false
}

//今日自动打卡是否成功,请确保auto_start.config文件存在
func TodayClockInSuccess() bool {
	clockInInfo, err := utils.ReverseRead("./auto_start.config", 2)
	if err != nil {
		log.Println("读取自动打卡信息失败！", err)
		fmt.Println("读取自动打卡信息失败！", err)
		return false
	}
	if len(clockInInfo) == 1 {
		return clockInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡成功"
	}
	if strings.TrimSpace(clockInInfo[0]) == "" {
		return clockInInfo[1] == time.Now().Format("2006/01/02")+" 自动打卡成功"
	}
	return clockInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡成功"
}

//今日自动打卡是否存在失败记录
func FailedLogExist() bool {
	clockInInfo, err := utils.ReverseRead("./auto_start.config", 2)
	if err != nil {
		log.Println("读取自动打卡信息失败！", err)
		fmt.Println("读取自动打卡信息失败！", err)
		return false
	}
	if len(clockInInfo) == 1 {
		return clockInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡失败"
	}
	if strings.TrimSpace(clockInInfo[0]) == "" {
		return clockInInfo[1] == time.Now().Format("2006/01/02")+" 自动打卡失败"
	}
	return clockInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡失败"
}
