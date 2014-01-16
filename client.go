package xmpp

import (
	"bufio"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

type XMPPClient struct {
	incoming chan []byte
	outgoing chan []byte
	conn     net.Conn

	bufrw *bufio.ReadWriter

	xmlDecoder *xml.Decoder

	handlers *map[string]func(*XMPPClient, interface{}) error
}

func (c *XMPPClient) Write() {
	for data := range c.outgoing {
		_, err := c.bufrw.Write(data)
		if err != nil {
			c.CloseStream()
			log.Fatal(err)
			break
		}
        c.bufrw.Flush()
		log.Printf("S: %s\n", string(data))
	}
}

func (c *XMPPClient) Read() {
	if err := c.DecodeXMLStreamElements(); err != nil {
		log.Fatal(err)
	}
	c.CloseStream()
}

func NewClient(conn net.Conn, handler *map[string]func(*XMPPClient, interface{}) error) *XMPPClient {
	bufrw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	c := &XMPPClient{
		incoming:   make(chan []byte),
		outgoing:   make(chan []byte),
		conn:       conn,
		bufrw:      bufrw,
		xmlDecoder: xml.NewDecoder(bufrw),
		handlers:   handler,
	}
	go c.Read()
	go c.Write()
	return c
}

func (c *XMPPClient) CloseStream() error {
	c.outgoing <- []byte(stream_end_fmt)
	return c.conn.Close()
}

func (c *XMPPClient) Response(obj interface{}) error {
	resp, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	c.outgoing <- resp
	return nil
}

func (c *XMPPClient) CallHandler(tag string, arg interface{}) error {

	if handler, ok := (*c.handlers)[tag]; ok {
		return handler(c, arg)
	}

	return nil
}

func (c *XMPPClient) ResponseStreamHeader(from, to, langcode string) error {

	sheader := &XMPPStream{
		From:    from,
		To:      to,
		XmlLang: langcode,
		Version: "1.0",
		Id:      generate_random_id(),
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
					ErrorType: XMPPStreamErrorInvalidNamespace{},
				}
			}
		} else if attr.Name.Space == "xmlns" && attr.Name.Local == "stream" {
			if attr.Value != XMLNS_STREAM {
				streamError = &XMPPStreamError{
					ErrorType: XMPPStreamErrorInvalidNamespace{},
				}
			}
		} else if attr.Name.Local == "version" {
			if attr.Value != "1.0" {
				streamError = &XMPPStreamError{
					ErrorType: XMPPStreamErrorUnsupportedVersion{},
				}
			}
			recv_stream.Version = attr.Value
		} else if attr.Name.Space == "xml" && attr.Name.Local == "lang" {
			recv_stream.XmlLang = attr.Value
		}
	}

	return recv_stream, streamError
}

func (c *XMPPClient) DecodeXMLStreamElements() error {
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
                    tag = HANDLER_SASL_AUTH
				case "response":
					element = &XMPPSASLResponse{}
                    tag = HANDLER_SASL_RESPONSE
				case "iq":
					element = &XMPPClientIQ{}
                    tag = HANDLER_CLIENT_IQ
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
	return nil
}

func (c *XMPPClient) nextStartElement() (xml.StartElement, error) {
	for {
		token, err := c.xmlDecoder.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			return t, nil
		}
	}
}

func (c *XMPPClient) DecodeElement() (interface{}, error) {

	var elem interface{}

	token, err := c.nextStartElement()
	if err != nil {
		return elem, err
	}

	switch token.Name.Space + " " + token.Name.Local {
	case "stream features":
		elem = &XMPPStreamFeatures{}
	case XMLNS_XMPP_SASL + " auth":
		elem = &XMPPSASLAuth{}
	case XMLNS_XMPP_SASL + " challenge":
		elem = &XMPPSASLChallenge{}
	case XMLNS_XMPP_SASL + " response":
		elem = &XMPPSASLResponse{}

	}

	err = c.xmlDecoder.DecodeElement(elem, &token)

	return elem, err
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
