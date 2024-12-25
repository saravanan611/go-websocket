package websocketutil

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"websocket/base"
)

func (wsserver *WebSocketServer) GetData(w http.ResponseWriter, r *http.Request) {

	if !strings.EqualFold(http.MethodPost, r.Method) {
		fmt.Fprint(w, base.CreateMsg(base.GErrorFlag, "Method not allowed"))
		return
	}

	lBody, lErr := io.ReadAll(r.Body)
	if lErr != nil {
		log.Println(lErr.Error())
		fmt.Fprint(w, base.CreateMsg(base.GErrorFlag, lErr.Error()))
		return
	}
	lErr = wsserver.Publish(r, lBody)
	if lErr != nil {
		log.Println(lErr.Error())
		fmt.Fprint(w, base.CreateMsg(base.GErrorFlag, lErr.Error()))
		return
	}
	fmt.Fprint(w, base.CreateMsg(base.GSuccessFlag, "Data send Successfully"))

}
