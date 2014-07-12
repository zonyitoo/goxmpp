package xmpp

import ()

type ServerConfig struct {
    ListenAddr    string
    ServerName    string
    StreamVersion *StreamVersion
}

func NewDefaultServerConfig() *ServerConfig {
    return &ServerConfig{
        ListenAddr:    ":5222",
        ServerName:    "",
        StreamVersion: &StreamVersion{Major: 1, Minor: 0},
    }
}
