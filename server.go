package xmpp

import (
    // "encoding/xml"
    "log"
    "net"
)

type Server struct {
    dispatcher *ServerDispatcher
    listener   net.Listener
    config     *ServerConfig
}

func NewServer(conf *ServerConfig) *Server {
    listener, err := net.Listen("tcp", conf.ListenAddr)
    if err != nil {
        panic(err)
    }
    server := &Server{
        dispatcher: &ServerDispatcher{},
        listener:   listener,
        config:     conf,
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
        NewC2SStream(trans, s)
    }

    log.Println("Server exited")
}

func (s *Server) Run() {
    log.Printf("Server listening %s", s.config.ListenAddr)
    s.doAccept()
}
