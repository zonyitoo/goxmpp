package protocol

import (
    "encoding/xml"
)

const (
    XMLNS_JABBER_X_OOB  = "jabber:x:oob"
    XMLNS_JABBER_IQ_OOB = "jabber:iq:oob"
)

type XMPPXOutOfBandData struct {
    XMLName xml.Name `xml:"jabber:x:oob x"`
    URL     string   `xml:"url"`
    Desc    string   `xml:"desc,omitempty"`
}

type XMPPOutofBandDataQuery struct {
    XMLName xml.Name `xml:"jabber:iq:oob query"`
    URL     string   `xml:"url"`
    Desc    string   `xml:"desc,omitempty"`
    SID     string   `xml:"sid,attr,omitempty"`
}
