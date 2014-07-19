package main

import "sort"
import "encoding/json"

/* System Collection UIDs. */
const ALL_SYS_COL uint32 = 0
const UNK_SYS_COL uint32 = 1
const OFF_SYS_COL uint32 = 2

type Collection struct {
	Name      string
	UID       uint32
	Computers []string
	sub       *subscription
}

func NewCollection(name string, uid uint32) *Collection {
	c := &Collection{
		Name:      name,
		UID:       uid,
		Computers: make([]string, 0),
		sub:       NewSubscription(),
	}
	return c
}

func (col *Collection) AddComputer(MAC string) {
	if col.ContainsComputer(MAC) {
		return
	}
	col.Computers = append(col.Computers, MAC)
	sort.Strings(col.Computers)
	rJSON, err := json.Marshal(WSResponse{"add-compR", CompStore.Get(MAC).Info})
	checkError(err)
	col.sub.Send(rJSON)
}

func (col *Collection) RemoveComputer(MAC string) bool {
	result := sort.SearchStrings(col.Computers, MAC)
	if col.Computers[result] == MAC {
		rJSON, err := json.Marshal(WSResponse{"rm-compR", MAC})
		checkError(err)
		col.sub.Send(rJSON)
		col.Computers = append(col.Computers[:result], col.Computers[result+1:]...)
		return true
	} else {
		return false
	}
}

func (col *Collection) ContainsComputer(MAC string) bool {
	result := sort.SearchStrings(col.Computers, MAC)
	if result >= len(col.Computers) {
		return false
	}
	if col.Computers[result] == MAC {
		return true
	} else {
		return false
	}
}

func (col *Collection) Subscribe(c *connection) {
	col.sub.Subscribe(c)
}

func (col *Collection) Unsubscribe(c *connection) {
	col.sub.Unsubscribe(c)
}

func (col *Collection) PrintComputers() {
	for _, c := range col.Computers {
		println(c)
	}
}
