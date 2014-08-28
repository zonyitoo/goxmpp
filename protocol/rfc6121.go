package protocol

import (
    "encoding/xml"
)

const (
    XMLNS_JABBER_IQ_ROSTER = "jabber:iq:roster"
)

// RFC6121
type XMPPStanzaIQRosterQuery struct {
    XMLName xml.Name                 `xml:"jabber:iq:roster query"`
    Ver     string                   `xml:"ver,attr,omitempty"`
    Item    []XMPPStanzaIQRosterItem `xml:",omitempty"`
}

type XMPPStanzaIQRosterItem struct {
    XMLName      xml.Name `xml:"item"`
    JID          string   `xml:"jid,attr"`
    Name         string   `xml:"name,attr,omitempty"`
    Subscription string   `xml:"subscription,attr,omitempty"`
    Ask          string   `xml:"ask,attr,omitempty"`
    Groups       []string `xml:"group,omitempty"`
    Approved     bool     `xml:"approved,omitempty"`
}

const (
    // the user does not have a subscription to the contact's
    // presence, and the contact does not have a subscription to the
    // user's presence; this is the default value, so if the subscription
    // attribute is not included then the state is to be understood as
    // "none"
    XMPP_IQ_ROSTER_ITEM_SUBSCRIPTION_TYPE_NONE = "none"
    // the user has a subscription to the contact's presence, but the
    // contact does not have a subscription to the user's presence
    XMPP_IQ_ROSTER_ITEM_SUBSCRIPTION_TYPE_TO = "to"
    // the contact has a subscription to the user's presence, but the
    // user does not have a subscription to the contact's presence
    XMPP_IQ_ROSTER_ITEM_SUBSCRIPTION_TYPE_FROM = "from"
    // the user and the contact have subscriptions to each other's
    // presence (also called a "mutual subscription")
    XMPP_IQ_ROSTER_ITEM_SUBSCRIPTION_TYPE_BOTH   = "both"
    XMPP_IQ_ROSTER_ITEM_SUBSCRIPTION_TYPE_REMOVE = "remove"
)

const (
    XMPP_IQ_ROSTER_ITEM_ASK_SUBSCRIBE = "subscribe"
)
