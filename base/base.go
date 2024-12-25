package base

import "fmt"

const (
	GBaseAPI     = "/socket"
	GErrorFlag   = "E"
	GSuccessFlag = "S"
	GConnName    = "conn"
	GActionName  = "action"
)

var GPortNo = "80"

func CreateMsg(pStatus, pMsg string) (lResp string) {
	return fmt.Sprintf(`{"status":"%s","msg":"%s"}`, pStatus, pMsg)
}
