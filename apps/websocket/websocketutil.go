package websocketutil

import (
	_err "errors"
	_nhtp "net/http"
	_str "strings"
	_syn "sync"
	_tim "time"
	"websocket/base"

	_ws "github.com/gorilla/websocket"
)

/*
***************************************************************************************

WebSocketServer structur is used to carry the user connection data ,update data and connection prodecter
==================
Purpose :
upgrader variable =>it is used to upgrade an HTTP connection to a WebSocket connection. This is a common task when building WebSocket-based applications
ClientsMutex variable =>it is used for implementing mutual exclusion, which allows you to protect shared resources from concurrent access by multiple goroutines
connections variable => It will maintain links connecting user's in matrix format for easy data sending and won't bother unwanted connections with check-ins.
==================
Author : Saravanan
Date : 06-sep-2023

****************************************************************************************
*/
type WebSocketServer struct {
	upgrader     _ws.Upgrader
	connections  map[string]map[string]*_ws.Conn
	ClientsMutex _syn.Mutex
}

/*
***************************************************************************************

Purpose : This method is used to set the origan of the websocket at a same time inserlaze the default value to WebSocketServer2 connection structur
Arguments : origan name
===================
output : WebSocketServer2 instent have a default value for connection
===================
Author : Saravanan
Date : 06-sep-2023

***************************************************************************************
*/
func NewWebSocketServer(pOrigan string) *WebSocketServer {
	lConn := &WebSocketServer{
		// create a defauld upgrader data
		upgrader: _ws.Upgrader{
			ReadBufferSize:  0,
			WriteBufferSize: 0,
			// check the allow origan of the request
			CheckOrigin: func(r *_nhtp.Request) bool {
				if _str.Contains(_str.ToLower(r.Host), _str.ToLower(pOrigan)) || _str.EqualFold(pOrigan, "*") {
					return true
				}
				return false
			},
		},
		// create a matrix map to carry the connection id based on Subscribe
		connections: make(map[string]map[string]*_ws.Conn),
	}

	go lConn.ConnectCheck()
	return lConn
}

/*
\
***************************************************************************************

Purpose : This method is used to create a connection to user based on Subscribe
Arguments : request and response of the api call
Queryparameter : (connection and action) for Subscribe
===================
output : create a user connection based on Subscribe in matrix formet

	 map:[
		  connection@@Action:[remote Address:connection id,remote Address:connection id,remote Address:connection id,..],
	 	  connection@@Action:[remote Address:connection id,remote Address:connection id,remote Address:connection id,..]
		  ,..
	     ]

===================
error : will occur based on the unsatisfactory condition
===================
Author : Saravanan
Date : 06-sep-2023

***************************************************************************************
*/
func (wss *WebSocketServer) Subscribe(pResp _nhtp.ResponseWriter, pReq *_nhtp.Request) (lErr error) {
	// upgrade an HTTP connection to a WebSocket connection
	conn, lErr := wss.upgrader.Upgrade(pResp, pReq, nil)
	if lErr != nil {
		return lErr
	}
	// read a url to get a connection and action type
	msgType := _str.ToUpper(pReq.URL.Query().Get(base.GConnName))
	Action := _str.ToUpper(pReq.URL.Query().Get(base.GActionName))
	// check the connection and action type are not empty
	if Action == "" || msgType == "" {
		return _err.New("there are no values in query params")
	}
	// get the remote address of the request
	lID := pReq.RemoteAddr
	// create a Subscribe parent map to carry the connection id
	if wss.connections[msgType+"@@"+Action] == nil {
		wss.connections[msgType+"@@"+Action] = make(map[string]*_ws.Conn)
	}
	// create a connection in parent map in form of [remote Address:connection id,...]
	wss.ClientsMutex.Lock()
	wss.connections[msgType+"@@"+Action][lID] = conn
	wss.ClientsMutex.Unlock()
	return nil
}

/*
***************************************************************************************

Purpose : This method is used to Publish the data or message to the user based on Subscribe
Arguments : request of the api call and data thet want to send user in formet of([] byte)
Heaser values : the request contains the value of Subscribe to send a data or message for which user's
===================
output :
===================
error : will occur based on the unsatisfactory condition
===================
Author : Saravanan
Date : 06-sep-2023

***************************************************************************************
*/

func (wss *WebSocketServer) Publish(pReq *_nhtp.Request, pPublishData []byte) error {
	// read the header to find which Subscriber to publish the data
	lConnectionType := _str.ToUpper(pReq.Header.Get(base.GConnName))
	lActionType := _str.ToUpper(pReq.Header.Get(base.GActionName))
	// check the Subscribe data will empty or not
	if lConnectionType == "" || lActionType == "" {
		return _err.New("missing required values in header")
	}
	// get the child map in parent map based on Subscribe
	// check the Subscribe acction will be All or Some other to send a data to connection id's
	lChildMap, flag := wss.connections[lConnectionType+"@@"+lActionType]
	if !flag {
		return _err.New("no connection type available")
	}
	wss.Sender(lConnectionType+"@@"+lActionType, lChildMap, pPublishData)
	return nil
}

/*
***************************************************************************************

Purpose : This method is a sub-method of Publish2 that sends a message or data to subscribers while also removing connections where the connection id's is in the state of being inactive.
Arguments : Subscribe chanel name ,child map that carry the connection id's,data in ([]byte)formet
===================
output : N/A
===================
error :cleare connection data in temp storage
===================
Author : Saravanan
Date : 06-sep-2023

***************************************************************************************
*/
func (wss *WebSocketServer) Sender(pParentKey string, pChildMap map[string]*_ws.Conn, pPublishData []byte) {
	var wg _syn.WaitGroup
	// sending a data to N user's at a same time and remove the in-active connection id's at a same time in child map
	for lIP, lConn := range pChildMap {
		wg.Add(1)
		go func(lIP string, lConn *_ws.Conn) {
			defer wg.Done()
			err := lConn.WriteMessage(_ws.TextMessage, pPublishData)
			wss.ClientsMutex.Lock()
			if err != nil {
				delete(wss.connections[pParentKey], lIP)
			}
			wss.ClientsMutex.Unlock()
		}(lIP, lConn)
	}
	wg.Wait()
	wss.ClientsMutex.Lock()
	// remove the in-active Subscribe in parent map
	if len(wss.connections[pParentKey]) == 0 {
		delete(wss.connections, pParentKey)
	}
	wss.ClientsMutex.Unlock()
}

/*
***************************************************************************************

Purpose : This method is used to removing connections every given time one's where the connection id's is in the state of being inactive.
Arguments :N/A
===================
output :N/A
===================
error :N/A
===================
Author : Saravanan
Date : 06-sep-2023

***************************************************************************************
*/

func (wss *WebSocketServer) ConnectCheck() {
	for {
		var wg _syn.WaitGroup
		// check the connection id's every given time one's where it active or not
		for lParentName, lSubMap := range wss.connections {
			// sending a PingMessage to N user's at a same time and remove the in-active connection id's at a same time in child map
			for lChildName, conn := range lSubMap {
				wg.Add(1)

				go func(pChildName, pParentName string, pConn *_ws.Conn) {
					defer wg.Done()
					err := pConn.WriteControl(_ws.PingMessage, []byte{}, _tim.Now().Add(1*_tim.Second))
					if err != nil {
						pConn.Close()
						delete(wss.connections[pParentName], pChildName)
					}
				}(lChildName, lParentName, conn)
			}
			wg.Wait()
			// remove the in-active Subscribe in parent map
			if len(wss.connections[lParentName]) == 0 {
				delete(wss.connections, lParentName)
			}
		}
		// Adjust the ping interval as needed
		_tim.Sleep(10 * _tim.Second)
	}
}
