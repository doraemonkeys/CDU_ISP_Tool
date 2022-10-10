package ISP

import "net/http"

type ISP_Tool interface {
	Login() error
	CheckIn() error
}

type CDUStudetForISP struct {
	CheckInStudent
}

func (cs *CDUStudetForISP) Login() error {
	cs.Client = &http.Client{}
	return nil
}

func (cs *CDUStudetForISP) CheckIn() error {
	return nil
}
