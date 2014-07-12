package xmpp

import ()

type ServerConfig struct {
    ListenAddr     string
    ServerName     string
    StreamVersion  *StreamVersion
    ServerType     int
    UseTLS         bool
    SASLMechanisms []string
    PingEnabled    bool // XEP-0199
}

const (
    SERVER_TYPE_C2S = iota
    SERVER_TYPE_S2S
)

func NewDefaultServerConfig() *ServerConfig {
    return &ServerConfig{
        ListenAddr:     ":5222",
        ServerName:     "",
        StreamVersion:  &StreamVersion{Major: 1, Minor: 0},
        ServerType:     SERVER_TYPE_C2S,
        UseTLS:         true,
        SASLMechanisms: []string{},
        PingEnabled:    true,
    }
}
