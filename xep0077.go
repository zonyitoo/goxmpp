package xmpp

import (
    "encoding/xml"
)

const (
    XMLNS_JABBER_IQ_REGISTER = "jabber:iq:register"
)

type XMPPStanzaIQRegister struct {
    XMLName      xml.Name                        `xml:"jabber:iq:register query"`
    Instructions string                          `xml:"instructions,omitempty"`
    Registered   *XMPPStanzaIQRegisterRegistered `xml:",omitempty"`
    Remove       *XMPPStanzaIQRegisterRemove     `xml:",omitempty"`

    Username string `xml:"username,omitempty"`
    Password string `xml:"password,omitempty"`
    EMail    string `xml:"email,omitempty"`
    // ... omited fields ...
}

type XMPPStanzaIQRegisterRegistered struct {
    XMLName xml.Name `xml:"registered"`
}

type XMPPStanzaIQRegisterRemove struct {
    XMLName xml.Name `xml:"remove"`
}

type XMPPStreamFeatureRegister struct {
    XMLName xml.Name `xml:"http://jabber.org/features/iq-register register"`
}
