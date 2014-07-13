package xmpp

import (
    log "github.com/cihub/seelog"
)

type XMPPEventHandler func(Stream, interface{}) bool

func DefaultC2SStreamHeaderHandler(s Stream, _x interface{}) bool {
    if s.State()&STREAM_STAT_STARTED != 0 {
        err := XMPPStreamError{
            InvalidXML: &XMPPStreamErrorInvalidXML{},
        }
        log.Errorf("JID:`%s` stream is already started", s.JID())
        s.SendErrorAndClose(err)
        return true
    }

    x, _ := _x.(*XMPPStream)

    if x.Xmlns != XMLNS_JABBER_CLIENT {
        err := XMPPStreamError{
            InvalidNamespace: &XMPPStreamErrorInvalidNamespace{},
        }
        s.SendErrorAndClose(err)
        return true
    } else {
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
    NO_FROM_ERROR:
    }

    return defaultStreamHeaderHandler(s, x)
}

func DefaultS2SStreamHeaderHandler(s Stream, _x interface{}) bool {
    if s.State()&STREAM_STAT_STARTED != 0 {
        err := XMPPStreamError{
            InvalidXML: &XMPPStreamErrorInvalidXML{},
        }
        log.Errorf("JID:`%s` stream is already started", s.JID())
        s.SendErrorAndClose(err)
        return true
    }

    x, _ := _x.(*XMPPStream)

    if x.Xmlns != XMLNS_JABBER_SERVER {
        err := XMPPStreamError{
            InvalidNamespace: &XMPPStreamErrorInvalidNamespace{},
        }
        s.SendErrorAndClose(err)
        return true
    } else {
        if x.From != "" {
            if from_jid, err := NewJIDFromString(x.From); err != nil {
                err := XMPPStreamError{
                    InvalidFrom: &XMPPStreamErrorInvalidFrom{},
                }
                s.SendErrorAndClose(err)
                return true
            } else {
                s.SetJID(from_jid)
            }
        }
    }

    return defaultStreamHeaderHandler(s, x)
}

func defaultStreamHeaderHandler(s Stream, x *XMPPStream) bool {
    if x.Version == "" {
        x.Version = "0.9"
    }

    version := &StreamVersion{}
    version.FromString(x.Version)

    if !s.Config().StreamVersion.GreaterOrEqualTo(version) {
        err := XMPPStreamError{
            UnsupportedVersion: &XMPPStreamErrorUnsupportedVersion{},
        }
        s.SendErrorAndClose(err)
        return true
    }

    s.StartStream(STREAM_TYPE_CLIENT, x.To, x.From, s.Config().StreamVersion.String(), x.XMLLang)
    s.SetState(s.State() | STREAM_STAT_STARTED)

    feature := &XMPPStreamFeatures{}
    if s.State()&STREAM_STAT_TLS_PROCEED == 0 && s.Config().UseTLS {
        feature.StartTLS = &XMPPStartTLS{
            Required: &XMPPRequired{},
        }
    } else if s.State()&STREAM_STAT_SASL_SUCCEEDED == 0 {
        feature.SASLMechanisms = &XMPPSASLMechanisms{
            Mechanisms: s.Config().SASLMechanisms,
        }
    } else {
        feature.Bind = &XMPPBind{}
    }
    s.SendElement(feature)
    return true
}

func DefaultStreamEndHandler(s Stream, _ interface{}) bool {
    s.EndStream()
    return true
}

func DefaultPingServerHandler(s Stream, _x interface{}) bool {
    x, _ := _x.(*XMPPStanzaIQ)

    if x.Ping != nil {
        if x.To == s.Config().ServerName {
            iq := XMPPStanzaIQ{
                Type: XMPP_STANZA_IQ_TYPE_RESULT,
                Id:   x.Id,
                From: x.To,
                To:   x.From,
                Ping: &XMPPStanzaIQPing{},
            }

            if !s.Config().PingEnabled {
                iq.Error = &XMPPStanzaError{
                    Type: XMPP_STANZA_ERROR_TYPE_CANCEL,
                    XMPPStanzaErrorGroup: XMPPStanzaErrorGroup{
                        ServiceUnavailable: &XMPPStanzaErrorServiceUnavailable{},
                    },
                }
                s.SendElement(iq)
                return true
            }

            switch x.Type {
            case XMPP_STANZA_IQ_TYPE_GET:
                // Just return a result
            case XMPP_STANZA_IQ_TYPE_RESULT:
                // TODO
            case XMPP_STANZA_IQ_TYPE_ERROR:
                // TODO
            default:
                iq.Error = &XMPPStanzaError{
                    Type: XMPP_STANZA_ERROR_TYPE_CANCEL,
                }
            }
            s.SendElement(iq)
        } else {
            // Dispatch to the destination
        }
        return true
    } else {
        return false
    }
}
