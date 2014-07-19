package main

import "strconv"
import "encoding/json"
import "github.com/tortis/mechadmin/types"

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
			// A1 contains the ID of the collection
		ColStore.Get(uint32(colUID64)).Subscribe(c)
			// Subscribe the connection to the collection.
		MACs := ColStore.Get(uint32(colUID64)).Computers
			// Get the MACs of the computers in this collection
		comps := make([]types.Status, len(MACs))
			// Create a slice of computer Status
		for i,mac := range(MACs) {
			comps[i] = CompStore.Get(mac).Info
		}
			// Retrieve the computers from the CompStore using the mac
		rJSON, err := json.Marshal(WSResponse{"list-compR", comps})
		checkError(err)
		c.send <- rJSON
	case "new-col":
	default:
		c.send <- []byte("UR")
	}
}
