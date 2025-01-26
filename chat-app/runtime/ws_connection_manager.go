package runtime

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type userWSConnMap map[string]*websocket.Conn
type ConnectionHandler struct {
	wsmap           userWSConnMap
	connectionMutex sync.Mutex
}

func NewConnectionHandler() *ConnectionHandler {
	mh := &ConnectionHandler{
		wsmap:           make(userWSConnMap),
		connectionMutex: sync.Mutex{},
	}
	return mh
}

func (ch *ConnectionHandler) AddUserConn(userId string, conn *websocket.Conn) {
	ch.connectionMutex.Lock()
	defer ch.connectionMutex.Unlock()
	ch.wsmap[userId] = conn
}

func (ch *ConnectionHandler) CloseUserConn(userId string) {
	ch.connectionMutex.Lock()
	defer ch.connectionMutex.Unlock()
	err := ch.wsmap[userId].Close()
	delete(ch.wsmap, userId)
	if err != nil {
		return
	}
}

func (ch *ConnectionHandler) GetUserConn(userId string) *websocket.Conn {
	return ch.wsmap[userId]
}

func (ch *ConnectionHandler) Close() {
	ch.connectionMutex.Lock()
	defer ch.connectionMutex.Unlock()
	for userId, conn := range ch.wsmap {
		err := conn.Close()
		if err != nil {
			log.Println("close connection error: for userId: ", err, userId)
		}
	}
	ch.wsmap = make(userWSConnMap)
}
