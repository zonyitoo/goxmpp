package xmpp

import ()

type ServerConfig struct {
    ListenAddr           string
    ServerName           string
    StreamHandlerFactory func() StreamHandler
    StreamVersion        *StreamVersion
}

func NewDefaultServerConfig() *ServerConfig {
    return &ServerConfig{
        ListenAddr:           ":5222",
        ServerName:           "",
        StreamHandlerFactory: DefaultStreamHandlerFactory,
        StreamVersion:        &StreamVersion{Major: 1, Minor: 0},
    }
}
