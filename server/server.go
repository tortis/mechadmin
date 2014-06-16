package main

import "net"
import "log"
import "encoding/json"
import "time"
import "github.com/tortis/cstatus/types"

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var ColStore *CollectionStore

func main() {
	var buffer []byte = make([]byte, 1024, 1024)

	addr, err := net.ResolveUDPAddr("udp", ":69")
	sock, err := net.ListenUDP("udp", addr)
	checkError(err)

	var s types.Status

	/*------------ Testing-----------------*/
	//root.NewCol(`Learning Spaces`)
	//lsOU := root.NewOU(`Learning Spaces`)
	//lsOU.NewCol(`COE`)
	//lsOU.NewCol(`COA`)

	//println("Creating collection named 'test'")
	//tCol := NewCollection("Test", uint32(1234))
	//println("Adding computer named 'COMP1'")
	//tCol.AddComputer("COMP1")
	//println("Adding computer named 'COMP2'")
	//tCol.AddComputer("COMP2")
	//println("Printing list of computers.")
	//tCol.PrintComputers()
	//println("Removing computer named 'COMP2'")
	//tCol.RemoveComputer("COMP2")
	//println()
	//println("Printing list of computers.")
	//tCol.PrintComputers()

	/*------------End Testing--------------*/

	go StartWebServer(root)
	for {
		println("Waiting for status packet.")
		rlen, _, err := sock.ReadFromUDP(buffer)
		checkError(err)

		json.Unmarshal(buffer[0:rlen], &s)
		UpdateComputer(&s)

		println("Computer: " + s.CN)
		println("User: " + s.UD + "\\" + s.UN)
		print("Active: ")
		println(s.A)
		println("Time: " + s.T.Format(time.ANSIC))
	}
}
