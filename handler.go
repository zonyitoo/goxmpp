package xmpp

import (
    log "github.com/cihub/seelog"
)

type XMPPEventHandler func(Stream, interface{}) bool

func DefaultStreamHeaderHandler(s Stream, _x interface{}) bool {
    if s.State()&STREAM_STAT_STARTED != 0 {
        err := XMPPStreamError{
            InvalidXML: &XMPPStreamErrorInvalidXML{},
        }
        log.Errorf("JID:`%s` stream is already started", s.JID())
        s.SendErrorAndClose(err)
        return true
    }

    x, _ := _x.(*XMPPStream)

    if x.Xmlns == XMLNS_JABBER_CLIENT {
        if x.From != "" {
            if from_jid, err := NewJIDFromString(x.From); err != nil {
                goto CLIENT_FROM_ERROR
            } else if from_jid.Domain != x.To {
                goto CLIENT_FROM_ERROR
            } else {
                s.SetJID(from_jid)
            }

            goto NO_FROM_ERROR
        CLIENT_FROM_ERROR:
            err := XMPPStreamError{
                InvalidFrom: &XMPPStreamErrorInvalidFrom{},
            }
            s.SendErrorAndClose(err)
            return true
        }
    } else {
        err := XMPPStreamError{
            InvalidNamespace: &XMPPStreamErrorInvalidNamespace{},
        }
        s.SendErrorAndClose(err)
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
        s.SendErrorAndClose(err)
        return true
    }

    s.StartStream(STREAM_TYPE_CLIENT, x.To, x.From, s.ServerConfig().StreamVersion.String(), x.XMLLang)
    s.SetState(s.State() | STREAM_STAT_STARTED)
    return true
}

func DefaultStreamEndHandler(s Stream, _ interface{}) bool {
    s.EndStream()
    return true
}
