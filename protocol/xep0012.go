package protocol

import (
    "encoding/xml"
)

const XMLNS_JABBER_IQ_LAST = "jabber:iq:last"

type XMPPStanzaIQLastActivityQuery struct {
    XMLName xml.Name `xml:"jabber:iq:last query"`
    Seconds uint     `xml:"seconds,attr,omitempty"`
    Status  string   `xml:",chardata"`
}
