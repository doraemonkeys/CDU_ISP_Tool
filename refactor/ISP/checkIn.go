package ISP

import (
	"fmt"
)

// 单一职责原则，一个类只负责一项职责
func (cs *CheckInStudent) CheckIn() error {
	var isp_toll ISP_Tool = &CDUStudetForISP{CheckInStudent: *cs}
	err := isp_toll.Login()
	if err != nil {
		return err
	}
	fmt.Println("登录成功")
	err = isp_toll.CheckIn()
	if err != nil {
		return err
	}
	fmt.Println("打卡成功")
	return nil
}
