package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	// var doc string
	// log.Printf("%s %s User-Agent: %+v", r.Method, r.RequestURI, r.Header["User-Agent"])
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
	tmpl := template.Must(template.ParseFiles(self.docroot + "/html/page.tmpl"))
	var html bytes.Buffer
	err := tmpl.Execute(&html, self.docData)
	if err != nil {
		log.Fatalf("Error template execution: %s", err)
	}
	w.Write([]byte(html.String()))
}

func (self reqHandler) css(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	cssDir := self.docroot + "/css/"
	w.Write(concatFiles(cssDir, ".css"))
}

func concatFiles(dir string, ext string) []byte {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var concatedFiles []byte
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ext) {
			// log.Printf("Reading %v", file.Name())
			f, err := ioutil.ReadFile(dir + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			concatedFiles = append(concatedFiles[:], f[:]...)
		}
	}
	return concatedFiles
}

func (self reqHandler) js(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/js")
	jsDir := self.docroot + "/js/"
	w.Write(concatFiles(jsDir, ".js"))
}

func reader(ws *websocket.Conn) {
	pongWait := 60 * time.Second
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		// log.Println("ws: Received pong")
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
			// log.Println("ws: File changed")
			ws.SetWriteDeadline(time.Now().Add(deadlineTime))
			err := ws.WriteMessage(websocket.TextMessage, []byte(content))
			if err != nil {
				log.Println(err)
				return
			}

		case <-pingTicker.C:
			// log.Printf("sending ping")
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
