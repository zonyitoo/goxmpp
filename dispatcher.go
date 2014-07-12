package xmpp

import (
    "errors"
)

const (
    EVENT_STREAM_HEADER = iota
    EVENT_STREAM_TLS_NEGOCIATION
    EVENT_STREAM_SASL_NEGOCIATION
    EVENT_STREAM_FEATURE
    EVENT_STREAM_STANZA_INFO_QUERY
    EVENT_STREAM_STANZA_MESSAGE
    EVENT_STREAM_STANZA_PRESENCE
    EVENT_STREAM_ERROR
    EVENT_STREAM_END

    // Extension
    EVENT_STREAM_COMPRESSION // XEP-0138

    EVENT_INPOSSIBLE
)

var (
    EventDispatcherUnrecognizedEventError = errors.New("Unrecognized event")
    EventDispatcherIgnoredEventError      = errors.New("Ignored event")
)

type StreamEventDispatcher struct {
    handlers map[int][]XMPPEventHandler
}

func NewStreamEventDispatcher() *StreamEventDispatcher {
    return &StreamEventDispatcher{
        handlers: make(map[int][]XMPPEventHandler),
    }
}

func (ed *StreamEventDispatcher) Dispatch(s Stream, x interface{}) error {
    switch x.(type) {
    case *XMPPStream:
        return ed.dispatch(EVENT_STREAM_HEADER, s, x)
    case *XMPPStreamEnd:
        return ed.dispatch(EVENT_STREAM_END, s, x)
    case *XMPPStreamFeatures:
        return ed.dispatch(EVENT_STREAM_FEATURE, s, x)
    case *XMPPStreamError:
        return ed.dispatch(EVENT_STREAM_ERROR, s, x)
    case *XMPPStartTLS, *XMPPTLSAbort, *XMPPTLSProceed, *XMPPTLSFailure:
        return ed.dispatch(EVENT_STREAM_TLS_NEGOCIATION, s, x)
    case *XMPPSASLAuth, *XMPPSASLChallenge, *XMPPSASLResponse, *XMPPSASLAbort, *XMPPSASLFailure, *XMPPSASLSuccess:
        return ed.dispatch(EVENT_STREAM_SASL_NEGOCIATION, s, x)
    case *XMPPStanzaIQ:
        return ed.dispatch(EVENT_STREAM_STANZA_INFO_QUERY, s, x)
    case *XMPPStanzaPresence:
        return ed.dispatch(EVENT_STREAM_STANZA_PRESENCE, s, x)
    case *XMPPStanzaMessage:
        return ed.dispatch(EVENT_STREAM_STANZA_MESSAGE, s, x)

    // Extensions
    // XEP-0138
    case *XMPPStreamCompressionCompress, *XMPPStreamCompressionCompressed, *XMPPStreamCompressionFailure:
        return ed.dispatch(EVENT_STREAM_COMPRESSION, s, x)
    }

    return EventDispatcherUnrecognizedEventError
}

func (ed *StreamEventDispatcher) dispatch(ev int, s Stream, x interface{}) error {
    if hdls, ok := ed.handlers[ev]; !ok {
        return EventDispatcherIgnoredEventError
    } else {
        for _, hdl := range hdls {
            if hdl(s, x) {
                return nil
            }
        }
        return EventDispatcherIgnoredEventError
    }
}

func (ed *StreamEventDispatcher) AddHandlerForEvent(ev int, hdl XMPPEventHandler) {
    if hdl == nil {
        panic("Handler should not be nil")
    }

    if ev < 0 || ev >= EVENT_INPOSSIBLE {
        panic("Impossible event code")
    }

    if _, ok := ed.handlers[ev]; !ok {
        ed.handlers[ev] = []XMPPEventHandler{hdl}
    } else {
        ed.handlers[ev] = append([]XMPPEventHandler{hdl}, ed.handlers[ev]...)
    }
}
