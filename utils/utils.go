package utils

import (
	"ISP_Tool/model"
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/fatih/color"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 自动检测html编码,不会减少缓冲器的内容
func DetermineEncodingbyPeek(r *bufio.Reader) (encoding.Encoding, error) {
	tempbytes, err := r.Peek(1024)
	if err != nil && err != io.EOF {
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
	content, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(content)), nil
}

// 初始化client
func Get_client() (http.Client, error) {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return http.Client{Jar: jar}, nil
}

// 获取当前的执行路径(包含可执行文件名称)
// C:\Users\*\AppData\Local\Temp\*\exe\main.exe
// (读取命令行的方式，可能得不到想要的路径)
func GetCurrentPath() (string, error) {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

// 获取当前文件的详细路径
// D:/Go/workspace/port/network_learn/server/server.go
func CurrentFilePath() (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("can not get current file info")
	}
	return file, nil
}

func NetWorkStatus() bool {
	timeout := time.Duration(time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("https://www.baidu.com")
	if err != nil {
		log.Println("测试网络连接出现问题！", err)
		return false
	}
	defer resp.Body.Close()
	log.Println("Net Status , OK", resp.Status)
	return true
}

// 从文件末尾按行读取文件。
// name:文件路径 lineNum:读取行数(超过文件行数则读取全文)。
// 最后一行为空也算读取了一行,会返回此行为空串,若全是空格也会原样返回。
// 返回的每一行都不包含换行符号。
func ReverseRead(name string, lineNum uint) ([]string, error) {
	//打开文件
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//获取文件大小
	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fs.Size()

	var offset int64 = -1   //偏移量，初始化为-1，若为0则会读到EOF
	char := make([]byte, 1) //用于读取单个字节
	lineStr := ""           //存放一行的数据
	buff := make([]string, 0, 100)
	for (-offset) <= fileSize {
		//通过Seek函数从末尾移动游标然后每次读取一个字节，offset为偏移量
		file.Seek(offset, io.SeekEnd)
		_, err := file.Read(char)
		if err != nil {
			return buff, err
		}
		if char[0] == '\n' {
			//判断文件类型为unix(LF)还是windows(CRLF)
			file.Seek(-2, io.SeekCurrent) //io.SeekCurrent表示游标放置于当前位置，逆向偏移2个字节
			//读完一个字节后游标会自动正向偏移一个字节
			file.Read(char)
			if char[0] == '\r' {
				offset-- //windows跳过'\r'
			}
			lineNum-- //到此读取完一行
			buff = append(buff, lineStr)
			lineStr = ""
			if lineNum == 0 {
				return buff, nil
			}
		} else {
			lineStr = string(char) + lineStr
		}
		offset--
	}
	buff = append(buff, lineStr)
	return buff, nil
}

// 回车后返回true
func PressToContinue(ch chan bool) {
	fmt.Scanf("\n")
	ch <- true
	close(ch)
}

// 不等待执行完毕就返回,如果params中有转义字符需要自己处理,
// dir为cmd命令执行的位置,传入空值则为默认路径。
func Cmd_NoWait(dir string, params []string) (cmd *exec.Cmd, err error) {
	cmd = exec.Command("cmd")
	cmd_in := bytes.NewBuffer(nil)
	cmd.Stdin = cmd_in
	if dir != "" {
		cmd.Dir = dir
	}
	command := ""
	for i := 0; i < len(params); i++ {
		command = command + params[i]
		if i != len(params)-1 {
			command += " "
		}
	}
	cmd_in.WriteString(command + "\n")
	err = cmd.Start() //不等待执行完毕就返回
	if err != nil {
		return cmd, err
	}
	//等待cmd已经读取指令
	for cmd_in.Len() != 0 {
		time.Sleep(time.Microsecond * 10)
	}
	return cmd, nil
}

func GbkToUtf8(b []byte) []byte {
	tfr := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(tfr)
	if e != nil {
		return nil
	}
	return d
}

// attributes描述了后面每个字符串的颜色属性，attributes与strs长度必须相同,
// 注意不要忘了带上空格和换行。
func ColorPrint(attributes []color.Attribute, strs ...string) {
	for k, str := range strs {
		if attributes[k] != 0 {
			color.Set(attributes[k])
			fmt.Print(str)
			color.Unset()
		} else {
			fmt.Print(str)
		}
	}
}

// GB18030
func Utf8ToANSI(text string) string {
	return mahonia.NewEncoder("GB18030").ConvertString(text)
}

// 不会对内容转码
func Fetch(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	//在浏览器中找到request.Header(请求头)中的User-Agent,把值复制下来
	//add key value
	request.Header.Add("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	//模拟客户端发送请求
	response, err := http.DefaultClient.Do(request)
	//response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(response.StatusCode))
	}
	return io.ReadAll(response.Body)
}

// 比较两个版本号(格式为v1.0.0)v1,v2的大小,如果v1>v2返回1，v1<v2返回-1，v1=v2返回0
func CompareVersion(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")
	v1 = strings.TrimPrefix(v1, "V")
	v2 = strings.TrimPrefix(v2, "V")
	v1Arr := strings.Split(v1, ".")
	v2Arr := strings.Split(v2, ".")
	for i := 0; i < len(v1Arr); i++ {
		v1Int, _ := strconv.Atoi(v1Arr[i])
		v2Int, _ := strconv.Atoi(v2Arr[i])
		if v1Int > v2Int {
			return 1
		} else if v1Int < v2Int {
			return -1
		}
	}
	return 0
}

func GetUpdateInfo() (model.Update, error) {
	var update model.Update
	//获取最新版本信息
	resp, err := Fetch("https://gitee.com/doraemonkey/json_isp/raw/master/update.txt")
	if err != nil {
		log.Println("获取更新信息失败", err)
		color.Red("获取更新信息失败!")
		return update, err
	}
	//解析json
	err = json.Unmarshal(resp, &update)
	if err != nil {
		log.Println("解析更新信息失败", err)
		color.Red("解析更新信息失败!")
		return update, err
	}
	return update, nil
}

// filename为文件存储的路径(可省略)和文件名,记得校验md5
func DownloadFile(url string, filename string) error {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("user-agent", model.UserAgent)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("访问url失败,err:%w", err)
	}
	defer resp.Body.Close()
	// 创建一个文件用于保存
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	//io.Copy() 方法将副本从 src 复制到 dst ，直到 src 达到文件末尾 ( EOF ) 或发生错误，
	//然后返回复制的字节数和复制时遇到的第一个错误(如果有)。
	//将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// 获取文件md5(字母小写)
func GetFileMd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	//将[]byte转成16进制的字符串表示
	//var hex string = "48656c6c6f"//(hello)
	//其中每两个字符对应于其ASCII值的十六进制表示,例如:
	//0x48 0x65 0x6c 0x6c 0x6f = "Hello"
	//fmt.Printf("%x\n", hash.Sum(nil))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 等待执行完毕才返回,不反回输出
func CmdNoOutput(dir string, params []string) error {
	cmd := exec.Command("cmd")
	cmd_in := bytes.NewBuffer(nil)
	cmd.Stdin = cmd_in
	if dir != "" {
		cmd.Dir = dir
	}
	command := ""
	for i := 0; i < len(params); i++ {
		command = command + params[i]
		if i != len(params)-1 {
			command += " "
		}
	}
	cmd_in.WriteString(command + "\n")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// 等待动画(用完记得换行+'\n')
func WaitAnimation(ctx context.Context) {
	n := 100
	for i := 0; i < n; i++ {
		fmt.Printf("%s", "-")
	}
	//光标回到第一行
	fmt.Printf("\r")
	attributes := []color.Attribute{color.FgGreen}
	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			ColorPrint(attributes, ">")
			time.Sleep(time.Second / 4)
			if i == n-1 {
				fmt.Printf("\r")
				for i := 0; i < n; i++ {
					fmt.Printf("%s", "-")
				}
				fmt.Printf("\r")
				i = -1
			}
		}
	}
}
