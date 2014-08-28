package protocol

import (
    "encoding/xml"
)

const XMLNS_PROTOCOL_FEATURE_NEGOCIATION = "http://jabber.org/protocol/feature-neg"

type XMPPProtocolFeatureNegociation struct {
    XMLName xml.Name   `xml:"http://jabber.org/protocol/feature-neg feature"`
    XData   *XMPPXData `xml:""`
}
