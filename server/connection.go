package main

import "code.google.com/p/go.net/websocket"
import "encoding/json"

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

func (c *connection) reader() {
	var buffer []byte = make([]byte, 1024)
	for {
		/* Read and block. */
		br, err := c.ws.Read(buffer)
		if err != nil {
			break
		}

		/* If the buffer contains the string 'bye', close connection. */
		if string(buffer[0:br]) == "bye" {
			break
		}

		/* Attempt to unmarshal the buffer into a WSRequest. */
		println(string(buffer[0:br]))
		var req WSRequest
		err = json.Unmarshal(buffer[0:br], &req)
		if err != nil {
			println(err)
			break
		}

		req.Handle(c)
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
