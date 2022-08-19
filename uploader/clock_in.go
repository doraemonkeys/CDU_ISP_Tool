package uploader

import (
	"ISP_Tool/model"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

func ISP_Clock_In(client *http.Client, user model.UserInfo) error {
	today := time.Now().Local().Format("2006年1月2日")
	apiUrl := "https://xsswzx.cdu.edu.cn/ispstu/com_user/projecthealth_add.asp"
	// URL param
	queryData := url.Values{}
	queryData.Set("id", user.UserNonce)
	queryData.Set("id2", today)
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
		log.Printf("parse url requestUrl failed, err:%v\n", err)
		return err
	}
	u.RawQuery = queryData.Encode() // URL encode
	// 构造请求
	param := url.Values{}
	param.Set("action", "add")
	param.Set("area", user.Area)
	param.Set("city", user.City)
	param.Set("province", user.Province)
	param.Set("fare", "否")
	param.Set("kesou", "否")
	param.Set("wls", "否")
	param.Set("wuhan", "否")
	data := param.Encode()
	req4, _ := http.NewRequest("POST", u.String(), strings.NewReader(data))
	req4.Header.Set("authority", "xsswzx.cdu.edu.cn")
	req4.Header.Set("content-type", "application/x-www-form-urlencoded")
	req4.Header.Set("user-agent", model.UserAgent)

	resp4, err := client.Do(req4)
	if err != nil {
		log.Println("发起ISP登记请求失败！可能是ISP结构发生变化，请联系开发者。", err)
		fmt.Println("发起ISP登记请求失败！可能是ISP结构发生变化，请联系开发者。", err)
		return err
	}
	content4, err := ioutil.ReadAll(resp4.Body)
	if err != nil {
		log.Println("读取ISP登记返回信息失败！", err)
		fmt.Println("读取ISP登记返回信息失败！", err)
		return err
	}
	re := regexp.MustCompile("提交成功")
	match := re.Find(content4)
	if match != nil {
		return nil
	}
	re = regexp.MustCompile("登记已存在")
	match = re.Find(content4)
	if match != nil {
		log.Println(user.UserID, "健康登记打卡已存在")
		color.Green("%s %s", user.UserID, "健康登记打卡已存在")
		return errors.New("健康登记打卡已存在")
	}
	return errors.New("CDU-ISP 健康登记打卡 失败")
}
