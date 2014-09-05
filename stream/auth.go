package stream

import (
    "github.com/zonyitoo/goxmpp/protocol"
)

type SASLAuthenticateHandler func(*protocol.XMPPSASLAuth, Streamer) bool

type SASLAuthenticator struct {
    handlers map[string]SASLAuthenticateHandler
}

func NewSASLAuthenticator() *SASLAuthenticator {
    return &SASLAuthenticator{
        handlers: make(map[string]SASLAuthenticateHandler),
    }
}

func (a *SASLAuthenticator) SetMechanism(name string, handler SASLAuthenticateHandler) {
    a.handlers[name] = handler
}

func (a *SASLAuthenticator) CallMechanism(name string, auth *protocol.XMPPSASLAuth, s Streamer) bool {
    if handler, ok := a.handlers[name]; !ok {
        return false
    } else {
        return handler(auth, s)
    }
}

func (a *SASLAuthenticator) Mechanisms() []string {
    mech := make([]string, len(a.handlers))
    idx := 0
    for k, _ := range a.handlers {
        mech[idx] = k
        idx++
    }
    return mech
}
