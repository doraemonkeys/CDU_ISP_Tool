package runningLog

import (
	"log"

	"github.com/go-toast/toast"
)

func Inform(content string) {
	notification := toast.Notification{
		AppID:   "CUD_ISP_TOOL",
		Title:   "Information",
		Message: content,
	}
	err := notification.Push()
	if err != nil {
		log.Println("系统通知发送失败！", err)
	}
}
