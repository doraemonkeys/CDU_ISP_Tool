package view

import (
	"ISP_Tool/model"
	"fmt"
)

func Menu() {
	fmt.Println("************************************************************************")
	fmt.Println("                    成都大学ISP健康打卡小工具                 	           ")
	fmt.Println("                                                                        ")
	fmt.Println("    更新地址:https://github.com/Doraemonkeys/CDU_ISP_Tool/releases       ")
	fmt.Println("************************************************************************")
}

func Warn() {
	fmt.Println("                  WARNNING                             ")
	fmt.Println("使用时请确认你正处于低风险地区且未与阳性患者发生密切接触！  ")
	fmt.Println("隐瞒风险所造成的一切后果由使用者自己承担！                 ")
}

func Clock_IN_Success(user model.UserInfo) {
	fmt.Println("**************************************************************")
	fmt.Printf(">>>>>>>>>>>>>>>ID  %s", user.UserID)
	fmt.Printf("    健康登记打卡  成功！                           \n")
	fmt.Println("**************************************************************")
}

func Clock_IN_Failed(user model.UserInfo) {
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	fmt.Printf(">>>>>>>>>>>>>>>ID  %s", user.UserID)
	fmt.Printf("    健康登记打卡  失败！                           \n")
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
}

func Auto_Clock_IN_Success() {
	fmt.Println("**************************************************************")
	fmt.Printf(">>>>>>>>>>>>>>> 已经自动打卡成功！                           \n")
	fmt.Println("**************************************************************")
}

func EndSlect() {
	fmt.Println("[0]  删除一个账号")
	fmt.Println("[1]  修改账号密码")
	fmt.Println("[2]  批量添加账号")
	fmt.Println("[3]  为添加的用户重新打卡")
	fmt.Println("[4]  清空账号")
	fmt.Println("[5]  开启/关闭每日定时自动打卡")
	fmt.Println()
	fmt.Println(">>>>>>>>>>>无操作120秒后自动退出<<<<<<<<<<<")
	fmt.Println("请选择 【0 - 5】:")
}
