package xmpp

import (
    "encoding/xml"
)

type XMPPProtocolDiscoInfoQuery struct {
    XMLName    xml.Name                         `xml:"http://jabber.org/protocol/disco#info query"`
    Identities []*XMPPProtocolDiscoInfoIdentity `xml:",omitempty"`
    Features   []*XMPPProtocolDiscoInfoFeature  `xml:",omitempty"`
    Node       string                           `xml:"node,attr,omitempty"`

    // XEP-0128
    XData *XMPPXData `xml:",omitempty"`
}

type XMPPProtocolDiscoInfoIdentity struct {
    XMLName  xml.Name `xml:"identity"`
    Category string   `xml:"category,attr"`
    Type     string   `xml:"type,attr"`
    Name     string   `xml:"type,attr,omitempty"`
}

type XMPPProtocolDiscoInfoFeature struct {
    XMLName xml.Name `xml:"feature"`
    Var     string   `xml:"var,attr"`
}

type XMPPProtocolDiscoItemQuery struct {
    XMLName xml.Name                      `xml:"http://jabber.org/protocol/disco#item query"`
    Items   []*XMPPProtocolDiscoItemsItem `xml:",omitempty"`
    Node    string                        `xml:"node,attr,omitempty"`
}

type XMPPProtocolDiscoItemsItem struct {
    XMLName xml.Name `xml:"item"`
    JID     string   `xml:"jid,attr"`
    Name    string   `xml:"name,attr,omitempty"`
    Node    string   `xml:"node,attr,omitempty"`
}
