package xmpp

import (
    log "github.com/cihub/seelog"
    "net"
)

type Server struct {
    streamDispatcher *StreamEventDispatcher
    listener         net.Listener
    config           *ServerConfig
}

func NewServer(conf *ServerConfig) *Server {
    listener, err := net.Listen("tcp", conf.ListenAddr)
    if err != nil {
        panic(err)
    }
    server := &Server{
        streamDispatcher: NewStreamEventDispatcher(),
        listener:         listener,
        config:           conf,
    }
    server.addDefaultHandlers()
    return server
}

func (server *Server) addDefaultHandlers() {
    server.AddHandlerForEvent(EVENT_STREAM_HEADER, DefaultStreamHeaderHandler)
    server.AddHandlerForEvent(EVENT_STREAM_END, DefaultStreamEndHandler)
}

func (s *Server) doAccept() {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            log.Error(err)
            break
        }

        log.Infof("Client %s connected", conn.RemoteAddr().String())

        trans := NewTCPTransport(conn)
        trans.SetReadTimeout()
        NewC2SStream(trans, s)
    }

    log.Info("Server exited")
}

func (s *Server) Run() {
    log.Infof("Server listening %s", s.config.ListenAddr)
    s.doAccept()
}

func (s *Server) AddHandlerForEvent(ev int, hdl XMPPEventHandler) {
    s.streamDispatcher.AddHandlerForEvent(ev, hdl)
}
