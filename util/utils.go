package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/encoding"
)

//自动检测html编码,不会减少缓冲器的内容
func DetermineEncodingbyPeek(r *bufio.Reader) (encoding.Encoding, error) {
	tempbytes, err := r.Peek(1024)
	if err != nil {
		return nil, err
	}
	e, _, _ := charset.DetermineEncoding(tempbytes, "")
	return e, nil
}

func GetIPV4() (string, error) {
	resp, err := http.Get("https://ipv4.netarm.com")
	if err != nil {
		log.Println("获取公网IP失败", err)
		fmt.Println("获取公网IP失败", err)
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return strings.TrimSpace(string(content)), nil
}

// 初始化client
func Get_client() (http.Client, error) {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return http.Client{Jar: jar}, nil
}
