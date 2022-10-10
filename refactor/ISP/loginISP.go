package ISP

import "net/http"

func (cs *CDUStudetForISP) Login() error {
	cs.Client = &http.Client{}
	return nil
}

func (cs *CDUStudetForISP) CheckIn() error {
	return nil
}
