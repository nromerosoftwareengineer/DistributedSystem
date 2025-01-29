package api

import (
	"encoding/json"
	"go_proj/runtime"
	"log"
	"net/http"
)

type HttpHandler struct {
	AppContext *runtime.AppContext
	Redis      *runtime.Redis
}

func NewHttpHandler(appContext *runtime.AppContext, redis *runtime.Redis) *HttpHandler {
	return &HttpHandler{
		AppContext: appContext,
		Redis:      redis,
	}
}

func (httpHandler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var group runtime.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	groupJSON, err := json.Marshal(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpHandler.Redis.Set(group.GroupName, groupJSON)
	groupString, err := httpHandler.Redis.Get(group.GroupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("group json %s\n", groupString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}
