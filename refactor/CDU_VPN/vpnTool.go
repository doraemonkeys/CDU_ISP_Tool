package CDU_VPN

import (
	"CDU_Tool/model"
	"net/http"
)

type VPN_Tool struct {
	Stu    *VPN_Student
	Client *http.Client
}

// 单一职责原则，一个类只负责一项职责
type VPN_Student struct {
	model.CDU_Student
}
