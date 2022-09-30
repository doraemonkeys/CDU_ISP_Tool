package login

import (
	"ISP_Tool/model"
	"ISP_Tool/utils"
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/text/transform"
)

// 登陆客户端

func LoginISP(client *http.Client, user model.UserInfo) error {
	content, err := Get_Login_Page(client, user)
	if err != nil {
		return err
	}
	code, err := Get_ISP_Login_code(content)
	if err != nil {
		return err
	}
	fmt.Println("成功获取登录验证码:", code)
	log.Println("成功获取登录验证码:", code)
	parame, err := getPostField(content)
	if err != nil {
		return err
	}
	parame.Set(model.All.Login.Input1Field, user.UserID)
	parame.Set(model.All.Login.Input2Field, user.UserPwd)
	parame.Set(model.All.Login.Input3Field, code)
	// parame := url.Values{}
	// parame.Set(model.All.Login.Input1Field, user.UserID)
	// parame.Set(model.All.Login.Input2Field, user.UserPwd)
	// parame.Set(model.All.Login.Input3Field, code)
	// for _, v := range model.All.Login.Other {
	// 	parame.Add(v.Field, v.Value)
	// }
	data := parame.Encode()
	var request *http.Request
	if model.UseVPN {
		request, _ = http.NewRequest(model.All.Login.Head.Method,
			model.All.VPN_ISP_BaseURL+model.All.Login.LoginUrl, strings.NewReader(data))
	} else {
		request, _ = http.NewRequest(model.All.Login.Head.Method,
			model.All.DirectBaseURL+model.All.Login.LoginUrl, strings.NewReader(data))
	}
	request.Header.Set("authority", model.All.Login.Head.Authority)
	request.Header.Set("content-type", model.All.Login.Head.Content_type)
	request.Header.Set("user-agent", model.UserAgent)
	if model.UseVPN {
		request.Header.Set("referer", model.All.VPN_ISP_BaseURL+model.All.Login.Head.Referer)
	} else {
		request.Header.Set("referer", model.All.DirectBaseURL+model.All.Login.Head.Referer)
	}
	// 发起登录请求
	resp, err := client.Do(request)
	if err != nil {
		log.Println("发起ISP登录请求失败！可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("发起ISP登录请求失败！可能是ISP结构发生变化，请联系开发者。")
		return err
	}
	bodyReader := bufio.NewReader(resp.Body)
	//自动检测html编码
	e, err := utils.DetermineEncodingbyPeek(bodyReader)
	if err != nil {
		log.Println("登录返回界面检测html编失败，请联系开发者。", err)
		fmt.Println("登录返回界面检测html编失败，请联系开发者。", err)
		return err
	}
	//转码utf-8
	utf8BodyReader := transform.NewReader(bodyReader, e.NewDecoder())

	content, err = ioutil.ReadAll(utf8BodyReader)
	if err != nil {
		log.Println("读取登录返回界面失败！", err)
		fmt.Println("读取登录返回界面失败！", err)
		return err
	}
	re := regexp.MustCompile(model.All.Regexp.PwdErrorRe)
	match := re.FindSubmatch(content)
	if match != nil {
		log.Println("账号或者密码错误，请修改。", "账号：", user.UserID, "密码：", user.UserPwd)
		fmt.Println("账号或者密码错误，请修改。", "账号：", user.UserID, "密码：", user.UserPwd)
		return errors.New("账号或者密码错误")
	}
	re = regexp.MustCompile(model.All.Regexp.IsNotStudentRe)
	match = re.FindSubmatch(content)
	if match != nil {
		log.Println("账号或者密码错误，请修改。", "账号：", user.UserID, "密码：", user.UserPwd)
		fmt.Println("账号或者密码错误，请修改。", "账号：", user.UserID, "密码：", user.UserPwd)
		return errors.New("账号或者密码错误")
	}
	//验证码错误
	re = regexp.MustCompile(`验证码`)
	match = re.FindSubmatch(content)
	if match != nil {
		log.Println("验证码错误!")
		fmt.Println("验证码错误!")
		return errors.New("验证码错误")
	}
	return nil
}

func Get_Login_Page(client *http.Client, user model.UserInfo) ([]byte, error) {
	maxTry := 3
	var content []byte
	var err error
	for i := 0; i < maxTry; i++ {
		statusCode := 0
		content, statusCode, err = Fetch_ISP_Login_Page(client)
		log.Println("第", i, "次", "LoginPage status code:", statusCode)
		if err != nil {
			log.Println("访问ISP登录界面失败！尝试使用VPN登录", err)
			fmt.Println("访问ISP登录界面失败,尝试使用VPN登录......")
			//使用VPN登录
			content, statusCode, err = Fetch_ISP_Login_Pag_VPN(client, user)
			if err != nil {
				log.Println("VPN登录失败！", "code:", statusCode, "err:", err)
				color.Red("VPN登录失败:", err)
				return nil, err
			}
			model.UseVPN = true
			break
		}
		if len(content) < 20 { //或者 statusCode != 200
			err = errors.New("len(content) too short")
			log.Println("访问ISP登录界面失败！", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
	if err != nil {
		fmt.Println("访问ISP登录界面失败！", err)
		//将页面内容写入到文件用于debug
		ioutil.WriteFile("loginError.html", content, 0644)
		return nil, err
	}
	return content, nil
}
