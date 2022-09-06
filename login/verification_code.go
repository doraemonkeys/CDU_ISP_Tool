package login

import (
	"ISP_Tool/model"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/fatih/color"
)

func Get_ISP_Login_code(client *http.Client) (string, error) {

	req, _ := http.NewRequest("GET", model.All.Login.LoginWebUrl, nil)
	req.Header.Set("User-Agent", model.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		re := regexp.MustCompile(model.All.Regexp.Ipv6Re)
		if re.FindString(err.Error()) != "" {
			log.Println("访问ISP登录界面失败！ISP不支持Ipv6。")
			color.Red("访问ISP登录界面失败！ISP不支持Ipv6。")
		} else {
			log.Println("访问ISP登录界面失败！可能是ISP结构发生变化，请联系开发者。")
			fmt.Println("访问ISP登录界面失败！可能是ISP结构发生变化，请联系开发者。")
		}
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取登录页面失败！", err)
		fmt.Println("读取登录页面失败！", err)
		return "", err
	}
	//fmt.Println(string(content))
	re := regexp.MustCompile(model.All.Regexp.VerificationCodeRe)
	substr := re.FindSubmatch(content)
	if substr == nil {
		log.Println("获取登录验证码失败！可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("获取登录验证码失败！可能是ISP结构发生变化，请联系开发者。")
		//将页面内容写入到文件用于debug
		ioutil.WriteFile("loginError.html", content, 0644)
		return "", errors.New("substr == nil")
	}
	return string(substr[1]), nil
}
