package xmpp

import (
	"errors"
	"strings"
)

// RFC6120

type BareJID struct {
	Local  string
	Domain string
}

type JID struct {
	Resource string
	BareJID
}

func (this *BareJID) ToString() string {
	return this.Local + "@" + this.Domain
}

func (this *JID) ToString() string {
	return this.BareJID.ToString() + "/" + this.Resource
}

func NewJID(jid string) (*JID, error) {
	p1 := strings.Split(jid, "@")
	if len(p1) != 2 {
		return nil, errors.New("Invalid jid: " + jid)
	}

	p2 := strings.SplitN(p1[1], "/", 2)
	objpnt := &JID{
		BareJID: BareJID{
			Local:  strings.TrimSpace(p1[0]),
			Domain: strings.TrimSpace(p2[0]),
		},
	}
	if len(p2) == 2 {
        objpnt.Resource = strings.TrimSpace(p2[1])
	}

    return objpnt, nil
}
