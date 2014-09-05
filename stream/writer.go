package stream

import (
    "encoding/xml"
    "github.com/zonyitoo/goxmpp/protocol"
    "io"
    "sync"
)

type Writer struct {
    transport io.Writer
    wchan     chan []byte
    wgroup    sync.WaitGroup
}

func NewWriter(transport io.Writer) *Writer {
    sw := &Writer{
        transport: transport,
        wchan:     make(chan []byte),
    }
    sw.wgroup.Add(1)
    go sw.send()
    return sw
}

func (sw *Writer) send() {
    defer sw.wgroup.Done()
    for data := range sw.wchan {
        if _, err := sw.transport.Write(data); err != nil {
            break
        }
    }
}

func (sw *Writer) Write(data []byte) (int, error) {
    return len(data), sw.SendBytes(data)
}

func (sw *Writer) SendBytes(data []byte) error {
    sw.wchan <- data
    return nil
}

func (sw *Writer) Close() error {
    if err := sw.SendBytes([]byte(protocol.XMPPStreamEndFmt)); err != nil {
        return err
    }
    return sw.Destroy()
}

func (sw *Writer) Destroy() error {
    close(sw.wchan)
    sw.wgroup.Wait()
    return nil
}

func (sw *Writer) Open(stream *protocol.XMPPStream) error {
    sw.SendBytes([]byte(xml.Header))
    header := protocol.GenXMPPStreamHeader(stream)
    if err := sw.SendBytes([]byte(header)); err != nil {
        return err
    }
    return nil
}

func (sw *Writer) SendElement(elem protocol.Protocol) error {
    data, err := xml.Marshal(elem)
    if err != nil {
        return err
    }
    return sw.SendBytes(data)
}
