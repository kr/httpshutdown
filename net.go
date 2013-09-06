package httpshutdown

import (
	"net"
	"sync"
)

// listener wraps a net.Listener and lets you wait when for
// all returned connections to be closed.
type listener struct {
	net.Listener
	w sync.WaitGroup
}

func (l *listener) Accept() (c net.Conn, e error) {
	c, e = l.Listener.Accept()
	if e == nil {
		l.w.Add(1)
		c = &conn{Conn: c, w: &l.w}
	}
	return
}

func (l *listener) wait() {
	l.w.Wait()
}

// conn wraps a net.Conn and decrements the WaitGroup
// when the connection is closed.
type conn struct {
	net.Conn
	w    *sync.WaitGroup
	once sync.Once
}

func (c *conn) Read(p []byte) (n int, e error) {
	n, e = c.Conn.Read(p)
	if ne, ok := e.(net.Error); e != nil && !(ok && ne.Temporary()) {
		c.once.Do(c.w.Done)
	}
	return
}

func (c *conn) Write(p []byte) (n int, e error) {
	n, e = c.Conn.Write(p)
	if ne, ok := e.(net.Error); e != nil && !(ok && ne.Temporary()) {
		c.once.Do(c.w.Done)
	}
	return
}

func (c *conn) Close() error {
	c.once.Do(c.w.Done)
	return c.Conn.Close()
}
