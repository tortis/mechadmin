package main

import (
	"code.google.com/p/go.net/websocket"
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Body  template.HTML
}

var templates = template.Must(template.ParseFiles("www/main.html"))

func handler(w http.ResponseWriter, r *http.Request, c *CollectionTree) {
	p := &Page{Title: "Computers", Body: template.HTML(c.GenerateHTML())}
	err := templates.ExecuteTemplate(w, "main.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn, ct *CollectionTree) {
	c := &connection{send: make(chan []byte, 256), ws: ws, colTree: ct}
	wsHub.register <- c
	defer func() { wsHub.unregister <- c }()
	go c.writer()
	c.reader()
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *CollectionTree), c *CollectionTree) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, c)
	}
}

func makeWSHandler(fn func(*websocket.Conn, *CollectionTree), c *CollectionTree) websocket.Handler {
	return func(ws *websocket.Conn) {
		fn(ws, c)
	}
}

func StartWebServer(c *CollectionTree) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("www"))))
	http.HandleFunc("/lab/", makeHandler(handler, c))
	http.Handle("/ws/", makeWSHandler(EchoServer, c))
	go wsHub.run()
	http.ListenAndServe(":8080", nil)
}
