package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func startHttp(addr string, markdownFile string, wsCh chan string) {
	s := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 14,
	}
	http.Handle("/", &reqHandler{MarkdownFile: markdownFile, docroot: "./docroot"})
	http.HandleFunc("/ws", wsChServer(wsCh))

	log.Fatal(s.ListenAndServe())
}

type reqHandler struct {
	docroot      string
	MarkdownFile string
	docData      DocData
}

type DocData struct {
	Filename      string
	ParsedContent template.HTML
}

func (self reqHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.RequestURI {
	default:
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	case "/":
		self.html(w, r)
	case "/app.css":
		self.css(w, r)
	case "/app.js":
		self.js(w, r)
	}
}

func (self reqHandler) html(w http.ResponseWriter, r *http.Request) {
	self.docData = DocData{
		ParsedContent: template.HTML(readMd(self.MarkdownFile)),
		Filename:      self.MarkdownFile,
	}
	var t = template.Must(template.New("page").Parse("html"))
	tmpl := template.Must(t.Parse(PAGE))
	var html bytes.Buffer
	err := tmpl.Execute(&html, self.docData)
	if err != nil {
		log.Fatalf("Error template execution: %s", err)
	}
	w.Write([]byte(html.String()))
}

func (self reqHandler) css(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(CSS))
}

func (self reqHandler) js(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/js")
	w.Write([]byte(JS))
}

func reader(ws *websocket.Conn) {
	pongWait := 60 * time.Second
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, wsCh chan string) {
	pingPeriod := 50 * time.Second //must be less than pongWait
	deadlineTime := 10 * time.Second
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case content := <-wsCh:
			ws.SetWriteDeadline(time.Now().Add(deadlineTime))
			err := ws.WriteMessage(websocket.TextMessage, []byte(content))
			if err != nil {
				log.Println(err)
				return
			}

		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(deadlineTime))
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func wsChServer(wsCh chan string) func(http.ResponseWriter, *http.Request) {
	var upgrader = websocket.Upgrader{}
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			_, ok := err.(websocket.HandshakeError)
			if !ok {
				log.Println(err)
			}
			return
		}

		go writer(ws, wsCh)
		reader(ws)
	}
}
