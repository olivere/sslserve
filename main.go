package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	hostname string

	upgrader = websocket.Upgrader{}
)

func main() {
	hostname, _ = os.Hostname()

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/websocket", ws)
	httpSrv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}
	log.Fatal(httpSrv.ListenAndServe())
}

func home(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	fmt.Fprintf(w, "Hostname: %s\nCurrent time is %v in %v\n", hostname, now, now.Location())
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read: %v", err)
			break
		}
		log.Printf("recv: %s", msg)
		err = conn.WriteMessage(mt, msg)
		if err != nil {
			log.Printf("ws write: %v", err)
			break
		}
	}
}
