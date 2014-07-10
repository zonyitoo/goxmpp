package xmpp

import (
    "encoding/xml"
)

type XMPPStanzaIQDiscoInfo struct {
    XMLName   xml.Name                         `xml:"http://jabber.org/protocol/disco#info query"`
    Identitiy []*XMPPStanzaIQDiscoInfoIdentity `xml:",omitempty"`
    Features  []*XMPPStanzaIQDiscoInfoFeature  `xml:",omitempty"`
    Node      string                           `xml:"node,attr,omitempty"`

    // XEP-0128
    XData *XMPPXData `xml:",omitempty"`
}

type XMPPStanzaIQDiscoInfoIdentity struct {
    XMLName  xml.Name `xml:"identity"`
    Category string   `xml:"category,attr"`
    Type     string   `xml:"type,attr"`
    Name     string   `xml:"type,attr"`
}

type XMPPStanzaIQDiscoInfoFeature struct {
    XMLName xml.Name `xml:"feature"`
    Var     string   `xml:"var,attr"`
}

type XMPPStanzaIQDiscoItems struct {
    XMLName xml.Name                      `xml:"http://jabber.org/protocol/disco#item query"`
    Items   []*XMPPStanzaIQDiscoItemsItem `xml:",omitempty"`
}

type XMPPStanzaIQDiscoItemsItem struct {
    XMLName xml.Name `xml:"item"`
    JID     string   `xml:"jid,attr"`
    Name    string   `xml:"name,attr,omitempty"`
    Node    string   `xml:"node,attr,omitempty"`
}
