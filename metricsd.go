package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/KawaiiDevs/kawaii-metrics/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Message contains some information about the command used
type Message struct {
	CommandName string `json:"command_name"`
	IsNSFW bool `json:"is_nsfw"`
}

var mainLog = log.New()

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
} 

// HandleWS handles /ws endpoint
func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		mainLog.WithError(err).Errorln("encountered error while trying to upgrade /ws")
		return
	}

	defer conn.Close()
	// begin main loop
	for {
		message := &Message{}
		err = conn.ReadJSON(message)
		if err != nil {
			mainLog.WithError(err).Errorln("encountered error while trying to read json")
			return
		}
		labels := prom.Labels{"name": message.CommandName}
		if message.IsNSFW {
			prometheus.NSFWCommandsTotal.With(labels).Inc()
		} else {
			prometheus.NormalCommandsTotal.With(labels).Inc()
		}
	}
}

// HandleMetrics handles /metrics endpoint
func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	contentType, err := prometheus.WriteMetrics(w, r.Header.Get("Accept"))
	if err != nil {
		mainLog.WithError(err).Errorln("failed to generate prom metrics")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "text/plain; charset=utf8")
		fmt.Fprint(w, "500 Internal Server Error")
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", contentType)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", HandleWS)
	r.HandleFunc("/prometheus/metrics", HandleMetrics)
	err := http.ListenAndServe("0.0.0.0:8080", r) // TODO
	if err != nil {
		mainLog.WithError(err).Errorln("couldn't listen")
	}
}
