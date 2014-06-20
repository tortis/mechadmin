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

func handler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Computers", Body: template.HTML(root.GenerateHTML())}
	err := templates.ExecuteTemplate(w, "main.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Echo the data received on the WebSocket.
func WSHandler(ws *websocket.Conn) {
	c := &connection{send: make(chan []byte, 256), ws: ws}
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

func StartWebServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("www/static"))))
	http.HandleFunc("/", handler)
	http.Handle("/ws/", websocket.Handler(WSHandler))
	go wsHub.run()
	http.ListenAndServe(":8080", nil)
}
