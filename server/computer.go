package main

import "github.com/tortis/cstatus/types"
import "math/rand"

type Computer struct {
	Info types.Status
	UID  uint32
}

var allComputers map[string]*Computer = make(map[string]*Computer)

func init() {
	// Try to load computers from file.
}

func loadComputers() {

}

func saveComputers() {

}

func DeleteComputer(name string) {
	//for _,v := range allCollections {
	//	v.RemoveComputer(name)
	//}
	delete (allComputers,  name)
}

func UpdateComputer(stat *types.Status) *Computer {
	_, present := allComputers[stat.CN]
	if (present == true) {
		allComputers[stat.CN].Info = *stat
		return allComputers[stat.CN]
	} else {
		c := &Computer {
			Info:	*stat,
			UID:	rand.Uint32(),
		}
		allComputers[stat.CN] = c
		return c
	}
}

func (comp *Computer) GenerateHTML() string {
	h := `<li class="computer-item">` + comp.Info.CN + `</li>`
	return h
}
