package main

import "net"
import "log"
import "os"
import "bytes"
import "encoding/json"
import "time"
import "os/exec"
import "regexp"
import "flag"
import "github.com/tortis/mechadmin/types"

var ip = flag.String("s", "localhost", "Specify the hostname or IP of the management server.")

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getActiveUser() (string, bool) {
	u := ""
	a := false
	session, _ := exec.Command("cmd", "/C", "query", "session").Output()
	// Pull result of "query session"
	exp, _ := regexp.Compile(`\s+`)
	// Create an expression that matches whitespace
	res := exp.Split(string(session), -1)
	// Split the session output on whitespace
	for index, element := range res { // Loop over the array
		if element == "Active" { // And if we find "Active"
			a = true
			u = res[index-2]
			break
		}
	}
	return u, a
}

func heartBeat(socket *net.UDPConn, host *net.UDPAddr, quit chan int, resync time.Duration) {
	var s types.Status
	var so string
	var jsonBuffer bytes.Buffer
	ifaces,err := net.Interfaces()
	checkError(err)
	mac := ifaces[1].HardwareAddr.String()
	println(len(ifaces))
	println(mac)
	lastSync := time.Now()
	enc := json.NewEncoder(&jsonBuffer)
	for {
		/* Get the easy stuff */
		s.UD = os.Getenv("USERDOMAIN")
		s.CN = os.Getenv("COMPUTERNAME")
		s.T = time.Now()
		s.UN, s.A = getActiveUser()
		s.MAC = mac
		/* Encode the status into JSON. */
		enc.Encode(&s)
		println(jsonBuffer.String())

		/* Check if the status has changed. */
		if so != jsonBuffer.String() || time.Since(lastSync) > resync {
			println("Sync required")
			s.T = time.Now()
			_, err := socket.WriteToUDP(jsonBuffer.Bytes(), host)
			checkError(err)
			lastSync = time.Now()
		} else {
			println("steady state")
		}

		/* Reset the encoding buffer. */
		so = jsonBuffer.String()
		jsonBuffer.Reset()
		time.Sleep(5 * time.Second)
	}
}

func main() {
	flag.Parse()
	println("Sending status packets to server at: "+*ip)
	host, err := net.ResolveUDPAddr("udp", *ip+":69")
	checkError(err)

	laddr, err := net.ResolveUDPAddr("udp", ":0")
	checkError(err)

	con, err := net.ListenUDP("udp", laddr)
	checkError(err)

	routineQuit := make(chan int)
	go heartBeat(con, host, routineQuit, time.Second*15)
	<-routineQuit
}
