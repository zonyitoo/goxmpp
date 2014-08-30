package stream

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "github.com/zonyitoo/goxmpp/protocol"
    // "io"
    "testing"
)

func Test_Writer(t *testing.T) {
    buf := &bytes.Buffer{}

    sw := NewWriter(buf)

    sheader := &protocol.XMPPStream{
        From:    "juliet@example.com",
        To:      "example.com",
        Id:      "abcd",
        Version: "1.0",
        XMLLang: "en",
        Xmlns:   protocol.XMLNS_JABBER_CLIENT,
    }
    assert.NoError(t, sw.Open(sheader))

    features := &protocol.XMPPStreamFeatures{
        StartTLS: &protocol.XMPPStartTLS{},
    }
    assert.NoError(t, sw.SendElement(features))

    assert.NoError(t, sw.Close())
    assert.Equal(t, `<stream:stream from='juliet@example.com' to='example.com' version='1.0' xml:lang='en' id='abcd' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams'><features xmlns="http://etherx.jabber.org/streams"><starttls xmlns="urn:ietf:params:xml:ns:xmpp-tls"></starttls></features></stream:stream>`, buf.String())
}
