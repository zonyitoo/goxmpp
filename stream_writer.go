package xmpp

import (
    "encoding/xml"
    "errors"
    "github.com/zonyitoo/goxmpp/protocol"
    "io"
    "sync"
)

type StreamWriter struct {
    transport  io.Writer
    wchan      chan []byte
    wgroup     sync.WaitGroup
    hasStarted bool
}

func NewStreamWriter(transport io.Writer) *StreamWriter {
    return &StreamWriter{
        transport:  transport,
        wchan:      nil,
        hasStarted: false,
    }
}

func (sw *StreamWriter) send() {
    defer sw.wgroup.Done()
    for data := range sw.wchan {
        if _, err := sw.transport.Write(data); err != nil {
            break
        }
    }
}

func (sw *StreamWriter) Write(data []byte) (int, error) {
    return len(data), sw.SendBytes(data)
}

func (sw *StreamWriter) SendBytes(data []byte) error {
    if !sw.hasStarted {
        return errors.New("Stream writer is closed")
    }
    sw.wchan <- data
    return nil
}

func (sw *StreamWriter) Close() error {
    if !sw.hasStarted {
        return errors.New("Stream writer is already closed")
    }
    if err := sw.SendBytes([]byte("</stream:stream>")); err != nil {
        return err
    }
    close(sw.wchan)
    sw.wgroup.Wait()
    sw.hasStarted = false
    return nil
}

func (sw *StreamWriter) Destroy() error {
    if !sw.hasStarted {
        return errors.New("Stream writer is already destroyed")
    }
    close(sw.wchan)
    sw.wgroup.Wait()
    sw.hasStarted = false
    return nil
}

func (sw *StreamWriter) Open(stream *protocol.XMPPStream) error {
    if sw.hasStarted {
        return errors.New("Stream writer is already opened")
    }
    header := protocol.GenXMPPStreamHeader(stream)
    sw.hasStarted = true
    sw.wchan = make(chan []byte, 1024)
    sw.wgroup.Add(1)
    go sw.send()

    if err := sw.SendBytes([]byte(header)); err != nil {
        return err
    }
    return nil
}

func (sw *StreamWriter) SendElement(elem protocol.Protocol) error {
    data, err := xml.Marshal(elem)
    if err != nil {
        return err
    }
    return sw.SendBytes(data)
}
