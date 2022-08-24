package engine

import (
	"ISP_Tool/fetcher"
	"ISP_Tool/login"
	"ISP_Tool/model"
	"ISP_Tool/uploader"
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
	userLocation, err := fetcher.GetLocation(user_no, &client)
	if err != nil {
		fmt.Println("ID", user.UserID, "获取地理位置失败！")
		return err
	}
	user.Location = userLocation
	err = uploader.ISP_Clock_In(&client, user)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	key_value, err := fetcher.CheckingAnomalies(user_no, &client)
	if err != nil {
		log.Println("健康登记已成功,但出现异常，将自动撤回本次打卡！")
		color.Red("健康登记已成功,但出现异常，将自动撤回本次打卡！")
		err2 := uploader.Cancel_Clock_In(key_value, &client)
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
