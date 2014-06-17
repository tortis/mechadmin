package main

import "net"
import "log"
import "encoding/json"
import "time"
import "github.com/tortis/mechadmin/types"
import "os"
import "os/signal"

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/* Global variables */
var ColStore *CollectionStore
var CompStore *ComputerStore
//var CompStore *ComputerStore
var root *CollectionTree
var wsHub = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func init() {
	ColStore = NewCollectionStore("colstore.gob")
	CompStore = NewComputerStore("compstore.gob")

	// Load root tree from file or create a new one.
	root = LoadTree("root.json")
	if root == nil {
		log.Fatal("Failed to create collection tree.")
	}

	// Register to recieve the interrup signal so that
	// the root tree can be saved before the program exits.
	// The root tree should still be saved at a regular interval.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if root.Save("root.json") != nil {
			println("\nError! The root collection tree could not be saved before the program stopped. Some data may have been lost.")
			os.Exit(1)
		} else {
			println("\nSuccessfully saved root collection tree.")
			os.Exit(0)
		}
	}()
}

func main() {
	var buffer []byte = make([]byte, 1024, 1024)

	addr, err := net.ResolveUDPAddr("udp", ":69")
	checkError(err)
	sock, err := net.ListenUDP("udp", addr)
	checkError(err)


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

	var s types.Status
	go StartWebServer()
	for {
		println("Waiting for status packet.")
		rlen, _, err := sock.ReadFromUDP(buffer)
		checkError(err)
		json.Unmarshal(buffer[0:rlen], &s)
		CompStore.UpdateOrAddComputer(&s)

		println("Computer: " + s.CN)
		println("User: " + s.UD + "\\" + s.UN)
		print("Active: ")
		println(s.A)
		println("Time: " + s.T.Format(time.ANSIC))
	}
}
