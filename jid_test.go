package xmpp

import (
    "testing"
)

const (
    BARE_JID_STR = "zonyitoo@gmail.com"
    FULL_JID_STR = "zonyitoo@gmail.com/resourcepath"
)

func Test_NewJIDFromString(t *testing.T) {
    bare_jid, err := NewXMPPJIDFromString(BARE_JID_STR)
    if err != nil {
        t.Error(err)
    }

    bare_jid_valid := &XMPPJID{
        Local:    "zonyitoo",
        Domain:   "gmail.com",
        Resource: "",
    }

    if *bare_jid != *bare_jid_valid || bare_jid.String() != BARE_JID_STR {
        t.Errorf("Error occurs while parsing %s", BARE_JID_STR)
    }

    full_jid, err := NewXMPPJIDFromString(FULL_JID_STR)
    if err != nil {
        t.Error(err)
    }

    full_jid_valid := &XMPPJID{
        Local:    "zonyitoo",
        Domain:   "gmail.com",
        Resource: "resourcepath",
    }

    if *full_jid != *full_jid_valid || full_jid.String() != FULL_JID_STR {
        t.Errorf("Error occurs while parsing %s", FULL_JID_STR)
    }
}
