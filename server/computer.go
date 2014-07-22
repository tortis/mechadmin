package main

import "github.com/tortis/mechadmin/types"

type Computer struct {
	Info types.Status
	UID  uint32
	Cols []uint32
}

func (comp *Computer) GenerateHTML() string {
	h := `<li class="computer-item">` + comp.Info.CN + `</li>`
	return h
}

func (comp *Computer) NotifyCollections() {
	for _, cid := range comp.Cols {
		ColStore.Get(cid).UpdateComputer(&comp.Info)
	}
}
