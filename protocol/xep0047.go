package protocol

import (
    "encoding/xml"
)

const XMLNS_JABBER_PROTOCOL_INBAND_BYTESTREAM = "http://jabber.org/protocol/ibb"

type XMPPProtocolInBandByteStreamOpen struct {
    XMLName   xml.Name `xml:"http://jabber.org/protocol/ibb open"`
    BlockSize uint     `xml:"block-size,attr"`
    SID       string   `xml:"sid,attr"`
    Stanza    string   `xml:"stanza,attr,omitempty"`
}

type XMPPProtocolInBandByteStreamData struct {
    XMLName xml.Name `xml:"http://jabber.org/protocol/ibb data"`
    Seq     uint     `xml:"seq,attr"`
    SID     string   `xml:"sid,attr"`
    Data    string   `xml:",chardata"`
}

type XMPPProtocolInBandByteStreamClose struct {
    XMLName xml.Name `xml:"http://jabber.org/protocol/ibb close"`
    SID     string   `xml:"sid,attr"`
}
