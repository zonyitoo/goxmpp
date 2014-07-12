package xmpp

import (
    "encoding/xml"
)

const XMLNS_PROTOCOL_STREAM_INITIATION = "http://jabber.org/protocol/si"

type XMPPProtocolStreamInitiation struct {
    XMLName            xml.Name                        `xml:"http://jabber.org/protocol/si si"`
    Id                 string                          `xml:"id,attr,omitempty"`
    MimeType           string                          `xml:"mime-type,attr,omitempty"`
    Profile            string                          `xml:"profile,attr,omitempty"`
    FeatureNegociation *XMPPProtocolFeatureNegociation `xml:",omitempty"`

    FileTransfer *XMPPProtocolStreamInitiationProfileFileTransfer `xml:",omitempty"` // XEP-0096
}
