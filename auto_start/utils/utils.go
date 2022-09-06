package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
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

//获取当前的执行路径(包含可执行文件名称)
//C:\Users\*\AppData\Local\Temp\*\exe\main.exe
//(读取命令行的方式，可能得不到想要的路径)
func GetCurrentPath() (string, error) {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
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

//从文件末尾按行读取文件。
//name:文件路径 lineNum:读取行数(超过文件行数则读取全文)。
//最后一行为空也算读取了一行,会返回此行为空串,若全是空格也会原样返回。
//返回的每一行都不包含换行符号。
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

//读取倒数第n行(n从1开始),
//若n大于文件行数则返回错误io.EOF。
func ReadStartWithLastLine(filename string, n int) (string, error) {
	//打开文件
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	//获取文件大小
	fs, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fs.Size()

	var offset int64 = -1   //偏移量，初始化为-1，若为0则会读到EOF
	char := make([]byte, 1) //用于读取单个字节
	lineStr := ""           //存放一行的数据
	lineCount := 0          //行数
	for (-offset) <= fileSize {
		//通过Seek函数从末尾移动游标然后每次读取一个字节，offset为偏移量
		file.Seek(offset, io.SeekEnd)
		_, err := file.Read(char)
		if err != nil {
			return "", err
		}
		if char[0] == '\n' {
			lineCount++
			if lineCount == n {
				return lineStr, nil
			}
			//防止偏移量-2后越界
			if fileSize-(-offset) >= 1 {
				//判断文件类型为unix(LF)还是windows(CRLF)
				file.Seek(-2, io.SeekCurrent) //io.SeekCurrent表示游标放置于当前位置，逆向偏移2个字节
				//读完一个字节后游标会自动正向偏移一个字节
				file.Read(char)
				if char[0] == '\r' {
					offset-- //windows跳过'\r'
				}
			}
			offset--
			continue
		}
		if lineCount == n-1 {
			lineStr = string(char) + lineStr
		}
		offset--
	}
	//到此文件已经从尾部读到头部
	if lineCount == n-1 {
		return lineStr, nil
	}
	return "", io.EOF
}

func PressToContinue(ch chan bool) {
	fmt.Scanf("\n")
	ch <- true
	close(ch)
}

//不等待执行完毕就返回
//如果params中有转义字符需要自己处理
//dir为cmd命令执行的位置,传入空值则为默认路径
func Cmd_NoWait(dir string, params []string) (string, error) {
	cmd := exec.Command("cmd")
	cmd_in := bytes.NewBuffer(nil)
	cmd.Stdin = cmd_in
	cmd_out := bytes.NewBuffer(nil)
	cmd.Stdout = cmd_out
	cmd.Dir = dir
	command := ""
	for i := 0; i < len(params); i++ {
		command = command + params[i]
		if i != len(params)-1 {
			command += " "
		}
	}
	cmd_in.WriteString(command + "\n")
	err := cmd.Start() //不等待执行完毕就返回
	if err != nil {
		return "", err
	}
	output := cmd_out.Bytes()
	return string(GbkToUtf8(output)), nil
}

func GbkToUtf8(b []byte) []byte {
	tfr := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(tfr)
	if e != nil {
		return nil
	}
	return d
}

//attributes描述了后面每个字符串的颜色属性，attributes与strs长度必须相同,
//注意不要忘了带上空格和换行。
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

//GB18030
func Utf8ToANSI(text string) string {
	return mahonia.NewEncoder("GB18030").ConvertString(text)
}

//获取昨天的日期
func GetYesterday() time.Time {
	return time.Now().AddDate(0, 0, -1)
}
