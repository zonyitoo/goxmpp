package xmpp

import (
    "encoding/xml"
)

const XMLNS_PROTOCOL_ENTITY_CAPABILITIES = "http://jabber.org/protocol/caps"

type XMPPProtocolEntityCapabilities struct {
    XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
    Hash    string   `xml:"hash,attr"`
    Node    string   `xml:"node,attr"`
    Ver     string   `xml:"ver,attr"`
    Ext     string   `xml:"ext,attr,omitempty"`
}
