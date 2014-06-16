package main

import "sort"

type Collection struct {
	Name      string
	UID       uint32
	Computers []string
}

func NewCollection(name string, uid uint32) *Collection {
	c := &Collection{
		Name:      name,
		UID:       uid,
		Computers: make([]string, 0),
	}
	return c
}

func (col *Collection) AddComputer(name string) {
	col.Computers = append(col.Computers, name)
	sort.Strings(col.Computers)
	wsHub.broadcast <- []byte("Adding new computer, " + name + " into collection, " + col.Name)
}

func (col *Collection) RemoveComputer(name string) bool {
	result := sort.SearchStrings(col.Computers, name)
	if col.Computers[result] == name {
		col.Computers = append(col.Computers[:result], col.Computers[result+1:]...)
		return true
	} else {
		return false
	}
}

func (col *Collection) ContainsComputer(name string) bool {
	if result := sort.SearchStrings(col.Computers, name); col.Computers[result] == name {
		return true;
	} else {
		return false;
	}
}

func (col *Collection) PrintComputers() {
	for _,c := range(col.Computers) {
		println(c)
	}
}
