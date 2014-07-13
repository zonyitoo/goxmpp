package xmpp

import (
    "encoding/xml"
    "errors"
    log "github.com/cihub/seelog"
    "net"
)

type ClientStream struct {
    transport Transport
    decoder   *Decoder
    wchan     chan []byte
    isClosed  bool
    state     int
    jid       *JID
    client    *Client
}

func NewClientStream(trans Transport, client *Client) *ClientStream {
    stream := &ClientStream{
        transport: trans,
        decoder:   NewDecoder(trans),
        wchan:     make(chan []byte),
        isClosed:  false,
        state:     STREAM_STAT_INIT,
        client:    client,
    }
    return stream
}

func (s *ClientStream) Run() {
    go s.asyncWrite()
    if err := s.initStream(); err != nil {
        if err != EventDispatcherIgnoredEventError {
            return
        }
    }
    go func() {
        for {
            if err := s.process(); err != nil {
                log.Error(err)
                break
            }
        }
    }()
}

func (s *ClientStream) asyncWrite() {
    for data := range s.wchan {
        _, err := s.transport.Write(data)

        if err != nil {
            log.Infof("%s %s\n", s.transport.RemoteAddr().String(), err)
            break
        }
    }

    if err := s.transport.Close(); err != nil {
        log.Error(err)
    }

    s.isClosed = true
    log.Infof("Client %s closed", s.transport.RemoteAddr().String())
}

func (s *ClientStream) initStream() error {
    s.client.streamDispatcher.DispatchEvent(EVENT_STREAM_INIT, s, nil)
    for {
        if err := s.process(); err != nil {
            log.Error(err)
            if err != EventDispatcherIgnoredEventError {
                return err
            }
        }

        if s.State()&STREAM_STAT_SASL_SUCCEEDED != 0 {
            break
        }
    }

    return nil
}

func (s *ClientStream) process() error {
    elem, err := s.decoder.GetNextElement()
    if err != nil {
        switch err.(type) {
        case net.Error:
            s.Close()
        default:
            resp := &XMPPStreamError{}
            if err == DecoderBadFormatError {
                resp.BadFormat = &XMPPStreamErrorBadFormat{}
            } else {
                resp.NotWellFormed = &XMPPStreamErrorNotWellFormed{}
            }
            s.SendErrorAndClose(resp)
        }
        log.Errorf("Decoding error: %s", err)
        return err
    }

    log.Debugf("Received from %s with %+v", s.transport.RemoteAddr(), elem)

    ev := s.client.streamDispatcher.EventOf(elem)

    if ev == EVENT_IMPOSSIBLE {
        s.SendErrorAndClose(&XMPPStreamError{
            InvalidXML: &XMPPStreamErrorInvalidXML{},
        })
        return errors.New("Impossible event")
    }

    if err := s.client.streamDispatcher.DispatchEvent(ev, s, elem); err != nil {
        log.Error(err)
    }

    if s.state == STREAM_STAT_CLOSED {
        return errors.New("Stream closed")
    }

    return nil
}

func (s *ClientStream) JID() *JID {
    return s.jid
}

func (s *ClientStream) SetJID(jid *JID) {
    s.jid = jid
}

func (s *ClientStream) SendBytes(data []byte) {
    if len(data) > 0 && data[len(data)-1] != '\n' {
        data = append(data, '\n')
    }
    s.wchan <- data
}

func (s *ClientStream) StartStream(stype int, from, to, version, lang string) {
    if s.State() == STREAM_STAT_INIT {
        s.SendBytes([]byte(xml.Header))
    }
    header := &XMPPStream{
        From:    from,
        To:      to,
        Id:      s.transport.Id(),
        Version: version,
        XMLLang: lang,
    }
    if stype == STREAM_TYPE_CLIENT {
        header.Xmlns = XMLNS_JABBER_CLIENT
    } else {
        header.Xmlns = XMLNS_JABBER_SERVER
    }
    s.SendBytes([]byte(GenXMPPStreamHeader(header)))

}

func (s *ClientStream) sendStreamEnd() {
    s.SendBytes([]byte(stream_end_fmt))
}

func (s *ClientStream) SendElement(elem interface{}) error {
    log.Debugf("Send %+v", elem)
    b, err := xml.Marshal(elem)
    if err != nil {
        return err
    }
    s.SendBytes(b)
    return nil
}

func (s *ClientStream) SendErrorAndClose(e interface{}) error {
    if s.State()&STREAM_STAT_STARTED == 0 {
        jidstr := ""
        if s.JID() != nil {
            jidstr = s.JID().String()
        }
        s.StartStream(STREAM_TYPE_CLIENT,
            s.Config().ServerName,
            jidstr,
            s.Config().StreamVersion.String(),
            "en")
    }
    err := s.SendElement(e)
    if err != nil {
        return err
    }
    s.EndStream()
    return nil
}

func (s *ClientStream) Config() *Config {
    return s.client.config
}

func (s *ClientStream) Close() {
    if !s.isClosed {
        close(s.wchan)
    }
    s.isClosed = true
    s.state = STREAM_STAT_CLOSED
}

func (s *ClientStream) EndStream() {
    s.sendStreamEnd()
    s.Close()
}

func (s *ClientStream) EnableReadTimeout() error {
    return s.transport.SetReadTimeout()
}

func (s *ClientStream) DisableReadTimeout() error {
    return s.transport.UnsetReadTimeout()
}

func (s *ClientStream) State() int {
    return s.state
}

func (s *ClientStream) SetState(state int) {
    s.state = state
}
