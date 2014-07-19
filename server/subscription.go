package main

type subscription struct {
	Subs	map[*connection]bool
}

func NewSubscription() *subscription {
	return &subscription{Subs:make(map[*connection]bool)}
}

func (s *subscription) Subscribe(c *connection) {
	s.Subs[c] = true
}

func (s *subscription) Unsubscribe(c *connection) {
	delete(s.Subs, c)
}

func (s *subscription) Send(msg []byte) {
	for c,a := range s.Subs {
		if a == true {
			select {	// Use select to prevent blocking on full channel.
			case c.send <- msg:
			default:
				delete(s.Subs, c)
					// If the msg cant be sent, unSubscribe
			}
		}
	}
}
