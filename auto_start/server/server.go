package server

import (
	"bufio"
	"fmt"
	"io"
	"ispTool_auto_start/model"
	"ispTool_auto_start/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Doraemonkeys/lanzou"
	"github.com/fatih/color"
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

//今日自动打卡失败次数
func FailedTimes() (int, error) {
	count := 0
	for i := 1; ; i++ {
		oneLine, err := utils.ReadStartWithLastLine("./auto_start.config", i)
		if err == io.EOF {
			return count, nil
		}
		if err != nil {
			log.Println("读取自动打卡信息失败！", err)
			fmt.Println("读取自动打卡信息失败！", err)
			return -1, err
		}
		if strings.TrimSpace(oneLine) == "" {
			continue
		}
		if strings.TrimSpace(oneLine) == time.Now().Format("2006/01/02")+" 自动打卡成功" {
			return 0, nil
		}
		if strings.TrimSpace(oneLine) == time.Now().Format("2006/01/02")+" 自动打卡失败" {
			count++
			continue
		}
		if strings.TrimSpace(oneLine) == utils.GetYesterday().Format("2006/01/02")+" 自动打卡失败" {
			return count, nil
		}
		if strings.TrimSpace(oneLine) == utils.GetYesterday().Format("2006/01/02")+" 自动打卡成功" {
			return count, nil
		}
	}
}

//检查是否有更新,有更新则直接更新
func CheckUpdate() {
	//删除旧的更新文件
	os.Remove("./update.bat")
	var updateInfo model.Update
	updateInfo, err := utils.GetUpdateInfo()
	if err != nil {
		return
	}
	if utils.CompareVersion(updateInfo.AutoStartProgramVersion, model.Version) != 1 {
		color.Green("当前版本为最新版本！")
		return
	}
	//有更新
	log.Println("检测到新版本,正在尝试更新!")
	err = Update(updateInfo)
	if err == nil {
		log.Println("更新成功！")
		os.Exit(0)
	}
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
		color.Red("获取当前路径失败！")
		return err
	}
	//获取文件的绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Println("获取当前路径失败！", err)
		color.Red("获取当前路径失败！")
		return err
	}
	absPath = filepath.Dir(absPath)
	//absPath = strings.Replace(absPath, "\\", "/", -1)
	programName := filepath.Base(path)
	//命令1
	cmd1 := "del " + programName
	//命令2
	cmd2 := "rename " + tempName + " " + programName
	//命令3
	cmd3 := "cmd /c start " + `C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup\isp_auto_start.vbs`
	f, err := os.Create("update.bat")
	if err != nil {
		log.Println("创建批处理文件失败！", err)
		color.Red("创建批处理文件失败！")
		return err
	}
	_, err = f.WriteString("ping -n 2 127.1>nul" + " & " + cmd1 + " & " + cmd2 + " & " + cmd3 + " & exit")
	if err != nil {
		log.Println("写入批处理文件失败！", err)
		color.Red("写入批处理文件失败！")
		return err
	}
	err = utils.CmdNoOutput(absPath, []string{"cmd /c start .\\update.bat", "&", "exit"})
	if err != nil {
		log.Println("更新失败！", err)
		color.Red("更新失败！")
		return err
	}
	return nil
}

func checkFile(tempName string, updateInfo model.Update) error {
	md5, err := utils.GetFileMd5(tempName)
	if err != nil {
		log.Println("获取更新文件MD5失败！", err)
		color.Red("获取更新文件MD5失败！")
		return err
	}
	if md5 != updateInfo.AutoStartProgramMd5 {
		log.Println("更新文件MD5校验失败！")
		color.Red("更新文件MD5校验失败！")
		return err
	}
	return nil
}

func downloadUpdate(updateInfo model.Update) (string, error) {
	tempName := "temp.exe"
	if updateInfo.AutoStartProgramDirectUrl != "" {
		err := utils.DownloadFile(updateInfo.AutoStartProgramDirectUrl, tempName)
		if err != nil {
			log.Println("下载更新文件失败！", err)
			color.Red("下载更新文件失败！")
			return "", err
		}
	}
	if updateInfo.AutoStartProgramDirectUrl == "" {
		directUrl, err := lanzou.GetDownloadUrl(updateInfo.LanzouUrl, updateInfo.LanzouPwd, updateInfo.AutoStartProgramName)
		if err != nil {
			log.Println("获取更新文件下载地址失败！", err)
			color.Red("获取更新文件下载地址失败！")
			return "", err
		}
		err = lanzou.Download(directUrl, tempName)
		if err != nil {
			log.Println("下载更新文件失败！", err)
			color.Red("下载更新文件失败！")
			return "", err
		}
	}
	return tempName, nil
}
