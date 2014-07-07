package xmpp

import (
    "errors"
    "fmt"
    "regexp"
)

var validJid *regexp.Regexp = regexp.MustCompile(`^(\w+)@(\w+\.\w+)+/?(\w+)?$`)

type BareJID struct {
    Local  string
    Domain string
}

type JID struct {
    BareJID
    Resource string
}

func ValidateJID(jidstr string) bool {
    return validJid.MatchString(jidstr)
}

func NewJID(local, domain, resource string) *JID {
    return &JID{
        BareJID: BareJID{
            Local:  local,
            Domain: domain,
        },
        Resource: resource,
    }
}

func NewJIDFromString(jid string) (*JID, error) {
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
    return NewJID(local, domain, resource), nil
}

func (jid *JID) String() string {
    if jid.Resource != "" {
        return fmt.Sprintf("%s@%s/%s", jid.Local, jid.Domain, jid.Resource)
    } else {
        return fmt.Sprintf("%s@%s", jid.Local, jid.Domain)
    }
}

func (barejid *BareJID) String() string {
    return fmt.Sprintf("%s@%s", barejid.Local, barejid.Domain)
}
