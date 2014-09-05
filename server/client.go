package server

import (
    "github.com/zonyitoo/goxmpp/stream"
    "net"
)

type Client interface {
    Stream() stream.Streamer
    Run()
}

type TCPClient struct {
    stream *stream.ServerClientStream
}

func NewTCPClient(conn net.Conn, a *stream.SASLAuthenticator, shandler stream.StanzaHandler) *TCPClient {
    return &TCPClient{
        stream: stream.NewServerClientStream(conn, a, shandler),
    }
}

func (c *TCPClient) Stream() stream.Streamer {
    return c.stream
}

func (c *TCPClient) Run() {
    c.stream.Run()
}
