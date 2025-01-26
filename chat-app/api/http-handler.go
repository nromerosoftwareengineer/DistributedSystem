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
	log.Printf("Received request to /create-group endpoint. Method: %s", r.Method)
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
