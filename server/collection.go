package main

import "github.com/petar/GoLLRB/llrb"

type Collection struct {
	Name      string
	UID       uint32
	Computers *llrb.LLRB
}

func NewCollection(name string, uid uint32) *Collection {
	c := &Collection{
		Name:      name,
		UID:       uid,
		Computers: llrb.New(),
	}
	return c
}

func (col *Collection) AddComputer(name string) {
	col.Computers.ReplaceOrInsert(llrb.String(name))
	wsHub.broadcast <- []byte("Adding new computer, " + name + " into collection, " + col.Name)
}

func (col *Collection) RemoveComputer(name string) {
	if (col.Computers.Delete(llrb.String(name)) != nil) {
		wsHub.broadcast <- []byte("Removing computer, " + name + " from collection, " + col.Name)
	}
}
