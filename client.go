package xmpp

import (
    "net"
)

type Client struct {
    streamDispatcher *StreamEventDispatcher
    config           *Config
    stream           Stream
}

func NewClient(conf *Config) *Client {
    conn, err := net.Dial("tcp", conf.ServerAddr)
    if err != nil {
        panic(err)
    }

    client := &Client{
        streamDispatcher: NewStreamEventDispatcher(),
        config:           conf,
    }
    client.stream = NewClientStream(NewTCPTransport(conn), client)
    return client
}

func (c *Client) AddHandlerForEvent(ev int, hdl XMPPEventHandler) {
    c.streamDispatcher.AddHandlerForEvent(ev, hdl)
}

func (c *Client) Run() {
    c.stream.Run()
}
