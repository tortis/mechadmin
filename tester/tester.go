package main

import "net"
import "log"
import "flag"
import "math/rand"
import "time"
import "strconv"

var ip = flag.String("s", "localhost", "Specify the hostname or IP of the management server.")

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func sender(toSend chan []byte, stop chan int, s *net.UDPConn, host *net.UDPAddr) {
	for {
		select {
		case msg := <-toSend:
			s.WriteToUDP(msg, host)
		case <-stop:
			return
		}
	}
}

func nameGen() string {
	b := make([]byte, 0)
	len := rand.Int()%10 + 5
	println(len)
	for c := len; c > 0; c-- {
		b = append(b, byte(rand.Int()%26+97))
	}
	return string(b[0:len])
}

func macGen() string {
	r := ""
	for i := 0; i < 6; i++ {
		r += strconv.FormatInt(rand.Int63()%255, 16)
		if (i < 5) {
			r += ":"
		}
	}
	return r
}

func main() {
	flag.Parse()
	rand.Seed(42)
	println("Sending status packets to server at: " + *ip)

	host, err := net.ResolveUDPAddr("udp", *ip+":69")
	checkError(err)
	laddr, err := net.ResolveUDPAddr("udp", ":0")
	checkError(err)
	con, err := net.ListenUDP("udp", laddr)
	checkError(err)

	routineQuit := make(chan int)
	send := make(chan []byte, 10)
	go sender(send, routineQuit, con, host)

	for num := 0; num < 100; num++ {
		go NewClient(nameGen(), nameGen(), macGen()).start(send, routineQuit)
		time.Sleep(time.Second * 2)
	}
	<-routineQuit
}
