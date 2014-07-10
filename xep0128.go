package xmpp

import (
    "encoding/xml"
)

const (
    XMLNS_JABBER_X_DATA = "jabber:x:data"
)

type XMPPXData struct {
    XMLName xml.Name          `xml:"jabber:x:data x"`
    Type    string            `xml:"type,attr,omitempty"`
    Fields  []*XMPPXDataField `xml:",omitempty"`
}

type XMPPXDataField struct {
    XMLName        xml.Name                       `xml:"field"`
    Var            string                         `xml:"var,attr"`
    Label          string                         `xml:"label,attr,omitempty"`
    Values         []string                       `xml:"value,omitempty"`
    OptionalValues []*XMPPXDataFieldOptionalValue `xml:",omitempty"`
}

type XMPPXDataFieldOptionalValue struct {
    XMLName xml.Name `xml:"option"`
    Value   string   `xml:"value"`
}
