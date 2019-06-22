package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	hostname string
	upgrader = websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	tmpl = template.Must(template.New("").Parse(homeTemplate))
)

func main() {
	var err error

	hostname, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/ws", ws)
	httpSrv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}
	log.Fatal(httpSrv.ListenAndServe())
}

func home(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	err := tmpl.Execute(w, struct {
		Hostname string
		Time     string
		Location string
	}{
		Hostname: hostname,
		Time:     now.Format(time.RFC3339),
		Location: now.Location().String(),
	})
	if err != nil {
		fmt.Fprintln(w, err)
	}
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

var (
	homeTemplate = `<!doctype html>
<html>
<head>
	<meta charset="utf-8">
	<title>SSLServe</title>
</head>
<body>
<h1>SSLServe</h1>
<div>
	<p>Hostname: {{.Hostname}}</p>
	<p>
		Time: {{.Time}}<br>
		Location: {{.Location}}
	</p>
</div>

<h1>WebSocket</h1>
<form>
	<p>
		<input type="text" id="input" value="Hello" />
		<button id="send">Send</button>
	</p>
</form>
<div id="output"></div>
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    var open = function() {
        if (ws) {
            return false;
        }
        ws = new WebSocket("wss://localhost:443/ws");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
	};
	
	open();
});
</script>
</body>
</html>`
)
