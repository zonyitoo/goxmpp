package xmpp

import (
    "encoding/xml"
)

type XMPPStreamFeatureCompression struct {
    XMLName xml.Name `xml:"http://jabber.org/features/compress compression"`
    Methods []string `xml:"method,omitempty"`
}

type XMPPStreamCreateCompress struct {
    XMLName xml.Name `xml:"http://jabber.org/features/compress ccompress"`
    Methods []string `xml:"method,omitempty"`
}

type XMPPStreamCompressed struct {
    XMLName xml.Name `xml:"http://jabber.org/features/compress compressed"`
}

type XMPPStreamCompressionFailure struct {
    XMLName xml.Name `xml:"http://jabber.org/features/compress failure"`
}

type XMPPStreamCompressionFailureUnsupportedMethod struct {
    XMLName xml.Name `xml:"unsupported-method"`
}

type XMPPStreamCompressionFailureSetupFailed struct {
    XMLName xml.Name `xml:"setup-failed"`
}
