package xmpp

import (
	"errors"
	"log"
	"net"
)

type XMPPClientInfo struct {
}

type XMPPOptions struct {
	SASLMechanisms     []string
	SASLServerHandlers []XMPPSASLServerMechanismHandler
	SASLClientHandlers []XMPPSASLClientMechanismHandler
	IQHandler          func(iq *XMPPClientIQ) (interface{}, error)
	MessageHandler     func(msg *XMPPClientMessage) (interface{}, error)
	PresenceHandler    func(prsc *XMPPClientPresence) (interface{}, error)
	Handlers           map[string]func(interface{}) (interface{}, error)
	ServerURL          string
	ListenIP           string
}

func StreamHeaderDefaultHandler(c *XMPPClient, s *XMPPStream) error {
	log.Printf("C: %+v", s)
	if err := c.ResponseStreamHeader(s.To, s.From, "zh"); err != nil {
		return err
	}

	if c.State == STATE_RESTART {
		if err := c.Response(XMPPStreamFeatures{
			Bind: &XMPPBind{},
		}); err != nil {
			return err
		}
	} else {
		if err := c.Response(XMPPStreamFeatures{
			SASLMechanisms: &XMPPSASLMechanisms{
				Mechanisms: c.server.opts.SASLMechanisms,
			},
		}); err != nil {
			return err
		}
		c.State = STATE_SASL_AUTH
	}

	return nil
}

func StreamNegociationDefaultHandler(c *XMPPClient, elem interface{}) error {

	switch e := elem.(type) {
	case *XMPPSASLAuth:
		log.Printf("%+v", e)

		found := false
		for index, mec := range c.server.opts.SASLMechanisms {
			if mec == e.Mechanism {
				found = true
				c.srvHandler = c.server.opts.SASLServerHandlers[index]
				name, ret, err := c.srvHandler.Auth(e.Data)
				c.BindJID = name
				c.State = STATE_SASL_AUTH
				if err := c.Response(ret); err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					return err
				}
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					c.srvHandler = nil
					return err
				} else {
					_, e1 := ret.(*XMPPSASLSuccess)
					if e1 {
						c.State = STATE_RESTART
					}
					_, e2 := ret.(XMPPSASLSuccess)
					if e2 {
						c.State = STATE_RESTART
					}
				}

				break
			}
		}

		if !found {
			log.Printf("%+v Err: Invalid Mechanism %s", c.conn.RemoteAddr(), e.Mechanism)
			if err := c.Response(XMPPSASLFailure{
				InvalidMechanism: &XMPPSASLErrorInvalidMechanism{},
			}); err != nil {
				return err
			}
		}
	case *XMPPSASLResponse:
		if c.srvHandler == nil {
			log.Printf("%+v Err: Invalid <response/>", c.conn.RemoteAddr())
			if err := c.Response(XMPPSASLFailure{
				MalformedRequest: &XMPPSASLErrorMalformedRequest{},
			}); err != nil {
				return err
			}
		}
		name, ret, err := c.srvHandler.Response(e.Data)
		if err := c.Response(ret); err != nil {
			log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
			return err
		}
		if err != nil {
			log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
			c.srvHandler = nil
			return err
		} else {
			c.BindJID = name
		}

	default:
		log.Printf("%+v Err: Unsupported XML", c.conn.RemoteAddr())
		if err := c.Response(XMPPSASLFailure{
			MalformedRequest: &XMPPSASLErrorMalformedRequest{},
		}); err != nil {
			return err
		}
	}

	return nil
}

func StreamIQDefaultHandler(iq *XMPPClientIQ) (interface{}, error) {

	if iq.Type == "set" {
		if iq.Bind != nil {
			log.Printf("From client bind %+v", iq.Bind)
			if iq.Bind.Resource != "" {
				log.Printf("Client bind resource %s", iq.Bind.Resource)

				return &XMPPClientIQ{
					Bind: &XMPPBind{
						Jid: "abc@abc.com/" + iq.Bind.Resource,
					},
					Type: "result",
					Id:   iq.Id,
				}, nil
			} else {
				genr := "balcony"
				log.Printf("Client asked server for resournce %s", genr)
				return &XMPPClientIQ{
					Bind: &XMPPBind{
						Jid: "abc@abc.com/" + genr,
					},
					Type: "result",
					Id:   iq.Id,
				}, nil
			}
		}
	} else if iq.Type == "result" {

	}

	return &XMPPClientIQ{
		Error: &XMPPStanzaError{
			BadRequest: &XMPPStanzaErrorBadRequest{},
		},
	}, errors.New("Bad Request")
}

func StreamPresenceDefaultHandler(presence *XMPPClientPresence) (interface{}, error) {

	log.Printf("Client Present %+v", presence)

	return nil, nil
}

func StreamMessageDefaultHandler(msg *XMPPClientMessage) (interface{}, error) {
	log.Printf("Client message %+v", msg)

	return nil, nil
}

type XMPPServer struct {
	entities map[*XMPPClient]*XMPPClientInfo
	listener net.Listener
	opts     *XMPPOptions

	streamHandler      func(*XMPPClient, *XMPPStream) error
	negociationHandler func(*XMPPClient, interface{}) error
}

func NewServer(opts *XMPPOptions) (*XMPPServer, error) {
	s := &XMPPServer{
		entities:           make(map[*XMPPClient]*XMPPClientInfo),
		opts:               opts,
		streamHandler:      StreamHeaderDefaultHandler,
		negociationHandler: StreamNegociationDefaultHandler,
	}

	listener, err := net.Listen("tcp", opts.ListenIP)
	if err != nil {
		return s, err
	}

	s.listener = listener

	// Default Handlers
	//s.SetHandler(TAG_STREAM, StreamHeaderDefaultHandler)
	//s.SetHandler(TAG_SASL_AUTH, StreamNegociationDefaultHandler)
	//s.SetHandler(TAG_CLIENT_IQ, StreamIQDefaultHandler)
	//s.SetHandler(TAG_CLIENT_PRESENCE, StreamPresenceDefaultHandler)

	return s, nil
}

func (s *XMPPServer) ServeForever() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		c := NewClient(conn, s)
		s.entities[c] = &XMPPClientInfo{}
	}
	return nil
}
