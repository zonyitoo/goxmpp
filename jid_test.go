package xmpp

import (
    "testing"
)

const (
    BARE_JID_STR   = "zonyitoo@gmail.com"
    FULL_JID_STR   = "zonyitoo@gmail.com/resourcepath"
    DOMAIN_JID_STR = "gmail.com"
)

func Test_NewJIDFromString(t *testing.T) {
    bare_jid, err := NewJIDFromString(BARE_JID_STR)
    if err != nil {
        t.Error(err)
        t.Fail()
    }

    bare_jid_valid := &JID{
        BareJID: BareJID{
            Local:  "zonyitoo",
            Domain: "gmail.com",
        },
        Resource: "",
    }

    if *bare_jid != *bare_jid_valid || bare_jid.String() != BARE_JID_STR {
        t.Errorf("Error occurs while parsing %s", BARE_JID_STR)
    }

    full_jid, err := NewJIDFromString(FULL_JID_STR)
    if err != nil {
        t.Error(err)
        t.Fail()
    }

    full_jid_valid := &JID{
        BareJID: BareJID{
            Local:  "zonyitoo",
            Domain: "gmail.com",
        },
        Resource: "resourcepath",
    }

    if *full_jid != *full_jid_valid || full_jid.String() != FULL_JID_STR {
        t.Errorf("Error occurs while parsing %s", FULL_JID_STR)
    }

    domain_jid, err := NewJIDFromString(DOMAIN_JID_STR)
    if err != nil {
        t.Error(err)
        t.Fail()
    }

    domain_jid_valid := &JID{
        BareJID: BareJID{
            Local:  "",
            Domain: "gmail.com",
        },
        Resource: "",
    }

    if *domain_jid != *domain_jid_valid || domain_jid.String() != DOMAIN_JID_STR {
        t.Errorf("Error occurs while parsing %s", DOMAIN_JID_STR)
    }

    _, err = NewJIDFromString("")
    if err == nil {
        t.Error("Empty should failed")
        t.Fail()
    }
}

func Test_JIDToString(t *testing.T) {
    bjid := BareJID{
        Local:  "test",
        Domain: "test.domain",
    }

    if bjid.String() != "test@test.domain" {
        t.Errorf("%s is not equals to %s", bjid.String(), "test@test.domain")
    }

    jid := NewJID("test", "test.domain", "resource")

    if jid.String() != "test@test.domain/resource" {
        t.Errorf("%s is not equals to %s", jid.String(), "test@test.domain/resource")
    }
}
