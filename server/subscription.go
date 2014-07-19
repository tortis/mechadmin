package main

type subscription struct {
	subs	map[*connection]bool
}

func NewSubscription() *subscription {
	return &subscription{subs:make(map[*connection]bool)}
}

func (s *subscription) Subscribe(c *connection) {
	s.subs[c] = true
}

func (s *subscription) Unsubscribe(c *connection) {
	delete(s.subs, c)
}

func (s *subscription) Send(msg []byte) {
	for c,a := range s.subs {
		if a == true {
			select {	// Use select to prevent blocking on full channel.
			case c.send <- msg:
			default:
				delete(s.subs, c)
					// If the msg cant be sent, unSubscribe
			}
		}
	}
}
