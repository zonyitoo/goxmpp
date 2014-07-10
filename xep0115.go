package xmpp

import (
    "encoding/xml"
)

type XMPPStanzaPresenceCaps struct {
    XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
    Hash    string   `xml:"hash,attr"`
    Node    string   `xml:"node,attr"`
    Ver     string   `xml:"ver,attr"`
}
