package student

import (
	"CDU_Tool/CDU_ISP"
	"CDU_Tool/CDU_VPN"
	"CDU_Tool/checkIn"
	"fmt"
	"net/http"
)

// 合成复用原则，对于继承和组合，优先使用组合(匿名结构体为继承)
type Stu_ISP_Tool struct {
	stu    *CDU_CheckInStudent
	Client *http.Client
	UseVPN bool //是否使用了VPN
}

// 实现层
func (tool *Stu_ISP_Tool) CheckIn(p checkIn.CheckInToolPerson) error {
	tool.UseVPN = false
	if _, ok := p.(*CDU_CheckInStudent); !ok {
		return fmt.Errorf("参数错误")
	}
	tool.stu = p.(*CDU_CheckInStudent)
	err := tool.login()
	err = fmt.Errorf("登录失败")
	if err != nil {
		err := CDU_ISP.TryLogin_VPN(tool.Client,
			CDU_VPN.VPN_Login_Person{
				UserID:  tool.stu.StudentId,
				VPN_Pwd: tool.stu.VpnPwd,
			})
		if err != nil {
			return err
		}
		tool.UseVPN = true
	}
	fmt.Println("打卡学生：", p)
	fmt.Println("是否使用VPN：", tool.UseVPN)
	return nil
}

func (tool *Stu_ISP_Tool) login() error {
	fmt.Println("登录学生：", tool.stu)
	return nil
}
