package main

import "strconv"
import "encoding/json"

type WSRequest struct {
	R  string
	A1 string
	A2 string
	A3 string
}

type WSResponse struct {
	R string
	D interface{}
}

func (r *WSRequest) Handle(c *connection) {
	switch r.R {
	case "list-comp":
		colUID64, _ := strconv.ParseUint(r.A1, 10, 32)
		rJSON, err := json.Marshal(WSResponse{"list-compR", ColStore.Get(uint32(colUID64)).Computers})
		checkError(err)
		c.send <- rJSON
	case "new-col":
	default:
		c.send <- []byte("UR")
	}
}
