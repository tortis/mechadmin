package main

import "time"
import "encoding/json"
import "github.com/tortis/mechadmin/types"

type Client struct {
	Name    string
	User    string
	Status  string
	Running bool
	MAC     string
	Resync  time.Duration
}

func NewClient(name, user, mac string) *Client {
	c := &Client{
		Name:    name,
		User:    user,
		Status:  "Hi, I'm new!",
		Running: true,
		Resync:  time.Second * 10,
		MAC:     mac,
	}
	return c
}

func (c *Client) start(send chan []byte, quit chan int) {
	var s types.Status
	for {
		s = types.Status{
			CN:  c.Name,
			UN:  c.User,
			UD:  c.Name,
			A:   c.Running,
			S:   c.Status,
			MAC: c.MAC,
			T:   time.Now(),
		}
		jsonBuffer, err := json.Marshal(s)
		checkError(err)
		println(string(jsonBuffer))
		send <- jsonBuffer
		time.Sleep(c.Resync)
	}
}
