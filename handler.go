package xmpp

import (
    "encoding/xml"
    "log"
)

type StanzaHandler interface {
    Message(Stream, *XMPPClientMessage)
    Present(Stream, *XMPPClientPresence)
    IQ(Stream, *XMPPClientIQ)
    Error(Stream, *XMPPStanzaError)
}

type StreamHandler interface {
    Header(Stream, *XMPPStream)
    TLSNegociation(Stream, interface{})
    SASLNegociation(Stream, interface{})
    Error(Stream, *XMPPStreamError)
    End(Stream)
}

type NegociationHandler interface {
}

type defaultStreamHandler struct {
    isAuthed bool
}

func DefaultStreamHandlerFactory() StreamHandler {
    return &defaultStreamHandler{
        isAuthed: false,
    }
}

func (h *defaultStreamHandler) Header(s Stream, x *XMPPStream) {
    if s.State() == STREAM_STAT_INIT {
        s.SendBytes([]byte(xml.Header))
    }
    var stype int
    if x.Xmlns == XMLNS_JABBER_CLIENT {
        stype = STREAM_TYPE_CLIENT
    } else if x.Xmlns == XMLNS_JABBER_SERVER {
        stype = STREAM_TYPE_SERVER
    } else {
        err := XMPPStreamError{
            InvalidNamespace: &XMPPStreamErrorInvalidNamespace{},
        }
        s.StartStream(XMLNS_JABBER_CLIENT, x.To, x.From, s.ServerConfig().StreamVersion.String(), x.XmlLang)
        s.SendElement(err)
        s.EndStream()
        return
    }

    if stype == STREAM_TYPE_CLIENT {
        if x.From != "" {
            if from_jid, err := NewJIDFromString(x.From); err != nil {
                goto CLIENT_FROM_ERROR
            } else if from_jid.Domain != x.To {
                goto CLIENT_FROM_ERROR
            }

            goto NO_FROM_ERROR
        CLIENT_FROM_ERROR:
            err := XMPPStreamError{
                InvalidFrom: &XMPPStreamErrorInvalidFrom{},
            }
            s.SendElement(err)
            s.EndStream()
            return
        }
    } else {
        if _, err := NewJIDFromString(x.From); err != nil {
            goto SERVER_FROM_ERROR
        }

        if _, err := NewJIDFromString(x.To); err != nil {
            goto SERVER_FROM_ERROR
        }

        goto NO_FROM_ERROR
    SERVER_FROM_ERROR:
        err := XMPPStreamError{
            InvalidFrom: &XMPPStreamErrorInvalidFrom{},
        }
        s.SendElement(err)
        s.EndStream()
        return
    }
NO_FROM_ERROR:

    if x.Version == "" {
        x.Version = "0.9"
    }

    version := &StreamVersion{}
    version.FromString(x.Version)

    if !s.ServerConfig().StreamVersion.GreaterOrEqualTo(version) {
        err := XMPPStreamError{
            UnsupportedVersion: XMPPStreamErrorUnsupportedVersion{},
        }
        s.StartStream(stype, x.To, x.From, s.ServerConfig().StreamVersion.String(), x.XmlLang)
        s.SendElement(err)
        s.EndStream()
    }

    s.StartStream(stype, x.To, x.From, s.ServerConfig().StreamVersion.String(), x.XmlLang)
}

func (h *defaultStreamHandler) TLSNegociation(s Stream, x interface{}) {

}

func (h *defaultStreamHandler) SASLNegociation(s Stream, x interface{}) {

}

func (h *defaultStreamHandler) Error(s Stream, err *XMPPStreamError) {

}

func (h *defaultStreamHandler) End(s Stream) {
    s.EndStream()
}
