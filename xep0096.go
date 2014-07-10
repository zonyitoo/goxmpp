package xmpp

import (
    "encoding/xml"
)

type XMPPProtocolStreamInit struct {
    XMLName  xml.Name                    `xml:"http://jabber.org/protocol/si si"`
    Id       string                      `xml:"id,attr"`
    MimeType string                      `xml:"mime-type,attr"`
    Profile  string                      `xml:"profile,attr"`
    File     *XMPPProtocolStreamInitFile `xml:",omitempty"`
}

type XMPPProtocolStreamInitFile struct {
    XMLName    xml.Name                         `xml:"http://jabber.org/protocol/si/profile/file-transfer file"`
    Name       string                           `xml:"name,attr,omitempty"`
    Size       uint                             `xml:"size,attr,omitempty"`
    Hash       string                           `xml:"hash,attr,omitempty"`
    Date       string                           `xml:"date,attr,omitempty"`
    Desc       string                           `xml:"desc,omitempty"`
    FeatureNeg *XMPPProtocoltreamInitFeatureNeg `xml:",omitempty"`
    Range      *XMPPProtocolStreamInitRange     `xml:",omitempty"`
}

type XMPPProtocoltreamInitFeatureNeg struct {
    XMLName xml.Name   `xml:"http://jabber.org/protocol/feature-neg feature"`
    XData   *XMPPXData `xml:",omitempty"`
}

type XMPPProtocolStreamInitRange struct {
    XMLName xml.Name `xml:"range"`
    Length  uint     `xml:"length,attr,omitempty"`
    Offset  uint     `xml:"offset,attr,omitempty"`
}
