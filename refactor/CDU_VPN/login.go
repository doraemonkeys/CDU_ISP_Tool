package CDU_VPN

import (
	"fmt"
	"net/http"
)

func (tool *VPN_Tool) Login() (*http.Client, error) {
	// 登录VPN
	fmt.Println("VPN登录学生：", tool.CDU_Student)
	if tool.VpnPwd == "" {
		return nil, fmt.Errorf("VPN密码为空")
	}
	return http.DefaultClient, nil
}
