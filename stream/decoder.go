package stream

import (
    "encoding/xml"
    "errors"
    "io"
    // "log"
    "bytes"
    "github.com/zonyitoo/goxmpp/protocol"
    "reflect"
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
    if startToken.Name == protocol.TAG_STREAM {
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
    } else {
        if t, ok := protocol.TAG_MAP[startToken.Name]; !ok {
            element = &protocol.XMPPCustom{}
        } else {
            element = reflect.New(t).Interface()
        }
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
                return &protocol.XMPPStreamEnd{}, nil
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
