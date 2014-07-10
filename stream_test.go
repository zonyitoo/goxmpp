package xmpp

import (
    "testing"
)

func Test_StreamVersion(t *testing.T) {
    sv := StreamVersion{}
    if err := sv.FromString("1.0"); err != nil {
        t.Error(err)
    }
    if sv.Major != 1 || sv.Minor != 0 {
        t.Errorf("%d.%d is not equal to 1.0", sv.Major, sv.Minor)
    }

    gt := StreamVersion{}
    if err := sv.FromString("2.13"); err != nil {
        t.Error(err)
    }

    lt := StreamVersion{}
    if err := sv.FromString("2.4"); err != nil {
        t.Error(err)
    }

    if gt.GreaterOrEqualTo(&lt) != true {
        t.Errorf("%s should greater than %s", gt.String(), lt.String())
    }
}
