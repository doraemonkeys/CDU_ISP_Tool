package login

import (
	"ISP_Tool/model"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func Get_ISP_Login_code(client *http.Client) (string, error) {
	req, _ := http.NewRequest("GET", "https://xsswzx.cdu.edu.cn/ispstu/com_user/weblogin.asp", nil)
	req.Header.Set("User-Agent", model.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("获取登录验证码失败！可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("获取登录验证码失败！可能是ISP结构发生变化，请联系开发者。")
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取登录页面失败！", err)
		fmt.Println("读取登录页面失败！", err)
		return "", err
	}
	//fmt.Println(string(content))
	re := regexp.MustCompile(model.VerificationCodeRe)
	substr := re.FindSubmatch(content)
	if substr == nil {
		log.Println("获取登录验证码失败！可能是ISP结构发生变化，请联系开发者。")
		fmt.Println("获取登录验证码失败！可能是ISP结构发生变化，请联系开发者。")
		return "", errors.New("substr == nil")
	}
	return string(substr[1]), nil
}
