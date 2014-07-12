package xmpp

import (
    "encoding/xml"
)

const XMLNS_STREAM_FEATURE_COMPRESSION = "http://jabber.org/features/compress"

type XMPPStreamFeatureCompression struct {
    XMLName xml.Name `xml:"http://jabber.org/features/compress compression"`
    Methods []string `xml:"method,omitempty"`
}

const (
    XMPP_STREAM_FEATURE_COMPRESSION_METHOD_ZLIB = "zlib"
    XMPP_STREAM_FEATURE_COMPRESSION_METHOD_LZW  = "lzw"
)

type XMPPStreamCompress struct {
    XMLName xml.Name `xml:"http://jabber.org/features/compress compress"`
    Methods []string `xml:"method,omitempty"`
}

type XMPPStreamCompressed struct {
    XMLName xml.Name `xml:"http://jabber.org/features/compress compressed"`
}

type XMPPStreamCompressionFailure struct {
    XMLName           xml.Name                                       `xml:"http://jabber.org/features/compress failure"`
    UnsupportedMethod *XMPPStreamCompressionFailureUnsupportedMethod `xml:",omitempty"`
    SetupFailed       *XMPPStreamCompressionFailureSetupFailed       `xml:",omitempty"`
    ProcessingFailed  *XMPPStreamCompressionFailureProcessingFailed  `xml:",omitempty"`
    XMPPStanzaErrorGroup
    Text *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamCompressionFailureUnsupportedMethod struct {
    XMLName xml.Name `xml:"unsupported-method"`
}

type XMPPStreamCompressionFailureSetupFailed struct {
    XMLName xml.Name `xml:"setup-failed"`
}

type XMPPStreamCompressionFailureProcessingFailed struct {
    XMLName xml.Name `xml:"processing-failed"`
}
