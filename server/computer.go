package main

import "github.com/tortis/mechadmin/types"

type Computer struct {
	Info types.Status
	UID  uint32
}

func (comp *Computer) GenerateHTML() string {
	h := `<li class="computer-item">` + comp.Info.CN + `</li>`
	return h
}
