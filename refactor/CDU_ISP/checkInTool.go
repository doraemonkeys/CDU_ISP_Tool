package CDU_ISP

import (
	"CDU_Tool/CDU_VPN"
	"fmt"
	"net/http"
)

// 合成复用原则，对于继承和组合，优先使用组合(匿名结构体为继承)
type ISP_Tool struct {
	Stu    *CDU_CheckInStudent
	Client *http.Client
}

// 迪米特法则，一个对象应该对其他对象有最少的了解(黑盒原则)
func (tool *ISP_Tool) login() error {
	fmt.Println("登录学生：", tool.Stu)
	return nil
}

// 实现层
func (tool *ISP_Tool) CheckIn() error {
	UseVPN := false
	err := tool.login()
	if err != nil {
		vpn_tool := CDU_VPN.VPN_Tool{
			Stu:    &CDU_VPN.VPN_Student{CDU_Student: tool.Stu.CDU_Student},
			Client: tool.Client,
		}
		client, err := vpn_tool.Login()
		if err != nil {
			return err
		}
		tool.Client = client
		UseVPN = true
	}
	fmt.Println("打卡学生：", tool.Stu)
	fmt.Println("是否使用VPN：", UseVPN)
	return nil
}
