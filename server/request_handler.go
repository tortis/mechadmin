package main

import "strconv"

type WSRequest struct {
	R  string
	A1 string
	A2 string
	A3 string
}

func (r *WSRequest) Handle(c *connection) {
	switch r.R {
	case "list-comp":
		colUID64, _ := strconv.ParseUint(r.A1, 10, 32)
		response := []byte("{temp: test,requested:" + strconv.FormatUint(colUID64, 10) + "}")
		c.send <- response
	case "new-col":
	default:
		c.send <- []byte("UR")
	}
}
