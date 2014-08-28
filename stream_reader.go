package xmpp

import (
    "github.com/zonyitoo/goxmpp/protocol"
    "io"
)

type StreamReader struct {
    decoder *Decoder
}

func NewStreamReader(r io.Reader) *StreamReader {
    return &StreamReader{
        decoder: NewDecoder(r),
    }
}

func (sr *StreamReader) NextElement() (protocol.Protocol, error) {
    return sr.decoder.GetNextElement()
}
