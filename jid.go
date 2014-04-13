package xmpp

import (
    "errors"
    "fmt"
    "regexp"
)

var validJid *regexp.Regexp = regexp.MustCompile(`^(\w+)@(\w+\.\w+)/?(\w+)?$`)

type XMPPJID struct {
    Local    string
    Domain   string
    Resource string
}

func NewXMPPJID(local, domain, resource string) *XMPPJID {
    return &XMPPJID{
        Local:    local,
        Domain:   domain,
        Resource: resource,
    }
}

func NewXMPPJIDFromString(jid string) (*XMPPJID, error) {
    if !validJid.MatchString(jid) {
        return nil, errors.New("Malformed JID")
    }

    substrs := validJid.FindStringSubmatch(jid)
    if len(substrs) < 3 {
        return nil, errors.New("Malformed JID")
    }

    local := substrs[1]
    domain := substrs[2]
    resource := ""
    if len(substrs) == 4 {
        resource = substrs[3]
    }
    return NewXMPPJID(local, domain, resource), nil
}

func (jid *XMPPJID) String() string {
    if jid.Resource != "" {
        return fmt.Sprintf("%s@%s/%s", jid.Local, jid.Domain, jid.Resource)
    } else {
        return fmt.Sprintf("%s@%s", jid.Local, jid.Domain)
    }
}
