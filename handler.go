package xmpp

import (
    "encoding/xml"
    // "log"
)

type XMPPEventHandler func(Stream, interface{}) bool

func DefaultStreamHeaderHandler(s Stream, _x interface{}) bool {
    x, _ := _x.(*XMPPStream)
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
        s.StartStream(STREAM_TYPE_CLIENT, x.To, x.From, s.ServerConfig().StreamVersion.String(), x.XMLLang)
        s.SendElement(err)
        s.EndStream()
        return true
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
            return true
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
        return true
    }
NO_FROM_ERROR:

    if x.Version == "" {
        x.Version = "0.9"
    }

    version := &StreamVersion{}
    version.FromString(x.Version)

    if !s.ServerConfig().StreamVersion.GreaterOrEqualTo(version) {
        err := XMPPStreamError{
            UnsupportedVersion: &XMPPStreamErrorUnsupportedVersion{},
        }
        s.StartStream(stype, x.To, x.From, s.ServerConfig().StreamVersion.String(), x.XMLLang)
        s.SendElement(err)
        s.EndStream()
    }

    s.StartStream(stype, x.To, x.From, s.ServerConfig().StreamVersion.String(), x.XMLLang)
    return true
}

func DefaultStreamEndHandler(s Stream, _ interface{}) bool {
    s.EndStream()
    return true
}
