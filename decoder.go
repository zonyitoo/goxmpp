package xmpp

import (
    "encoding/xml"
    "errors"
    "io"
    // "log"
    "bytes"
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

func (d *Decoder) ParseElement(startToken xml.StartElement) (interface{}, error) {
    var element interface{}
    switch startToken.Name {
    case TAG_STREAM:
        streamElem := &XMPPStream{}
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
    case TAG_STREAM_FEATURES:
        element = &XMPPStreamFeatures{}

    case TAG_TLS_START:
        element = &XMPPStartTLS{}
    case TAG_TLS_PROCEED:
        element = &XMPPTLSProceed{}
    case TAG_TLS_FAILURE:
        element = &XMPPTLSFailure{}
    case TAG_TLS_ABORT:
        element = &XMPPTLSAbort{}

    case TAG_SASL_SUCCESS:
        element = &XMPPSASLSuccess{}
    case TAG_SASL_RESPONSE:
        element = &XMPPSASLResponse{}
    case TAG_SASL_FAILURE:
        element = &XMPPSASLFailure{}
    case TAG_SASL_CHALLENGE:
        element = &XMPPSASLChallenge{}
    case TAG_SASL_AUTH:
        element = &XMPPSASLAuth{}

    case TAG_STANZA_IQ_CLIENT, TAG_STANZA_IQ_SERVER:
        element = &XMPPStanzaIQ{}
    case TAG_STANZA_PRESENCE_SERVER, TAG_STANZA_PRESENCE_CLIENT:
        element = &XMPPStanzaPresence{}
    case TAG_STANZA_MESSAGE_SERVER, TAG_STANZA_MESSAGE_CLIENT:
        element = &XMPPStanzaMessage{}

    // Extensions
    // XEP-0138
    case TAG_STREAM_COMPRESSION_COMPRESS:
        element = &XMPPStreamCompressionCompress{}
    case TAG_STREAM_COMPRESSION_FAILURE:
        element = &XMPPStreamCompressionFailure{}
    case TAG_STREAM_COMPRESSION_COMPRESSED:
        element = &XMPPStreamCompressionCompressed{}

    default:
        element = &XMPPCustom{}
    }

    if err := d.xmlDecoder.DecodeElement(element, &startToken); err != nil {
        return nil, err
    }

    return element, nil
}

func (d *Decoder) GetNextElement() (interface{}, error) {
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
            return t, nil
        case xml.EndElement:
            if t.Name == TAG_STREAM {
                return &XMPPStreamEnd{
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
