package xmpp

import (
    "encoding/xml"
    "log"
    "net"
)

type C2SStream struct {
    transport     Transport
    decoder       *Decoder
    wchan         chan []byte
    isClosed      bool
    state         int
    jid           *JID
    server        *Server
    streamHandler StreamHandler
}

func NewC2SStream(trans Transport, s *Server) *C2SStream {
    stream := &C2SStream{
        transport:     trans,
        decoder:       NewDecoder(trans),
        wchan:         make(chan []byte),
        isClosed:      false,
        state:         STREAM_STAT_INIT,
        server:        s,
        streamHandler: s.config.StreamHandlerFactory(),
    }
    go stream.asyncWrite()
    go stream.asyncProcess()
    return stream
}

func (s *C2SStream) asyncWrite() {
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

func (s *C2SStream) asyncProcess() {
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
                s.SendElement(resp)
                s.EndStream()
            }
            return
        }

        switch t := elem.(type) {
        case xml.ProcInst:

        // case *XMPPStream:
        //     if s.state == STREAM_STAT_INIT {
        //         jid, err := NewJIDFromString(t.From)
        //         if err != nil {
        //             resp := &XMPPStreamError{
        //                 InvalidFrom: &XMPPStreamErrorInvalidFrom{},
        //             }
        //             s.SendErrorAndEnd(resp)
        //             return
        //         }
        //         s.jid = jid
        //         s.StartStream(STREAM_TYPE_CLIENT, "", t.From, t.Version, t.XmlLang)
        //         s.state |= STREAM_STAT_STARTED
        //     } else {
        //         resp := &XMPPStreamError{
        //             BadFormat: &XMPPStreamErrorBadFormat{},
        //         }
        //         s.SendErrorAndEnd(resp)
        //         return
        //     }
        case *XMPPStreamEnd:
            s.streamHandler.End(s)
            return
        default:
            s.Dispatch(t)
        }

        if s.state == STREAM_STAT_CLOSED {
            break
        }
    }
}

func (s *C2SStream) ServerConfig() *ServerConfig {
    return s.server.config
}

func (s *C2SStream) JID() *JID {
    return s.jid
}

func (s *C2SStream) SetJID(jid *JID) {
    s.jid = jid
}

func (s *C2SStream) Dispatch(elem interface{}) {
    switch t := elem.(type) {
    case *XMPPStream:
        s.streamHandler.Header(s, t)
        s.state = STREAM_STAT_STARTED
    case *XMPPStartTLS:
        s.state = STREAM_STAT_TLS_NEGOCIATION
        s.streamHandler.TLSNegociation(s, t)
    case *XMPPTLSProceed:
        s.streamHandler.TLSNegociation(s, t)
    case *XMPPTLSAbort:
        s.streamHandler.TLSNegociation(s, t)
    case *XMPPTLSFailure:
        s.streamHandler.TLSNegociation(s, t)
    case *XMPPSASLAuth:
        s.state = STREAM_STAT_SASL_NEGOCIATION
        s.streamHandler.SASLNegociation(s, t)
    case *XMPPSASLAbort:
        s.state = STREAM_STAT_SASL_FAILURE
        s.streamHandler.SASLNegociation(s, t)
    case *XMPPSASLFailure:
        s.state = STREAM_STAT_SASL_FAILURE
        s.streamHandler.SASLNegociation(s, t)
    case *XMPPSASLSuccess:
        s.state = STREAM_STAT_SASL_SUCCEED
        s.streamHandler.SASLNegociation(s, t)
    case *XMPPSASLChallenge:
        s.streamHandler.SASLNegociation(s, t)

    default:

    }
}

func (s *C2SStream) SendBytes(data []byte) {
    if len(data) > 0 && data[len(data)-1] != '\n' {
        data = append(data, '\n')
    }
    s.wchan <- data
}

func (s *C2SStream) StartStream(stype int, from, to, version, lang string) {
    header := &XMPPStream{
        From:    from,
        To:      to,
        Id:      s.transport.Id(),
        Version: version,
        XmlLang: lang,
    }
    if stype == STREAM_TYPE_CLIENT {
        header.Xmlns = XMLNS_JABBER_CLIENT
    } else {
        header.Xmlns = XMLNS_JABBER_SERVER
    }
    s.SendBytes([]byte(GenXMPPStreamHeader(header)))

}

func (s *C2SStream) sendStreamEnd() {
    s.SendBytes([]byte(stream_end_fmt))
}

func (s *C2SStream) SendElement(elem interface{}) error {
    b, err := xml.Marshal(elem)
    if err != nil {
        return err
    }
    s.SendBytes(b)
    return nil
}

func (s *C2SStream) Close() {
    if !s.isClosed {
        close(s.wchan)
    }
    s.isClosed = true
    s.state = STREAM_STAT_CLOSED
}

func (s *C2SStream) EndStream() {
    s.sendStreamEnd()
    s.Close()
}

func (s *C2SStream) EnableReadTimeout() error {
    return s.transport.SetReadTimeout()
}

func (s *C2SStream) DisableReadTimeout() error {
    return s.transport.UnsetReadTimeout()
}

func (s *C2SStream) State() int {
    return s.state
}

func (s *C2SStream) SetState(state int) {
    s.state = state
}
