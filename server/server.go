package server

import (
    "github.com/zonyitoo/goxmpp/stream"
    "log"
    "net"
)

type Server interface {
    Accept() Client
    Serve()
}

type TCPServer struct {
    listener      net.Listener
    authenticator *stream.SASLAuthenticator
    shandler      stream.StanzaHandler
}

func NewTCPServer(listener net.Listener,
    a *stream.SASLAuthenticator, shandler stream.StanzaHandler) *TCPServer {
    return &TCPServer{
        listener:      listener,
        authenticator: a,
        shandler:      shandler,
    }
}

func (s *TCPServer) Accept() Client {
    conn, err := s.listener.Accept()
    if err != nil {
        panic(err)
    }
    return NewTCPClient(conn, s.authenticator, s.shandler)
}

func (s *TCPServer) Serve() {
    log.Printf("Server listening %+v", s.listener.Addr())
    for {
        c := s.Accept()
        log.Printf("Client %+v connected", c.Stream().RemoteAddr())
        go c.Run()
    }
}
