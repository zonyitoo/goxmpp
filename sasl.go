package xmpp

import (
    "errors"
    "strings"
)

type XMPPSASLAuthenticator interface {
    Auth(stream *Stream, data string)
    Challenge(stream *Stream, data string)
    Response(stream *Stream, data string)
    Abort(stream *Stream)
    Failure(stream *Stream, err interface{})
    Success(stream *Stream, data string)
}

type XMPPSASLAuthenticatorFactory func() XMPPSASLAuthenticator

type XMPPSASLManager struct {
    mechanisms map[string]XMPPSASLAuthenticatorFactory
}

func NewXMPPSASLManager() *XMPPSASLManager {
    return &XMPPSASLManager{
        mechanisms: make(map[string]XMPPSASLAuthenticatorFactory),
    }
}

func (m *XMPPSASLManager) SetMechanism(mech string, authenticator XMPPSASLAuthenticatorFactory) {
    upper := strings.ToUpper(mech)
    m.mechanisms[upper] = authenticator
}

func (m *XMPPSASLManager) GetAuthenticatorFactory(mech string) (XMPPSASLAuthenticatorFactory, error) {
    upper := strings.ToUpper(mech)
    if auth, ok := m.mechanisms[upper]; !ok {
        return nil, errors.New("Unsupported mechanism")
    } else {
        return auth, nil
    }
}
