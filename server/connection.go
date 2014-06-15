package main

import "code.google.com/p/go.net/websocket"
import "encoding/json"
import "strconv"

type connection struct {
	ws      *websocket.Conn
	send    chan []byte
	colTree *CollectionTree
}

func (c *connection) reader() {
	var buffer []byte = make([]byte, 1024)
	for {
		br, err := c.ws.Read(buffer)
		if err != nil {
			break
		}
		if string(buffer[0:br]) == "bye" {
			break
		}
		println(string(buffer[0:br]))
		var req WSRequest
		err = json.Unmarshal(buffer[0:br], &req)
		if err != nil {
			println(err)
			break
		}
		println(req.R)
		switch req.R {
		case "list-comp":
			colUID64, _ := strconv.ParseUint(req.A1, 10, 32)
			response := []byte("{temp: test,requested:"+strconv.FormatUint(colUID64, 10)+"}")
			c.send <- response
		case "eat-shit":
		}
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		_, err := c.ws.Write(message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}
