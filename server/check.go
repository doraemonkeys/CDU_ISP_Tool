package server

import (
	"ISP_Tool/model"
	"ISP_Tool/utils"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Doraemonkeys/lanzou"
	"github.com/fatih/color"
)

// 今日自动打卡是否成功,请确保auto_start.config文件存在,不存在可以调用CheckAutoStart()函数创建。
func TodayCheckInSuccess() bool {
	checkInInfo, err := utils.ReverseRead("./config/auto_start.config", 2)
	if err != nil {
		log.Println("读取自动打卡信息失败！", err)
		fmt.Println("读取自动打卡信息失败！", err)
		return false
	}
	if len(checkInInfo) == 1 {
		if checkInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡成功" {
			return true
		}
	}
	if strings.TrimSpace(checkInInfo[0]) == "" {
		return checkInInfo[1] == time.Now().Format("2006/01/02")+" 自动打卡成功"
	}
	if checkInInfo[0] == time.Now().Format("2006/01/02")+" 自动打卡成功" {
		return true
	}
	return false
}

// 是否已经设置为自启动
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

// 在打卡后的界面寻找异常关键字
func LookForKeyword(user model.UserInfo, content []byte) error {
	prefixRe := user.Province + `[ ]*\|` + user.City + `[ ]*\|` + user.Area
	re3 := regexp.MustCompile(prefixRe + model.All.Regexp.Today_statusRe)
	Today_status := re3.Find(content)
	if Today_status == nil {
		log.Println("检查今日打卡状态出错，可能是ISP结构发生改变！")
		//可能是第一次打卡,进行全局匹配是否出现异常
		re4 := regexp.MustCompile(prefixRe)
		loc := re4.FindIndex([]byte(content))
		if loc == nil {
			color.Red("无法检查打卡是否异常,建议去ISP看看！")
			return nil
		}
		Today_status = content[loc[0]:]
	}
	//下面匹配两次关键字(冗余操作防止意外)
	re4 := regexp.MustCompile(model.All.Regexp.AbnormalRe)
	match4 := re4.Find(Today_status)
	if match4 != nil {
		log.Println("检测到 异常 关键字！")
		return errors.New("健康登记出现异常")
	}
	re5 := regexp.MustCompile(model.All.Regexp.AbnormalColorRe)
	match5 := re5.FindAll(Today_status, -1)
	if len(match5) > 0 { //删除按钮也是红色的
		log.Println("检测到 " + model.All.Regexp.AbnormalColorRe + " 关键字！")
		return errors.New("健康登记出现异常")
	}
	re6 := regexp.MustCompile(`正常`)
	match6 := re6.FindAll(Today_status, -1)
	if len(match6) < 4 {
		log.Println("检测到 " + `正常` + " 关键字 少于4个！")
		return errors.New("健康登记出现异常")
	}
	//打卡无异常
	return nil
}

// 检查是否有更新,有更新则直接更新
func CheckUpdate() {
	//删除旧的更新文件
	os.Remove("./update.bat")
	var updateInfo model.Update
	updateInfo, err := utils.GetUpdateInfo()
	if err != nil {
		return
	}
	if utils.CompareVersion(updateInfo.MainProgramVersion, model.Version) != 1 {
		color.Green("当前版本为最新版本！")
		return
	}
	//有更新
	color.Yellow("检测到新版本 %s,正在尝试更新...", updateInfo.MainProgramVersion)
	log.Println("检测到新版本", updateInfo.MainProgramVersion, ",正在尝试更新!")
	ctx, cancel := context.WithCancel(context.Background())
	go utils.WaitAnimation(ctx)
	err = Update(updateInfo)
	cancel()
	fmt.Println()
	if err == nil {
		log.Println("更新成功！")
		color.Green("更新成功！程序即将退出！")
		os.Exit(0)
	}
	color.Red(err.Error())
}

func Update(updateInfo model.Update) error {
	//下载更新文件
	tempName, err := downloadUpdate(updateInfo)
	if err != nil {
		return err
	}
	//校验文件
	err = checkFile(tempName, updateInfo)
	if err != nil {
		return err
	}
	return updateAndRestart(tempName)
}

func updateAndRestart(tempName string) error {
	//获取文件路径
	path, err := utils.GetCurrentPath()
	if err != nil {
		log.Println("获取当前路径失败！", err)
		return fmt.Errorf("获取当前路径失败,%w", err)
	}
	//获取文件的绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Println("获取当前路径失败！", err)
		return fmt.Errorf("获取当前路径失败,%w", err)
	}
	absPath = filepath.Dir(absPath)
	programName := filepath.Base(path)
	//命令1
	cmd1 := "del " + programName
	//命令2
	cmd2 := "rename " + tempName + " " + programName
	//命令3
	cmd3 := "cmd /c start " + programName
	batFile := "update.bat"
	//命令4
	cmd4 := "del " + batFile
	f, err := os.Create(batFile)
	if err != nil {
		log.Println("创建批处理文件失败！", err)
		return fmt.Errorf("创建批处理文件失败,%w", err)
	}
	_, err = f.WriteString(`if "%1" == "h" goto begin
	mshta vbscript:createobject("wscript.shell").run("""%~nx0"" h",0)(window.close)&&exit
	:begin` + "\n")
	if err != nil {
		log.Println("写入批处理文件失败！", err)
		return fmt.Errorf("写入批处理文件失败,%w", err)
	}
	f.WriteString("ping -n 2 127.1>nul" + " & " + cmd1 + " & " + cmd2 + " & " + cmd3 + " & " + cmd4)
	f.Close()
	cmdStr := "cmd /c start .\\" + batFile
	err = utils.CmdNoOutput(absPath, []string{cmdStr, "&", "exit"})
	if err != nil {
		log.Println("执行cmd命令失败！", err)
		return fmt.Errorf("执行cmd命令失败,%w", err)
	}
	return nil
}

func checkFile(tempName string, updateInfo model.Update) error {
	md5, err := utils.GetFileMd5(tempName)
	if err != nil {
		log.Println("获取更新文件MD5失败！", err)
		return fmt.Errorf("获取更新文件MD5失败,%w", err)
	}
	md5 = strings.ToUpper(md5)
	updateInfo.MainProgramMd5 = strings.ToUpper(updateInfo.MainProgramMd5)
	if md5 != updateInfo.MainProgramMd5 {
		log.Println("更新文件MD5校验失败！")
		return errors.New("更新文件MD5校验失败")
	}
	return nil
}

func downloadUpdate(updateInfo model.Update) (string, error) {
	tempName := "temp"
	if updateInfo.MainProgramDirectUrl != "" {
		err := utils.DownloadFile(updateInfo.MainProgramDirectUrl, tempName)
		if err != nil {
			log.Println("下载更新文件失败！", err)
			return "", fmt.Errorf("下载更新文件失败,%w", err)
		}
	}
	if updateInfo.MainProgramDirectUrl == "" {
		directUrl, err := lanzou.GetDownloadUrl(updateInfo.LanzouUrl, updateInfo.LanzouPwd, updateInfo.MainProgramName)
		if err != nil {
			log.Println("获取更新文件下载地址失败！", err)
			return "", fmt.Errorf("获取更新文件下载地址失败,%w", err)
		}
		err = lanzou.Download(directUrl, tempName)
		if err != nil {
			log.Println("下载更新文件失败！", err)
			return "", fmt.Errorf("下载更新文件失败,%w", err)
		}
	}
	return tempName, nil
}
