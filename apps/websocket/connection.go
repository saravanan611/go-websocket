package websocketutil

import (
	"fmt"
	"net/http"
	"websocket/base"
)

func (wsserver *WebSocketServer) GetConnection(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		fmt.Fprint(w, base.CreateMsg(base.GErrorFlag, "Method not allowed"))
		return
	}

	lErr := wsserver.Subscribe(w, r)
	if lErr != nil {
		fmt.Fprint(w, base.CreateMsg(base.GErrorFlag, lErr.Error()))
		return
	}
}
