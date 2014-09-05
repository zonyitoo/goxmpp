package auth

import (
    "errors"
    "github.com/zonyitoo/goxmpp/stream"
)

type SASLAuthenticate struct{}

type SASLAuthenticator interface {
    Process(*stream.Reader, *stream.Writer, string) bool
}

type SASLAuthencatorFactory func() SASLAuthenticator

type SASLAuthenticateManager struct {
    mechanisms map[string]SASLAuthencatorFactory
}

func NewAuthenticateManager() *SASLAuthenticateManager {
    return &SASLAuthenticateManager{
        mechanisms: make(map[string]SASLAuthencatorFactory),
    }
}

func (am *SASLAuthenticateManager) GetAuthenticator(mechanism string) (SASLAuthenticator, error) {
    if auth, ok := am.mechanisms[mechanism]; !ok {
        return nil, errors.New("Mechanism does not exist")
    } else {
        return auth(), nil
    }
}

func (am *SASLAuthenticateManager) SetAuthenticator(mechanism string, a SASLAuthencatorFactory) {
    am.mechanisms[mechanism] = a
}
