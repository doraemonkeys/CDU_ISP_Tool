package engine

import (
	"ISP_Tool/fetcher"
	"ISP_Tool/login"
	"ISP_Tool/model"
	"ISP_Tool/uploader"
	"ISP_Tool/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

func Run(client http.Client, user model.UserInfo) error {
	err := login.LoginISP(&client, user)
	if err != nil {
		fmt.Println("ID", user.UserID, "登录失败！")
		return err
	}
	log.Println("登录ISP成功")
	fmt.Println("登录ISP成功")
	user_no, err := fetcher.Get_User_Nonce(&client)
	if err != nil {
		fmt.Println("ID", user.UserID, "获取ISP的用户识别码失败！")
		return err
	}
	fmt.Println("从isp获取到用户识别码：", user_no)
	log.Println("从isp获取到用户识别码：", user_no)
	user.UserNonce = user_no
	//选择打卡地址
	err = selectAddress(&user, &client)
	if err != nil {
		return err
	}
	err = uploader.ISP_CheckIn(&client, user)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	key_value, err := fetcher.CheckingAnomalies(user_no, &client)
	if err != nil {
		log.Println("健康登记已成功,但出现异常，将自动撤回本次打卡！")
		color.Red("健康登记已成功,但出现异常，将自动撤回本次打卡！")
		err2 := uploader.CancelCheckIn(key_value, &client)
		if err2 != nil {
			log.Println("自动撤回打卡失败！请前往ISP手动修改！")
			color.Red("自动撤回打卡失败！请前往ISP手动修改！")
		}
		if err2 == nil {
			log.Println("自动撤回打卡成功！")
			color.HiGreen("自动撤回打卡成功！")
		}
		return err
	}
	return nil
}

func selectAddress(user *model.UserInfo, client *http.Client) error {
	//若配置文件已设置地址，则优先使用配置文件地址
	if user.Location.Province == "" {
		userLocation, err := fetcher.GetLocation(*user, client)
		if err != nil {
			fmt.Println("ID", user.UserID, "获取地理位置失败！")
			return err
		}
		user.Location = userLocation
	} else {
		attributes := [5]color.Attribute{}
		attributes[0] = color.FgYellow
		utils.ColorPrint(attributes[:], "使用配置文件中设置的地址打卡", "，如果有错误请前往ISP手动修改！\n")
		log.Println("使用配置文件中设置的地址打卡，如果有错误请前往ISP手动修改！")
	}
	return nil
}
