package CDU_ISP

import (
	"CDU_Tool/CDU_VPN"
	"net/http"
)

func TryLogin_VPN(client *http.Client, person CDU_VPN.VPN_Login_Person) error {
	vpn_tool := CDU_VPN.Stu_VPN_Tool{
		BasicTool: &CDU_VPN.Basic_VPN_Tool{Client: client},
	}
	err := person.UseLoginTool(vpn_tool.BasicTool)
	if err != nil {
		return err
	}
	return nil
}
