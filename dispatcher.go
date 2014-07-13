package xmpp

import (
    "errors"
)

const (
    EVENT_STREAM_INIT = iota
    EVENT_STREAM_HEADER
    EVENT_STREAM_TLS_NEGOCIATION
    EVENT_STREAM_TLS_ABORTED
    EVENT_STREAM_TLS_FAILURE
    EVENT_STREAM_TLS_PROCEEDED
    EVENT_STREAM_SASL_NEGOCIATION
    EVENT_STREAM_SASL_SUCCESSED
    EVENT_STREAM_SASL_FAILURE
    EVENT_STREAM_SASL_ABORTED
    EVENT_STREAM_FEATURE
    EVENT_STREAM_STANZA_INFO_QUERY
    EVENT_STREAM_STANZA_MESSAGE
    EVENT_STREAM_STANZA_PRESENCE
    EVENT_STREAM_ERROR
    EVENT_STREAM_END

    // Extension
    EVENT_STREAM_COMPRESSION // XEP-0138

    EVENT_IMPOSSIBLE
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

func (ed *StreamEventDispatcher) EventOf(x interface{}) int {
    switch x.(type) {
    case *XMPPStream:
        return EVENT_STREAM_HEADER
    case *XMPPStreamEnd:
        return EVENT_STREAM_END
    case *XMPPStreamFeatures:
        return EVENT_STREAM_FEATURE
    case *XMPPStreamError:
        return EVENT_STREAM_ERROR
    case *XMPPStartTLS:
        return EVENT_STREAM_TLS_NEGOCIATION
    case *XMPPTLSAbort:
        return EVENT_STREAM_TLS_ABORTED
    case *XMPPTLSProceed:
        return EVENT_STREAM_TLS_PROCEEDED
    case *XMPPTLSFailure:
        return EVENT_STREAM_TLS_FAILURE
    case *XMPPSASLAuth, *XMPPSASLChallenge, *XMPPSASLResponse:
        return EVENT_STREAM_SASL_NEGOCIATION
    case *XMPPSASLAbort:
        return EVENT_STREAM_SASL_ABORTED
    case *XMPPSASLFailure:
        return EVENT_STREAM_SASL_FAILURE
    case *XMPPSASLSuccess:
        return EVENT_STREAM_SASL_SUCCESSED
    case *XMPPStanzaIQ:
        return EVENT_STREAM_STANZA_INFO_QUERY
    case *XMPPStanzaPresence:
        return EVENT_STREAM_STANZA_PRESENCE
    case *XMPPStanzaMessage:
        return EVENT_STREAM_STANZA_MESSAGE

    // Extensions
    // XEP-0138
    case *XMPPStreamCompressionCompress, *XMPPStreamCompressionCompressed, *XMPPStreamCompressionFailure:
        return EVENT_STREAM_COMPRESSION

    default:
        return EVENT_IMPOSSIBLE
    }
}

func (ed *StreamEventDispatcher) Dispatch(s Stream, x interface{}) error {
    ev := ed.EventOf(x)
    if ev != EVENT_IMPOSSIBLE {
        return ed.DispatchEvent(ev, s, x)
    }

    return EventDispatcherUnrecognizedEventError
}

func (ed *StreamEventDispatcher) DispatchEvent(ev int, s Stream, x interface{}) error {
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

    if ev < 0 || ev >= EVENT_IMPOSSIBLE {
        panic("Impossible event code")
    }

    if _, ok := ed.handlers[ev]; !ok {
        ed.handlers[ev] = []XMPPEventHandler{hdl}
    } else {
        ed.handlers[ev] = append([]XMPPEventHandler{hdl}, ed.handlers[ev]...)
    }
}
