package server

import (
	"ISP_Tool/model"
	"ISP_Tool/utils"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

//今日自动打卡是否成功,请确保auto_start.config文件存在
func TodayClockInSuccess() bool {
	clockInInfo, err := utils.ReverseRead("./config/auto_start.config", 2)
	if err != nil {
		log.Println("读取自动打卡信息失败！", err)
		fmt.Println("读取自动打卡信息失败！", err)
		return false
	}
	if len(clockInInfo) == 1 {
		if clockInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡成功" {
			return true
		}
	}
	if strings.TrimSpace(clockInInfo[0]) == "" {
		return clockInInfo[1] == time.Now().Format("2006/01/02")+" 自动打卡成功"
	}
	if clockInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡成功" {
		return true
	}
	return false
}

//是否已经设置为自启动
func CheckAutoStart() bool {
	autoStart, err := os.Open("./config/auto_start.config")
	if err != nil {
		log.Println("读取自启动信息失败！", err)
		autoStart, err := os.OpenFile("./config/auto_start.config", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
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
		fmt.Println("读取自动打卡信息失败！", err)
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

//在打卡后的界面寻找异常关键字
func LookForKeyword(content []byte) error {
	re3 := regexp.MustCompile(model.Today_statusRe)
	Today_status := re3.Find(content)
	if Today_status == nil {
		//可能是第一次打卡,进行全局匹配是否出现异常
		re := regexp.MustCompile("异常")
		match := re.Find(content)
		if match == nil {
			//打卡无异常
			return nil
		}
		return errors.New("健康登记出现异常")
	}
	//下面匹配两次关键字(冗余操作防止意外)
	re4 := regexp.MustCompile("异常")
	match4 := re4.Find(Today_status)
	if match4 != nil {
		return errors.New("健康登记出现异常")
	}
	re5 := regexp.MustCompile("color=red")
	match5 := re5.Find(Today_status)
	if match5 != nil {
		return errors.New("健康登记出现异常")
	}
	//打卡无异常
	return nil
}
