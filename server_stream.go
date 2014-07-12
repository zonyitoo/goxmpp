package xmpp

import (
    "encoding/xml"
    log "github.com/cihub/seelog"
    "net"
)

type C2SStream struct {
    transport Transport
    decoder   *Decoder
    wchan     chan []byte
    isClosed  bool
    state     int
    jid       *JID
    server    *Server
}

func NewC2SStream(trans Transport, s *Server) *C2SStream {
    stream := &C2SStream{
        transport: trans,
        decoder:   NewDecoder(trans),
        wchan:     make(chan []byte),
        isClosed:  false,
        state:     STREAM_STAT_INIT,
        server:    s,
    }
    go stream.asyncWrite()
    go stream.asyncProcess()
    return stream
}

func (s *C2SStream) asyncWrite() {
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

func (s *C2SStream) asyncProcess() {
    for {
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
            return
        }

        log.Debugf("Received from %s with %+v", s.transport.RemoteAddr(), elem)

        if err := s.server.streamDispatcher.Dispatch(s, elem); err != nil {
            log.Error(err)
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

func (s *C2SStream) SendBytes(data []byte) {
    if len(data) > 0 && data[len(data)-1] != '\n' {
        data = append(data, '\n')
    }
    s.wchan <- data
}

func (s *C2SStream) StartStream(stype int, from, to, version, lang string) {
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

func (s *C2SStream) SendErrorAndClose(e interface{}) error {
    if s.State()&STREAM_STAT_STARTED == 0 {
        jidstr := ""
        if s.JID() != nil {
            jidstr = s.JID().String()
        }
        s.StartStream(STREAM_TYPE_CLIENT,
            s.ServerConfig().ServerName,
            jidstr,
            s.ServerConfig().StreamVersion.String(),
            "en")
    }
    err := s.SendElement(e)
    if err != nil {
        return err
    }
    s.EndStream()
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
