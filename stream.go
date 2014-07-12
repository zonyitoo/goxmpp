package xmpp

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
)

const (
    STREAM_TYPE_CLIENT = iota
    STREAM_TYPE_SERVER
)

const (
    STREAM_STAT_INIT    = 0
    STREAM_STAT_STARTED = 1 << iota
    STREAM_STAT_TLS_PROCEED
    STREAM_STAT_SASL_SUCCEEDED
    STREAM_STAT_CLOSED
)

type Stream interface {
    ServerConfig() *ServerConfig
    JID() *JID
    SetJID(*JID)
    SendBytes([]byte)
    StartStream(stype int, from, to, version, lang string)
    SendElement(interface{}) error
    SendErrorAndClose(interface{}) error
    Close()
    EndStream()
    State() int
    SetState(int)
}

type StreamVersion struct {
    Major int
    Minor int
}

func (sv *StreamVersion) FromString(vstr string) error {
    sp := strings.Split(vstr, ".")
    if len(sp) != 2 {
        return errors.New("Malformed version string")
    }
    if major, err := strconv.Atoi(sp[0]); err != nil {
        return err
    } else {
        sv.Major = major
    }

    if sv.Major < 0 {
        return errors.New("Invalid major version")
    }

    if minor, err := strconv.Atoi(sp[1]); err != nil {
        return err
    } else {
        sv.Minor = minor
    }

    if sv.Minor < 0 {
        return errors.New("Invalid minor version")
    }

    return nil
}

func (sv *StreamVersion) String() string {
    return fmt.Sprintf("%d.%d", sv.Major, sv.Minor)
}

func (sv *StreamVersion) GreaterOrEqualTo(rhs *StreamVersion) bool {
    if sv.Major == rhs.Major {
        return sv.Minor >= rhs.Minor
    } else {
        return sv.Major >= rhs.Major
    }
}
