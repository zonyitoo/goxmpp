package xmpp

import (
	"bufio"
	"encoding/xml"
	//	"errors"
	//	"io"
	"log"
	"net"
	"strings"
)

const (
	STATE_INIT = iota
	STATE_START_TLS
	STATE_SASL_AUTH
	STATE_SASL_AUTH_CHALLENGE
	STATE_SASL_AUTH_RESPONSE
	STATE_SASL_AUTH_DONE
	STATE_RESTART
	STATE_RESOURCE_BINDING
)

type XMPPClient struct {
	incoming   chan []byte
	outgoing   chan []byte
	conn       net.Conn
	bufrw      *bufio.ReadWriter
	xmlDecoder *xml.Decoder
	State      int
	BindJID    string
	srvHandler XMPPSASLServerMechanismHandler
	cliHandler XMPPSASLClientMechanismHandler
	server     *XMPPServer
	Id         string
}

func (c *XMPPClient) Write() {
	for data := range c.outgoing {
		_, err := c.bufrw.Write(data)
		if err != nil {
			c.CloseStream()
			log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
			break
		}
		c.bufrw.Flush()
		log.Printf("%+v Resp: %s\n", c.conn.RemoteAddr(), string(data))
	}
}

func (c *XMPPClient) Process() {
PROCESS_LOOP:
	for {
		token, err := c.xmlDecoder.Token()
		if err != nil {
			log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			tag, elem, serr := c.decodeXMLStreamElements(&t)
			if serr != nil {
				c.Response(serr)
				c.CloseStream()
				break PROCESS_LOOP
			}

			log.Printf("%+v Recv: %s %+v", c.conn.RemoteAddr(), tag, elem)

			switch t := elem.(type) {
			case *XMPPStream:
				if err := c.server.streamHandler(c, t); err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
				}
			case *XMPPSASLAuth:
				err := c.server.negociationHandler(c, t)
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					break PROCESS_LOOP
				}
			case *XMPPSASLChallenge:
				err := c.server.negociationHandler(c, t)
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					break PROCESS_LOOP
				}
			case *XMPPSASLResponse:
				err := c.server.negociationHandler(c, t)
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					break PROCESS_LOOP
				}
			case *XMPPSASLSuccess:
				err := c.server.negociationHandler(c, t)
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					break PROCESS_LOOP
				}
			case *XMPPSASLFailure:
				err := c.server.negociationHandler(c, t)
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					break PROCESS_LOOP
				}
			case *XMPPClientIQ:
				ret, err := c.server.opts.IQHandler(t)
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
				}
				if err = c.Response(ret); err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					break PROCESS_LOOP
				}
			case *XMPPClientMessage:
				ret, err := c.server.opts.MessageHandler(t)
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
				}
				if err = c.Response(ret); err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					break PROCESS_LOOP
				}
			case *XMPPClientPresence:
				ret, err := c.server.opts.PresenceHandler(t)
				if err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
				}
				if err = c.Response(ret); err != nil {
					log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					break PROCESS_LOOP
				}
			default:
				if h, ok := c.server.opts.Handlers[tag]; ok {
					ret, err := h(t)
					if err != nil {
						log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
					}
					if err = c.Response(ret); err != nil {
						log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
						break PROCESS_LOOP
					}
				} else {
					log.Printf("%+v Err: Unrecognized XML %s", c.conn.RemoteAddr(), tag)
					if err := c.Response(XMPPStreamError{
						UndefinedCondition: &XMPPStreamErrorUndefinedCondition{},
					}); err != nil {

						log.Printf("%+v Err: %s", c.conn.RemoteAddr(), err)
						break PROCESS_LOOP
					}
				}
			}
		case xml.ProcInst:
			if serr := c.processXMLProcInst(&t); serr != nil {
				c.ResponseStreamHeader("", "", "en")
				c.Response(serr)
				c.CloseStream()
				break PROCESS_LOOP
			}
		case xml.EndElement:
			if t.Name.Local == "stream" && t.Name.Space == "stream" {
				c.CloseStream()
				break PROCESS_LOOP
			} else {
				log.Printf("%+v Unexpected EndElement: %+v",
					c.conn.RemoteAddr(), t)
			}
		}
	}

	c.CloseStream()
}

func NewClient(conn net.Conn, server *XMPPServer) *XMPPClient {
	bufrw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	c := &XMPPClient{
		incoming:   make(chan []byte),
		outgoing:   make(chan []byte),
		conn:       conn,
		bufrw:      bufrw,
		xmlDecoder: xml.NewDecoder(bufrw),
		State:      STATE_INIT,
		server:     server,
		Id:         generate_random_id(),
	}
	go c.Process()
	go c.Write()
	return c
}

func (c *XMPPClient) CloseStream() error {
	c.outgoing <- []byte(stream_end_fmt)
    delete(c.server.entities, c)
	return c.conn.Close()
}

func (c *XMPPClient) Response(obj interface{}) error {
    if obj == nil {
        // Blankspace ping
        c.outgoing <- []byte(" ")
    }
	resp, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	c.outgoing <- resp
	return nil
}

/*
func (c *XMPPClient) CallHandler(tag string, arg interface{}) error {

	if handler, ok := (c.server.handlers)[tag]; ok {
		return handler(c, arg)
	}

	return nil
}
*/

func (c *XMPPClient) ResponseStreamHeader(from, to, langcode string) error {

	sheader := &XMPPStream{
		From:    from,
		To:      to,
		XmlLang: langcode,
		Version: "1.0",
		Id:      c.Id,
	}

	c.outgoing <- []byte(make_stream_begin(sheader))

	return nil
}

func (c *XMPPClient) decodeStreamHeader(t *xml.StartElement) (*XMPPStream, *XMPPStreamError) {
	var streamError *XMPPStreamError = nil
	recv_stream := &XMPPStream{}
	for _, attr := range t.Attr {
		if attr.Name.Local == "from" {
			recv_stream.From = attr.Value
		} else if attr.Name.Local == "to" {
			recv_stream.To = attr.Value
		} else if attr.Name.Local == "xmlns" {
			if attr.Value != XMLNS_JABBER_CLIENT {
				streamError = &XMPPStreamError{
					InvalidNamespace: &XMPPStreamErrorInvalidNamespace{},
				}
			}
		} else if attr.Name.Space == "xmlns" && attr.Name.Local == "stream" {
			if attr.Value != XMLNS_STREAM {
				streamError = &XMPPStreamError{
					InvalidNamespace: &XMPPStreamErrorInvalidNamespace{},
				}
			}
			recv_stream.XMLName.Space = attr.Value
			recv_stream.XMLName.Local = "stream"
		} else if attr.Name.Local == "version" {
			if attr.Value != "1.0" {
				streamError = &XMPPStreamError{
					UnsupportedVersion: &XMPPStreamErrorUnsupportedVersion{},
				}
			}
			recv_stream.Version = attr.Value
		} else if attr.Name.Space == "xml" && attr.Name.Local == "lang" {
			recv_stream.XmlLang = attr.Value
		}
	}

	return recv_stream, streamError
}

func (c *XMPPClient) decodeXMLElement(t *xml.StartElement, obj interface{}) *XMPPStreamError {
	if err := c.xmlDecoder.DecodeElement(obj, t); err != nil {
		serr := &XMPPStreamError{
			BadFormat: &XMPPStreamErrorBadFormat{},
		}

		return serr
	}

	return nil
}

func (c *XMPPClient) processXMLProcInst(t *xml.ProcInst) *XMPPStreamError {
	var streamError *XMPPStreamError = nil
	if t.Target == "xml" {
		insts := strings.Split(string(t.Inst), " ")

		for _, inst := range insts {
			parts := strings.Split(inst, "=")
			if len(parts) == 2 {
				if parts[0] == "encoding" && strings.ToUpper(parts[1]) != "'UTF-8'" {
					streamError = &XMPPStreamError{
						UnsupportedEncoding: &XMPPStreamErrorUnsupportedEncoding{},
					}
					if err := c.Response(streamError); err != nil {
						c.CloseStream()
					}
				}
			}
		}
		c.outgoing <- []byte(xml.Header)
	}

	return streamError
}

func (c *XMPPClient) decodeXMLStreamElements(token *xml.StartElement) (string, interface{}, *XMPPStreamError) {
	var tag string = token.Name.Space + ":" + token.Name.Local
	if token.Name.Space == "" {
		for _, attr := range token.Attr {
			if attr.Name.Local == "xmlns" {
				tag = attr.Value + ":" + token.Name.Local
				break
			}
		}
	} else {
		for _, attr := range token.Attr {
			if attr.Name.Space == "xmlns" && attr.Name.Local == token.Name.Space {
				tag = attr.Value + ":" + token.Name.Local
				break
			}
		}
	}

	var obj interface{}
	switch tag {
	case TAG_STREAM:
		s, serr := c.decodeStreamHeader(token)
		return TAG_STREAM, s, serr
	case TAG_STREAM_FEATURES:
		obj = &XMPPStreamFeatures{}
	case TAG_TLS_START:
		obj = &XMPPStartTLS{}
	case TAG_TLS_PROCEED:
		obj = &XMPPTLSProceed{}
	case TAG_TLS_FAILURE:
		obj = &XMPPTLSFailure{}

	case TAG_SASL_AUTH:
		obj = &XMPPSASLAuth{}
	case TAG_SASL_CHALLENGE:
		obj = &XMPPSASLChallenge{}
	case TAG_SASL_RESPONSE:
		obj = &XMPPSASLResponse{}
	case TAG_SASL_SUCCESS:
		obj = &XMPPSASLSuccess{}
	case TAG_SASL_FAILURE:
		obj = &XMPPSASLFailure{}
	case TAG_CLIENT_IQ:
		obj = &XMPPClientIQ{}
	case TAG_CLIENT_PRESENCE:
		obj = &XMPPClientPresence{}
	case TAG_CLIENT_MESSAGE:
		obj = &XMPPClientMessage{}
	default:
		log.Printf("Cannot decode token: %+v", token)
		return "", "", &XMPPStreamError{
			NotWellFormed: &XMPPStreamErrorNotWellFormed{},
		}
	}

	return tag, obj, c.decodeXMLElement(token, obj)
	/*
		for {
			token, err := c.xmlDecoder.RawToken()
			if err != nil {
				return err
			} else if err == io.EOF {
				return err
			}

			log.Printf("DecodeXMLStreamElements: %+v", token)

			var streamError *XMPPStreamError = nil
			switch t := token.(type) {
			case xml.ProcInst:
				if t.Target == "xml" {
					insts := strings.Split(string(t.Inst), " ")

					for _, inst := range insts {
						parts := strings.Split(inst, "=")
						if len(parts) == 2 {
							if parts[0] == "encoding" && strings.ToUpper(parts[1]) != "'UTF-8'" {
								streamError = &XMPPStreamError{
									ErrorType: XMPPStreamErrorUnsupportedEncoding{},
								}
								if err := c.Response(streamError); err != nil {
									return err
								}
								return errors.New("Unsupported Encoding")
							}
						}
					}
					c.outgoing <- []byte(xml.Header)
				}
			case xml.StartElement:
				if t.Name.Local == "stream" {

					streamHeader, serr := c.decodeStreamHeader(&t)

					if serr != nil {
						c.Response(serr)
						return errors.New("Stream Header Decode Error")
					}

					if err := c.CallHandler("stream:stream", streamHeader); err != nil {
						return err
					}
				} else {
					var element interface{} = nil

					tag := t.Name.Local
					switch tag {
					case "starttls":
						element = &XMPPStartTLS{}
					case "auth":
						element = &XMPPSASLAuth{}
						tag = TAG_SASL_AUTH
					case "response":
						element = &XMPPSASLResponse{}
						tag = TAG_SASL_RESPONSE
					case "iq":
						element = &XMPPClientIQ{}
						tag = TAG_CLIENT_IQ
					default:
						continue
					}

					if err := c.xmlDecoder.DecodeElement(element, &t); err != nil {
						return err
					}

					if err := c.CallHandler(tag, element); err != nil {
						return err
					}
				}
			case xml.EndElement:
				if t.Name.Local == "stream" {
					return errors.New("Stream Closed")
				}
			}
		}
	*/
}

/*
func (c *XMPPClient) Negociation() error {
	features := XMPPStreamFeatures{
		SASLMechanisms: &XMPPSASLMechanisms{
			Mechanisms: []string{
				"PLAIN",
			},
		},
	}
	if err := c.Response(features); err != nil {
		c.Close()
		return err
	}
	// Expecting the Auth
	auth := XMPPSASLAuth{}
	c.xmlDecoder.DecodeElement(&auth, nil)
	log.Printf("%+v", auth)
	resp, serr := c.Handler(auth)
	if serr != nil {
		if err := c.Response(*serr); err != nil {
			c.Close()
			return err
		}
		c.Close()
		return errors.New("Stream Error")
	}
	if err := c.Response(resp); err != nil {
		c.Close()
		return err
	}
	// Success
	auth_success := XMPPSASLSuccess{}
	resp_auth_success, err := xml.Marshal(auth_success)
	if err != nil {
		panic(err)
	}
	c.outgoing <- resp_auth_success
	stream_token, err := c.xmlDecoder.Token()
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", stream_token)

	return nil
}
*/
