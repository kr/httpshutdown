package httpshutdown

import (
	"net"
	"sync"
)

// listener wraps a net.Listener and lets you wait when for
// all returned connections to be closed.
type listener struct {
	net.Listener
	w *sync.WaitGroup
}

func (l *listener) Accept() (net.Conn, error) {
	l.w.Add(1) // call Add before the event to be waited for (the connection)
	c, err := l.Listener.Accept()
	if err != nil {
		l.w.Done()
		return nil, err
	}
	return &conn{Conn: c, w: l.w}, nil
}

// conn wraps a net.Conn and decrements the WaitGroup
// when the connection is closed.
type conn struct {
	net.Conn
	w    *sync.WaitGroup
	once sync.Once
}

func (c *conn) Close() error {
	defer c.once.Do(c.w.Done)
	return c.Conn.Close()
}
