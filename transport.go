package xmpp

import (
    "code.google.com/p/go-uuid/uuid"
    "io"
    "net"
    "time"
)

type Transport interface {
    io.Reader
    io.Writer
    RemoteAddr() net.Addr
    Close() error
    Id() string
    Conn() net.Conn
    SetReadTimeout() error
    UnsetReadTimeout() error
}

type TCPTransport struct {
    conn              net.Conn
    id                string
    hasSetReadTimeout bool
}

// As an io.Writer
func (trans *TCPTransport) Write(p []byte) (int, error) {
    return trans.conn.Write(p)
}

// As an io.Reader
func (trans *TCPTransport) Read(p []byte) (int, error) {
    cnt, err := trans.conn.Read(p)
    if err != nil {
        if e, ok := err.(net.Error); ok && e.Timeout() {
            // Write space ping
            if _, err := trans.Write([]byte(" ")); err != nil {
                return 0, err
            }
            trans.SetReadTimeout()
            return 0, nil
        }
    } else {
        if trans.hasSetReadTimeout {
            trans.SetReadTimeout()
        }
    }
    return cnt, err
}

func (trans *TCPTransport) RemoteAddr() net.Addr {
    return trans.conn.RemoteAddr()
}

func (trans *TCPTransport) Close() error {
    return trans.conn.Close()
}

func (trans *TCPTransport) Id() string {
    return trans.id
}

func (trans *TCPTransport) Conn() net.Conn {
    return trans.conn
}

func (trans *TCPTransport) SetReadTimeout() error {
    err := trans.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
    if err != nil {
        return err
    }
    trans.hasSetReadTimeout = true
    return nil
}

func (trans *TCPTransport) UnsetReadTimeout() error {
    err := trans.conn.SetReadDeadline(time.Time{})
    if err != nil {
        return err
    }
    trans.hasSetReadTimeout = false
    return nil
}

func NewTCPTransport(conn net.Conn) *TCPTransport {
    return &TCPTransport{
        conn:              conn,
        id:                uuid.New(),
        hasSetReadTimeout: false,
    }
}
