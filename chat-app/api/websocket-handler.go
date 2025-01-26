package api

import (
	"github.com/gorilla/websocket"
	"go_proj/runtime"
	"log"
	"net/http"
)

type WebSocketHandler struct {
	Upgrader   websocket.Upgrader
	AppContext *runtime.AppContext
}

func NewWebSocketHandler(appContext *runtime.AppContext) *WebSocketHandler {
	return &WebSocketHandler{
		Upgrader:   websocket.Upgrader{},
		AppContext: appContext,
	}
}

func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	log.Printf("userId: %s is trying to connect, upgrading to websocket connection\n", userId)
	c := wsh.upgradeToWs(w, r)
	if c == nil {
		return
	}
	wsh.AppContext.CH.AddUserConn(userId, c)
	wsh.AppContext.MH.HandleMessageLoop(userId)
}

func (wsh WebSocketHandler) upgradeToWs(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	c, err := wsh.Upgrader.Upgrade(w, r, nil)
	if err == nil {
		return c
	}
	log.Printf("error %s when upgrading connection to websocket", err)
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write([]byte("Unable to upgrade to websocket"))
	if err != nil {
		log.Printf("error %s when trying to return error message to client", err)
	}
	return nil
}
