package runtime

import "log"

type AppContext struct {
	CH *ConnectionHandler
	MH *MessageHandler
}

func NewAppContext(redis *Redis) *AppContext {
	connectionHandler := NewConnectionHandler()
	return &AppContext{
		CH: connectionHandler,
		MH: NewMessageHandler(connectionHandler, redis),
	}
}

func (appContext *AppContext) Clean() {
	log.Printf("Cleaning up connection handler/message handler")
	appContext.CH.Close()
	appContext.MH.Close()
}
