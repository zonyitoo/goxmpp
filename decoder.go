package xmpp

import (
    "encoding/xml"
    "errors"
    "io"
    // "log"
    "bytes"
    "github.com/zonyitoo/goxmpp/protocol"
)

var (
    DecoderBadFormatError              = errors.New("Bad format")
    DecoderUnexpectedEndOfElementError = errors.New("Unexpected end of element")
    DecoderRestrictedXMLError          = errors.New("Restricted XML")
)

type Decoder struct {
    xmlDecoder *xml.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
    return &Decoder{
        xmlDecoder: xml.NewDecoder(r),
    }
}

func (d *Decoder) ParseElement(startToken xml.StartElement) (protocol.Protocol, error) {
    var element interface{}
    switch startToken.Name {
    case protocol.TAG_STREAM:
        streamElem := &protocol.XMPPStream{}
        for _, attr := range startToken.Attr {
            switch attr.Name {
            case xml.Name{Space: "", Local: "from"}:
                streamElem.From = attr.Value
            case xml.Name{Space: "", Local: "to"}:
                streamElem.To = attr.Value
            case xml.Name{Space: "", Local: "id"}:
                streamElem.Id = attr.Value
            case xml.Name{Space: "", Local: "version"}:
                streamElem.Version = attr.Value
            case xml.Name{Space: "http://www.w3.org/XML/1998/namespace", Local: "lang"}:
                streamElem.XMLLang = attr.Value
            case xml.Name{Space: "", Local: "xmlns"}:
                streamElem.Xmlns = attr.Value
            }
        }
        streamElem.XMLName = startToken.Name
        return streamElem, nil
    case protocol.TAG_STREAM_FEATURES:
        element = &protocol.XMPPStreamFeatures{}
    case protocol.TAG_STREAM_ERROR:
        element = &protocol.XMPPStreamError{}

    case protocol.TAG_TLS_START:
        element = &protocol.XMPPStartTLS{}
    case protocol.TAG_TLS_PROCEED:
        element = &protocol.XMPPTLSProceed{}
    case protocol.TAG_TLS_FAILURE:
        element = &protocol.XMPPTLSFailure{}
    case protocol.TAG_TLS_ABORT:
        element = &protocol.XMPPTLSAbort{}

    case protocol.TAG_SASL_SUCCESS:
        element = &protocol.XMPPSASLSuccess{}
    case protocol.TAG_SASL_RESPONSE:
        element = &protocol.XMPPSASLResponse{}
    case protocol.TAG_SASL_FAILURE:
        element = &protocol.XMPPSASLFailure{}
    case protocol.TAG_SASL_CHALLENGE:
        element = &protocol.XMPPSASLChallenge{}
    case protocol.TAG_SASL_AUTH:
        element = &protocol.XMPPSASLAuth{}

    case protocol.TAG_STANZA_IQ_CLIENT, protocol.TAG_STANZA_IQ_SERVER:
        element = &protocol.XMPPStanzaIQ{}
    case protocol.TAG_STANZA_PRESENCE_SERVER, protocol.TAG_STANZA_PRESENCE_CLIENT:
        element = &protocol.XMPPStanzaPresence{}
    case protocol.TAG_STANZA_MESSAGE_SERVER, protocol.TAG_STANZA_MESSAGE_CLIENT:
        element = &protocol.XMPPStanzaMessage{}

    // Extensions
    // XEP-0138
    case protocol.TAG_STREAM_COMPRESSION_COMPRESS:
        element = &protocol.XMPPStreamCompressionCompress{}
    case protocol.TAG_STREAM_COMPRESSION_FAILURE:
        element = &protocol.XMPPStreamCompressionFailure{}
    case protocol.TAG_STREAM_COMPRESSION_COMPRESSED:
        element = &protocol.XMPPStreamCompressionCompressed{}

    default:
        element = &protocol.XMPPCustom{}
    }

    if err := d.xmlDecoder.DecodeElement(element, &startToken); err != nil {
        return nil, err
    }

    return element, nil
}

func (d *Decoder) GetNextElement() (protocol.Protocol, error) {
    // Move to First StartElement
    for {
        token, err := d.xmlDecoder.Token()
        if err != nil {
            return nil, err
        }
        switch t := token.(type) {
        case xml.StartElement:
            return d.ParseElement(t)
        case xml.ProcInst:
            continue
        case xml.EndElement:
            if t.Name == protocol.TAG_STREAM {
                return &protocol.XMPPStreamEnd{
                    XMLName: t.Name,
                }, nil
            } else {
                return nil, DecoderUnexpectedEndOfElementError
            }
        case xml.CharData:
            if len(bytes.TrimSpace(t)) != 0 {
                return nil, DecoderBadFormatError
            }
        case xml.Comment:
            return nil, DecoderRestrictedXMLError
        }
    }
}
