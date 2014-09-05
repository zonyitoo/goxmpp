package stream

import (
    "code.google.com/p/go-uuid/uuid"
    "github.com/zonyitoo/goxmpp/protocol"
    "net"
)

type Streamer interface {
    Id() string
    Start() error
    RemoteAddr() net.Addr
    Writer() *Writer
    Reader() *Reader
    IsAnonymous() bool
    IsAuthenticated() bool
    SASLAuthenticator() *SASLAuthenticator
    Reset()
    Run()
    Close(bool) error
}

type ServerClientStream struct {
    conn            net.Conn
    id              string
    authenticator   *SASLAuthenticator
    writer          *Writer
    reader          *Reader
    isAnonymous     bool
    isAuthenticated bool
    stanzaHandler   StanzaHandler
}

func NewServerClientStream(conn net.Conn,
    authenticator *SASLAuthenticator, shandler StanzaHandler) *ServerClientStream {
    return &ServerClientStream{
        conn:            conn,
        id:              uuid.New(),
        reader:          NewReader(conn),
        writer:          NewWriter(conn),
        authenticator:   authenticator,
        isAuthenticated: false,
        isAnonymous:     false,
        stanzaHandler:   shandler,
    }
}

func (scs *ServerClientStream) Id() string {
    return scs.id
}

func (scs *ServerClientStream) Start() error {
    header, err := scs.Reader().NextElement()
    if err != nil {
        scs.Writer().SendElement(&protocol.XMPPStreamError{
            InvalidXML: &protocol.XMPPStreamErrorInvalidXML{},
        })
        scs.Close(true)
        return err
    }
    switch t := header.(type) {
    case *protocol.XMPPStream:
        if t.Xmlns != protocol.XMLNS_JABBER_CLIENT {
            scs.Writer().SendElement(&protocol.XMPPStreamError{
                InvalidNamespace: &protocol.XMPPStreamErrorInvalidNamespace{},
            })
            scs.Close(true)
            return err
        }
    default:
        scs.Writer().SendElement(&protocol.XMPPStreamError{
            BadFormat: &protocol.XMPPStreamErrorBadFormat{},
        })
        scs.Close(false)
        return err
    }
    scs.Writer().Open(&protocol.XMPPStream{
        Id:    scs.Id(),
        From:  "example.com",
        Xmlns: protocol.XMLNS_JABBER_CLIENT,
    })
    return nil
}

func (scs *ServerClientStream) RemoteAddr() net.Addr {
    return scs.conn.RemoteAddr()
}

func (scs *ServerClientStream) Writer() *Writer {
    return scs.writer
}

func (scs *ServerClientStream) Reader() *Reader {
    return scs.reader
}

func (scs *ServerClientStream) IsAnonymous() bool {
    return scs.isAnonymous
}

func (scs *ServerClientStream) IsAuthenticated() bool {
    return scs.isAuthenticated
}

func (scs *ServerClientStream) SASLAuthenticator() *SASLAuthenticator {
    return scs.authenticator
}

func (scs *ServerClientStream) Reset() {
    scs.reader = NewReader(scs.conn)
    scs.writer.Destroy()
    scs.writer = NewWriter(scs.conn)
}

func (scs *ServerClientStream) Run() {
    // Response Stream Header to Client
    if scs.Start() != nil {
        return
    }

    // TLS Negociation
    {
        // Send Feature
        tls := &protocol.XMPPStreamFeatures{
            StartTLS: &protocol.XMPPStartTLS{
                Required: &protocol.XMPPRequired{},
            },
        }
        scs.Writer().SendElement(tls)

        if resp, err := scs.Reader().NextElement(); err != nil {
            scs.Writer().SendElement(&protocol.XMPPStreamError{
                InvalidXML: &protocol.XMPPStreamErrorInvalidXML{},
            })
            scs.Close(true)
            return
        } else {
            switch resp.(type) {
            case *protocol.XMPPStartTLS:
                scs.Writer().SendElement(&protocol.XMPPTLSProceed{})
            case *protocol.XMPPTLSAbort:
                scs.Close(true)
                return
            default:
                scs.Writer().SendElement(&protocol.XMPPStreamError{
                    BadFormat: &protocol.XMPPStreamErrorBadFormat{},
                })
                scs.Close(false)
                return
            }
        }

        // TODO: Wrap TLS and Reset
        scs.Reset()
    }

    // Restart Stream
    if scs.Start() != nil {
        return
    }

    // TODO: Send Feature
    //       SASL Negociation
    {
        feature := &protocol.XMPPStreamFeatures{
            SASLMechanisms: &protocol.XMPPSASLMechanisms{
                Mechanisms: scs.authenticator.Mechanisms(),
            },
        }
        scs.Writer().SendElement(feature)

        // Waiting for XMPPSASLAuth
        if auth, err := scs.Reader().NextElement(); err != nil {
            scs.Writer().SendElement(&protocol.XMPPStreamError{
                InvalidXML: &protocol.XMPPStreamErrorInvalidXML{},
            })
            scs.Close(true)
            return
        } else {
            if t, ok := auth.(*protocol.XMPPSASLAuth); !ok {
                scs.Writer().SendElement(&protocol.XMPPStreamError{
                    BadFormat: &protocol.XMPPStreamErrorBadFormat{},
                })
                scs.Close(true)
                return
            } else {
                if !scs.authenticator.CallMechanism(t.Mechanism, t, scs) {
                    scs.Close(true)
                    return
                }

                scs.isAuthenticated = true
            }
        }

        // Reset
        scs.Reset()
    }

    // Restart Stream
    if scs.Start() != nil {
        return
    }

    // TODO: Send Features
    //       Send Optional Features
    features := &protocol.XMPPStreamFeatures{
        Bind: &protocol.XMPPBind{},
    }
    scs.Writer().SendElement(features)

    // TODO: Begin Stanza Exchanges
    for {
        elem, err := scs.Reader().NextElement()
        if err != nil {
            scs.Writer().SendElement(&protocol.XMPPStreamError{
                InvalidXML: &protocol.XMPPStreamErrorInvalidXML{},
            })
            scs.Close(true)
            return
        }

        switch t := elem.(type) {
        case *protocol.XMPPStanzaIQ:
            if scs.stanzaHandler.HandleIQ(t, scs) != nil {
                return
            }
        case *protocol.XMPPStanzaMessage:
            if scs.stanzaHandler.HandleMessage(t, scs) != nil {
                return
            }
        case *protocol.XMPPStanzaPresence:
            if scs.stanzaHandler.HandlePresence(t, scs) != nil {
                return
            }
        case *protocol.XMPPStreamEnd:
            scs.Close(true)
            return
        }
    }

    scs.Close(true)
}

func (scs *ServerClientStream) Close(withCloseTag bool) error {
    if withCloseTag {
        if err := scs.Writer().Close(); err != nil {
            return err
        }
    } else {
        if err := scs.Writer().Destroy(); err != nil {
            return err
        }
    }

    scs.conn.Close()
    return nil
}

type ServerServerStream struct {
    authenticator *SASLAuthenticator
}
