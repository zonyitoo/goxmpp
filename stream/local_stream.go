package stream

import (
    "code.google.com/p/go-uuid/uuid"
    "io"
    "net"
)

type emptyAddr struct{}

func (emptyAddr) Network() string {
    return ""
}

func (emptyAddr) String() string {
    return ""
}

type LocalStream struct {
    reader *Reader
    writer *Writer
    id     string
}

func NewLocalStream(r io.Reader, w io.Writer) *LocalStream {
    return &LocalStream{
        reader: NewReader(r),
        writer: NewWriter(w),
        id:     uuid.New(),
    }
}

func (ns *LocalStream) RemoteAddr() net.Addr {
    return emptyAddr{}
}

func (ns *LocalStream) Writer() *Writer {
    return ns.writer
}

func (ns *LocalStream) Reader() *Reader {
    return ns.reader
}

func (ns *LocalStream) Id() string {
    return ns.id
}
