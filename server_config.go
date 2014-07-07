package xmpp

import ()

type ServerConfig struct {
    listenAddr string
    serverName string
}

func NewDefaultServerConfig() *ServerConfig {
    return &ServerConfig{
        listenAddr: ":5222",
        serverName: "",
    }
}
