package xmpp

import (
    // "encoding/xml"
    "log"
    "net"
)

type Server struct {
    clients         map[BareJID][]*Stream
    listener        net.Listener
    config          *ServerConfig
    unbindedClients []*Stream
}

func NewServer(conf *ServerConfig) *Server {
    listener, err := net.Listen("tcp", conf.listenAddr)
    if err != nil {
        panic(err)
    }
    server := &Server{
        clients:  make(map[BareJID][]*Stream),
        listener: listener,
        config:   conf,
    }
    return server
}

func (s *Server) doAccept() {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            log.Println(err)
            break
        }

        log.Printf("Client %s connected", conn.RemoteAddr().String())

        trans := NewTCPTransport(conn)
        trans.SetReadTimeout()
        NewStream(trans)
        // s.unbindedClients = append(s.unbindedClients, stream)
    }

    log.Println("Server exited")
}

func (s *Server) Run() {
    log.Printf("Server listening %s", s.config.listenAddr)
    s.doAccept()
}
