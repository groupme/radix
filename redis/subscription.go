package redis

import "runtime"

//* Subscription

// Subscription is a structure for holding a Redis subscription for multiple channels.
type Subscription struct {
	client      *Client
	conn        *connection
	closerChan  chan struct{}
	msgHdlr     func(msg *Message)
}

// newSubscription returns a new Subscription or an error.
func newSubscription(client *Client, msgHdlr func(msg *Message)) (*Subscription, *Error) {
	var err *Error

	sub := &Subscription{
		client:      client,
		closerChan:  make(chan struct{}),
		msgHdlr:     msgHdlr,
	}

	// Connection handling
	sub.conn, err = sub.client.pool.pull()

	if err != nil {
		return nil, err
	}

	runtime.SetFinalizer(sub, (*Subscription).Close)
	go sub.backend()

	return sub, nil
}

// Subscribe subscribes to given channels or returns an error.
func (s *Subscription) Subscribe(channels ...string) *Error {
	return s.conn.subscribe(channels...)
}

// Unsubscribe unsubscribes from given channels or returns an error.
func (s *Subscription) Unsubscribe(channels ...string) *Error {
	return s.conn.unsubscribe(channels...)
}

// Psubscribe subscribes to given patterns or returns an error.
func (s *Subscription) Psubscribe(patterns ...string) *Error {
	return s.conn.psubscribe(patterns...)
}

// Punsubscribe unsubscribes from given patterns or returns an error.
func (s *Subscription) Punsubscribe(patterns ...string) *Error {
	return s.conn.punsubscribe(patterns...)
}

// Close closes the Subscription and returns its connection to the connection pool.
func (s *Subscription) Close() {
	runtime.SetFinalizer(s, nil)
	s.closerChan <- struct{}{}
	// Try to unsubscribe from all channels to reset the connection state back to normal
	err := s.conn.unsubscribe()
	if err != nil {
		s.conn.close()
		s.conn = nil
	}

	s.client.pool.push(s.conn)
}

func (s *Subscription) backend() {
	for {
		select {
		case <-s.closerChan:
			return
		case msg := <-s.conn.messageChan:
			s.msgHdlr(msg)
		}
	}
}
