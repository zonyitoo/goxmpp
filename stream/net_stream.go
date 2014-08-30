package stream

import (
    "code.google.com/p/go-uuid/uuid"
    "net"
)

type NetStream struct {
    conn   net.Conn
    reader *Reader
    writer *Writer
    id     string
}

func NewNetStream(conn net.Conn) *NetStream {
    return &NetStream{
        conn:   conn,
        reader: NewReader(conn),
        writer: NewWriter(conn),
        id:     uuid.New(),
    }
}

func (ns *NetStream) RemoteAddr() net.Addr {
    return ns.conn.RemoteAddr()
}

func (ns *NetStream) Writer() *Writer {
    return ns.writer
}

func (ns *NetStream) Reader() *Reader {
    return ns.reader
}

func (ns *NetStream) Id() string {
    return ns.id
}
