package xmpp

import (
    "encoding/xml"
    "log"
    "net"
)

const (
    STREAM_STAT_INIT    = 0
    STREAM_STAT_STARTED = 1 << iota
    STREAM_STAT_CLOSED
)

type Stream struct {
    transport Transport
    decoder   *Decoder
    wchan     chan []byte
    isClosed  bool
    state     int
    jid       *JID
}

func NewStream(trans Transport) *Stream {
    stream := &Stream{
        transport: trans,
        decoder:   NewDecoder(trans),
        wchan:     make(chan []byte),
        isClosed:  false,
        state:     STREAM_STAT_INIT,
    }
    go stream.asyncWrite()
    go stream.asyncProcess()
    return stream
}

func (s *Stream) ResetTransport(trans Transport) {
    if trans == nil {
        panic("Transport should not be nil")
    }

    s.transport = trans
    if !s.isClosed {
        close(s.wchan)
    }
    s.decoder = NewDecoder(s.transport)
    s.wchan = make(chan []byte)
    s.isClosed = false
    go s.asyncWrite()
    go s.asyncProcess()
}

func (s *Stream) asyncWrite() {
    for data := range s.wchan {
        _, err := s.transport.Write(data)

        if err != nil {
            log.Printf("%s %s\n", s.transport.RemoteAddr().String(), err)
            break
        }
    }

    if err := s.transport.Close(); err != nil {
        log.Println(err)
    }

    s.isClosed = true
    log.Printf("Client %s closed", s.transport.RemoteAddr().String())
}

func (s *Stream) asyncProcess() {
    for {
        elem, err := s.decoder.GetNextElement()
        if err != nil {
            switch err.(type) {
            case net.Error:
                s.Close()
            default:
                resp := &XMPPStreamError{}
                if err == DecoderBadFormat {
                    resp.BadFormat = &XMPPStreamErrorBadFormat{}
                } else {
                    resp.NotWellFormed = &XMPPStreamErrorNotWellFormed{}
                }
                s.SendErrorAndEnd(resp)
            }
            return
        }

        switch t := elem.(type) {
        case xml.ProcInst:
            if s.state != STREAM_STAT_INIT {
                resp := &XMPPStreamError{
                    BadFormat: &XMPPStreamErrorBadFormat{},
                }
                s.SendErrorAndEnd(resp)
                return
            } else {
                s.SendBytes([]byte(xml.Header))
            }
        case *XMPPStream:
            if s.state == STREAM_STAT_INIT {
                jid, err := NewJIDFromString(t.From)
                if err != nil {
                    resp := &XMPPStreamError{
                        InvalidFrom: &XMPPStreamErrorInvalidFrom{},
                    }
                    s.SendErrorAndEnd(resp)
                    return
                }
                s.jid = jid
                s.SendStreamHeader("", t.From, t.Version, t.XmlLang)
                s.state |= STREAM_STAT_STARTED
            } else {
                resp := &XMPPStreamError{
                    BadFormat: &XMPPStreamErrorBadFormat{},
                }
                s.SendErrorAndEnd(resp)
                return
            }
        case *XMPPStreamEnd:
            s.End()
            return
        default:
            if s.state&STREAM_STAT_STARTED == 0 {
                resp := &XMPPStreamError{
                    BadFormat: &XMPPStreamErrorBadFormat{},
                }
                s.SendErrorAndEnd(resp)
                return
            }
            s.SendElement(t)
        }
    }
}

func (s *Stream) SendBytes(data []byte) {
    s.wchan <- data
}

func (s *Stream) SendErrorAndEnd(e interface{}) {
    if s.state&STREAM_STAT_STARTED == 0 {
        s.SendStreamHeader("", "", "1.0", "en")
    }
    s.SendElement(e)
    s.End()
}

func (s *Stream) SendStreamHeader(from, to, version, lang string) {
    header := &XMPPStream{
        From:    from,
        To:      to,
        Id:      s.transport.Id(),
        Version: version,
        XmlLang: lang,
    }

    s.SendBytes([]byte(GenXMPPStreamHeader(header)))
}

func (s *Stream) SendStreamEnd() {
    s.SendBytes([]byte(stream_end_fmt))
}

func (s *Stream) SendElement(elem interface{}) error {
    b, err := xml.Marshal(elem)
    if err != nil {
        return err
    }
    s.SendBytes(b)
    return nil
}

func (s *Stream) Close() {
    if !s.isClosed {
        close(s.wchan)
    }
    s.isClosed = true
    s.state = STREAM_STAT_CLOSED
}

func (s *Stream) End() {
    s.SendStreamEnd()

    s.Close()
}

func (s *Stream) EnableReadTimeout() error {
    return s.transport.SetReadTimeout()
}

func (s *Stream) DisableReadTimeout() error {
    return s.transport.UnsetReadTimeout()
}
