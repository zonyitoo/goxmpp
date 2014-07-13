package xmpp

import ()

type Config struct {
    ListenAddr     string
    ServerAddr     string
    ServerName     string
    StreamVersion  *StreamVersion
    ServerType     int
    UseTLS         bool
    SASLMechanisms []string
    PingEnabled    bool // XEP-0199
    ClientJID      *JID
    ClientPassword string
}

const (
    SERVER_TYPE_C2S = iota
    SERVER_TYPE_S2S
)

func NewDefaultServerConfig() *Config {
    return &Config{
        ListenAddr:     ":5222",
        ServerName:     "",
        StreamVersion:  &StreamVersion{Major: 1, Minor: 0},
        ServerType:     SERVER_TYPE_C2S,
        UseTLS:         true,
        SASLMechanisms: []string{},
        PingEnabled:    true,
    }
}

func NewDefaultClientConfig(addr, jid, password string) *Config {
    jid_struct, err := NewJIDFromString(jid)
    if err != nil {
        panic(err)
    }
    return &Config{
        ServerAddr:     addr,
        StreamVersion:  &StreamVersion{Major: 1, Minor: 0},
        UseTLS:         true,
        SASLMechanisms: []string{},
        PingEnabled:    true,
        ClientJID:      jid_struct,
        ClientPassword: password,
    }
}
