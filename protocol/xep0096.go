package protocol

import (
    "encoding/xml"
)

const XMLNS_PROTOCOL_STREAM_INITIATION_PROFILE_FILE_TRANSFER = "http://jabber.org/protocol/si/profile/file-transfer"

type XMPPProtocolStreamInitiationProfileFileTransfer struct {
    XMLName xml.Name                           `xml:"http://jabber.org/protocol/si/profile/file-transfer file"`
    Name    string                             `xml:"name,attr,omitempty"`
    Size    uint                               `xml:"size,attr,omitempty"`
    Hash    string                             `xml:"hash,attr,omitempty"`
    Date    string                             `xml:"date,attr,omitempty"`
    Desc    string                             `xml:"desc,omitempty"`
    Range   *XMPPProtocolStreamInitiationRange `xml:",omitempty"`
}

type XMPPProtocolStreamInitiationRange struct {
    XMLName xml.Name `xml:"range"`
    Length  uint     `xml:"length,attr,omitempty"`
    Offset  uint     `xml:"offset,attr,omitempty"`
}
