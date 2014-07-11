package xmpp

import (
    "encoding/xml"
)

type XMPPStanzaIQPing struct {
    XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}
