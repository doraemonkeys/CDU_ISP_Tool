package teacher

import (
	"CDU_Tool/CDU_ISP"
	"CDU_Tool/CDU_VPN"
	"CDU_Tool/checkIn"
	"fmt"
	"net/http"
)

// 合成复用原则，对于继承和组合，优先使用组合(匿名结构体为继承)
type Teacher_ISP_Tool struct {
	teacher *CDU_CheckInTeacher
	Client  *http.Client
	UseVPN  bool //是否使用了VPN
}

// 实现层
func (tool *Teacher_ISP_Tool) CheckIn(p checkIn.CheckInToolPerson) error {
	tool.UseVPN = false
	if _, ok := p.(*CDU_CheckInTeacher); !ok {
		return fmt.Errorf("参数错误")
	}
	tool.teacher = p.(*CDU_CheckInTeacher)
	err := tool.login()
	if err != nil {
		err := CDU_ISP.TryLogin_VPN(tool.Client,
			CDU_VPN.VPN_Login_Person{
				UserID:  tool.teacher.TeacherId,
				VPN_Pwd: tool.teacher.VpnPwd,
			})
		if err != nil {
			return err
		}
		tool.UseVPN = true
	}
	fmt.Println("打卡老师：", p)
	fmt.Println("是否使用VPN：", tool.UseVPN)
	return nil
}

func (tool *Teacher_ISP_Tool) login() error {
	fmt.Println("登录老师：", tool.teacher)
	return nil
}
