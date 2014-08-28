package protocol

import (
    "encoding/xml"
)

const (
    XMLNS_JABBER_X_DATA = "jabber:x:data"
)

type XMPPXData struct {
    XMLName xml.Name         `xml:"jabber:x:data x"`
    Type    string           `xml:"type,attr,omitempty"`
    Fields  []XMPPXDataField `xml:",omitempty"`
    Title   string           `xml:"title,omitempty"`
}

const (
    XMPP_X_DATA_TYPE_CANCEL = "cancel"
    XMPP_X_DATA_TYPE_FORM   = "form"
    XMPP_X_DATA_TYPE_RESULT = "result"
    XMPP_X_DATA_TYPE_SUBMIT = "submit"
)

type XMPPXDataField struct {
    XMLName        xml.Name                      `xml:"field"`
    Var            string                        `xml:"var,attr,omitempty"`
    Label          string                        `xml:"label,attr,omitempty"`
    Type           string                        `xml:"type,attr,omitempty"`
    Values         []string                      `xml:"value,omitempty"`
    OptionalValues []XMPPXDataFieldOptionalValue `xml:",omitempty"`
    Required       *XMPPRequired                 `xml:",omitempty"`
    Desc           string                        `xml:"desc,omitempty"`
}

const (
    XMPP_X_DATA_FIELD_TYPE_BOOLEAN      = "boolean"
    XMPP_X_DATA_FIELD_TYPE_FIXED        = "fixed"
    XMPP_X_DATA_FIELD_TYPE_HIDDEN       = "hidden"
    XMPP_X_DATA_FIELD_TYPE_JID_MULTI    = "jid-multi"
    XMPP_X_DATA_FIELD_TYPE_JID_SINGLE   = "jid-single"
    XMPP_X_DATA_FIELD_TYPE_LIST_MULTI   = "list-multi"
    XMPP_X_DATA_FIELD_TYPE_LIST_SINGLE  = "list-single"
    XMPP_X_DATA_FIELD_TYPE_TEXT_MULTI   = "text-multi"
    XMPP_X_DATA_FIELD_TYPE_TEXT_PRIVATE = "text-private"
    XMPP_X_DATA_FIELD_TYPE_TEXT_SINGLE  = "text-single"
)

type XMPPXDataFieldOptionalValue struct {
    XMLName xml.Name `xml:"option"`
    Value   string   `xml:"value"`
    Label   string   `xml:"label,attr,omitempty"`
}

type XMPPXDataItem struct {
    XMLName xml.Name         `xml:"item"`
    Fields  []XMPPXDataField `xml:",omitempty"`
}

type XMPPXDataReported struct {
    XMLName xml.Name         `xml:"reported"`
    Fields  []XMPPXDataField `xml:",omitempty"`
}
