package stream

import (
    "net"
)

type Streamer interface {
    Id() string
    RemoteAddr() net.Addr
    Writer() *Writer
    Reader() *Reader
    Close() error
}
