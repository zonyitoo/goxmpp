package xmpp

import (
    "encoding/xml"
    "log"
    "net"
)

type ConnectionHandler func(*Connection) bool
type ProcessingHandler func(*Connection, interface{}) bool

type Server struct {
    listener    net.Listener
    connHandler ConnectionHandler
    procHandler ProcessingHandler
}

type Connection struct {
    server  *Server
    conn    net.Conn
    decoder *Decoder
    wchan   chan []byte
}

func NewServer(addr string, connHandler ConnectionHandler, procHandler ProcessingHandler) *Server {
    l, err := net.Listen("tcp", addr)
    if err != nil {
        log.Panic(err)
    }

    if procHandler == nil {
        log.Panic("ProcessingHandler cannot be nil")
    }
    return &Server{
        listener:    l,
        connHandler: connHandler,
        procHandler: procHandler,
    }
}

func (s *Server) Run() error {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            return nil
        }

        log.Printf("%s connected", conn.RemoteAddr().String())
        c := NewConnection(s, conn)

        if s.connHandler != nil && !s.connHandler(c) {
            conn.Close()
            continue
        }

        go c.process()
        go c.write()
    }
}

func NewConnection(s *Server, conn net.Conn) *Connection {
    return &Connection{
        server:  s,
        conn:    conn,
        decoder: NewDecoder(conn),
        wchan:   make(chan []byte),
    }
}

func (c *Connection) process() {
    for {
        elem, err := c.decoder.GetNextElement()
        if err != nil {
            log.Printf("%s Err: %s", c.conn.RemoteAddr().String(), err)
            c.EndStream()
            return
        }

        switch t := elem.(type) {
        case *XMPPStreamEnd:
            c.EndStream()
            return
        case xml.ProcInst:
            c.Write([]byte(xml.Header))
        default:
            if !c.server.procHandler(c, t) {
                c.EndStream()
                return
            }
        }
    }
}

func (c *Connection) write() {
    for b := range c.wchan {
        c.conn.Write(b)
    }
    log.Printf("%s closed", c.conn.RemoteAddr().String())
}

func (c *Connection) Close() {
    close(c.wchan)
    c.conn.Close()
}

func (c *Connection) Write(b []byte) {
    c.wchan <- b
}

func (c *Connection) EndStream() {
    c.Write([]byte("</stream:stream>"))
    c.Close()
}

func (c *Connection) RemoteAddr() net.Addr {
    return c.conn.RemoteAddr()
}
