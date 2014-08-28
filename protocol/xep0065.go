package protocol

import (
    "encoding/xml"
)

const (
    XMLNS_JABBER_PROTOCOL_BYTESTREAM = "http://jabber.org/protocol/bytestreams"
)

type XMPPProtocolByteStreamQuery struct {
    XMLName  xml.Name                        `xml:"http://jabber.org/protocol/bytestreams query"`
    SID      string                          `xml:"sid,attr"`
    Mode     string                          `xml:"mode,attr,omitempty"`
    DstAddr  string                          `xml:"dstaddr,attr,omitempty"`
    Host     *XMPPProtocolByteStreamHost     `xml:",omitempty"`
    HostUsed *XMPPProtocolByteStreamHostUsed `xml:",omitempty"`
    Activate string                          `xml:"activate,omitempty"`
}

const (
    XMPP_PROTOCOL_BYTE_STREAM_MODE_TCP = "tcp"
    XMPP_PROTOCOL_BYTE_STREAM_MODE_UDP = "udp"
)

type XMPPProtocolByteStreamHost struct {
    XMLName xml.Name `xml:"streamhost"`
    JID     string   `xml:"jid,attr"`
    Host    string   `xml:"host,attr"`
    Port    uint     `xml:"port,attr,omitempty"`
}

type XMPPProtocolByteStreamHostUsed struct {
    XMLName xml.Name `xml:"streamhost-used"`
    JID     string   `xml:"jid,attr"`
}

type XMPPProtocolByteStreamUDPSuccess struct {
    XMLName xml.Name `xml:"udpsuccess"`
    DstAddr string   `xml:"dstaddr,attr"`
}
