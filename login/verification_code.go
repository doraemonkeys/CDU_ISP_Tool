package login

import (
	"ISP_Tool/model"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/dlclark/regexp2"
)

func Get_ISP_Login_code(content []byte) (string, error) {
	//fmt.Println(string(content))
	re := regexp.MustCompile(model.All.Regexp.VerificationCodeRe)
	substr := re.FindSubmatch(content)
	if substr == nil {
		log.Println("获取登录验证码失败！可能是ISP结构发生变化。")
		fmt.Println("获取登录验证码失败！可能是ISP结构发生变化。")
		return "", errors.New("substr == nil")
	}
	return string(substr[1]), nil
}

func Fetch_ISP_Login_Page(client *http.Client) ([]byte, int, error) {
	req, _ := http.NewRequest("GET", model.All.DirectBaseURL+model.All.Login.LoginWebUrl, nil)
	req.Header.Set("User-Agent", model.UserAgent)
	//设置超时时间1.5s
	client.Timeout = time.Second + time.Second/2
	resp, err := client.Do(req)
	//恢复超时时间
	client.Timeout = 0
	if err != nil {
		return nil, 0, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取登录页面失败！", err)
		fmt.Println("读取登录页面失败！", err)
		return nil, resp.StatusCode, err
	}
	return content, resp.StatusCode, nil
}

func getPostField(content []byte) (url.Values, error) {
	var parame = url.Values{}
	re := regexp2.MustCompile(model.All.Login.MatchParamesRe, 0)
	rematch, err := re.FindStringMatch(string(content))
	if err != nil {
		log.Println("正则表达式匹配出现问题！")
		fmt.Println("正则表达式匹配出现问题！")
		return nil, errors.New("正则表达式匹配出现问题！")
	}
	for rematch != nil {
		parame.Add(rematch.GroupByNumber(2).String(), rematch.GroupByNumber(4).String())
		rematch, err = re.FindNextMatch(rematch)
		if err != nil {
			log.Println("正则表达式匹配出现问题！")
			fmt.Println("正则表达式匹配出现问题！")
			return nil, errors.New("正则表达式匹配出现问题！")
		}
	}
	//打印parame的所有key
	count := 0
	for k, v := range parame {
		for range v {
			fmt.Printf("%s ", k)
			count++
		}
	}
	//打印parame的数量
	fmt.Println("共", count, "个字段")
	return parame, nil
}
