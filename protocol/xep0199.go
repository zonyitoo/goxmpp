package protocol

import (
    "encoding/xml"
)

const XMLNS_XMPP_PING = "urn:xmpp:ping"

type XMPPStanzaIQPing struct {
    XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}
