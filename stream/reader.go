package stream

import (
    "github.com/zonyitoo/goxmpp/protocol"
    "io"
)

type Reader struct {
    decoder *Decoder
}

func NewReader(r io.Reader) *Reader {
    return &Reader{
        decoder: NewDecoder(r),
    }
}

func (sr *Reader) NextElement() (protocol.Protocol, error) {
    return sr.decoder.GetNextElement()
}
