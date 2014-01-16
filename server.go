package xmpp

import (
	"net"
)

type XMPPServer struct {
	clients []*XMPPClient

	listener net.Listener

	handlers map[string]func(*XMPPClient, interface{}) error
}

func NewServer(laddr string) (*XMPPServer, error) {
	s := &XMPPServer{}
    s.handlers = make(map[string]func(*XMPPClient, interface{}) error)

	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		return s, err
	}

	s.listener = listener

	return s, nil
}

func (s *XMPPServer) AddHandler(tag string, h func(*XMPPClient, interface{}) error) {
    s.handlers[tag] = h
}

func (s *XMPPServer) ServeForever() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		c := NewClient(conn, &s.handlers)
		s.clients = append(s.clients, c)
	}
	return nil
}
