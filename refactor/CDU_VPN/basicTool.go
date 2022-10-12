package CDU_VPN

import (
	"CDU_Tool/login"
	"fmt"
	"net/http"
)

type VPN_Login_Person struct {
	UserID  string
	VPN_Pwd string
}

type Basic_VPN_Tool struct {
	person *VPN_Login_Person
	Client *http.Client
}

func (tool *Basic_VPN_Tool) Login(p login.LoginToolPerson) error {
	_, ok := p.(*VPN_Login_Person)
	if !ok {
		return fmt.Errorf("参数错误")
	}
	tool.person = p.(*VPN_Login_Person)
	fmt.Println("VPN登录人：", tool.person)
	// 登录VPN
	return nil
}

func (p *VPN_Login_Person) UseLoginTool(tool login.LoginTool) error {
	err := tool.Login(p)
	if err != nil {
		return err
	}
	return nil
}
