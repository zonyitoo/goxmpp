package xmpp

import (
    "encoding/xml"
)

const (
    XMLNS_JABBER_IQ_REGISTER = "jabber:iq:register"
)

type XMPPStanzaIQRegisterQuery struct {
    XMLName      xml.Name                        `xml:"jabber:iq:register query"`
    Instructions string                          `xml:"instructions,omitempty"`
    Registered   *XMPPStanzaIQRegisterRegistered `xml:",omitempty"`
    Remove       *XMPPStanzaIQRegisterRemove     `xml:",omitempty"`

    Username string `xml:"username,omitempty"`
    Password string `xml:"password,omitempty"`
    Nick     string `xml:"nick,omitempty"`
    Name     string `xml:"name,omitempty"`
    First    string `xml:"first,omitempty"`
    Last     string `xml:"last,omitempty"`
    EMail    string `xml:"email,omitempty"`
    Address  string `xml:"address,omitempty"`
    City     string `xml:"city,omitempty"`
    State    string `xml:"state,omitempty"`
    Zip      string `xml:"zip,omitempty"`
    Phone    string `xml:"phone,omitempty"`
    URL      string `xml:"url,omitempty"`
    Date     string `xml:"date,omitempty"`
    Misc     string `xml:"misc,omitempty"`
    Text     string `xml:"text,omitempty"`
    Key      string `xml:"key,omitempty"`
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
