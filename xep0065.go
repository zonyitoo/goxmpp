package xmpp

import (
    "encoding/xml"
)

const (
    XMLNS_JABBER_PROTOCOL_BYTESTREAM = "http://jabber.org/protocol/bytestreams"
)

type XMPPProtocolByteStreams struct {
    XMLName xml.Name `xml:"http://jabber.org/protocol/bytestreams query"`
    SID     string   `xml:"sid,attr"`
    Mode    string   `xml:"mode,attr"`
}

type XMPPProtocolByteStreamsHost struct {
    XMLName  xml.Name `xml:"streamhost"`
    JID      string   `xml:"jid,attr"`
    Host     string   `xml:"host,attr"`
    ZeroConf string   `xml:"zerocofn,attr"`
}
