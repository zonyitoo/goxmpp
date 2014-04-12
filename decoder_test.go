package xmpp

import (
    "bytes"
    "encoding/xml"
    "testing"
)

const (
    xmpp_stream_sample = `
<?xml version='1.0'?>
<stream:stream
    from='juliet@im.example.com'
    to='im.example.com'
    version='1.0'
    xml:lang='en'
    xmlns='jabber:client'
    xmlns:stream='http://etherx.jabber.org/streams'>

    <stream:features>
        <starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>
            <required/>
        </starttls>
    </stream:features>

</stream:stream>
`
)

func Test_Decoder(t *testing.T) {
    r := bytes.NewBufferString(xmpp_stream_sample)
    decoder := NewDecoder(r)

    var value interface{}
    var err error
    value, err = decoder.GetNextElement()
    if err != nil {
        t.Error(err)
    }
    if _, ok := value.(xml.ProcInst); !ok {
        t.Error("Error occurs while decoding ProcInst")
    }
    value, err = decoder.GetNextElement()
    if err != nil {
        t.Error(err)
    }
    if _, ok := value.(*XMPPStream); !ok {
        t.Error("Error occurs while decoding Stream header")
    }
    value, err = decoder.GetNextElement()
    if err != nil {
        t.Error(err)
    }
    if _, ok := value.(*XMPPStreamFeatures); !ok {
        t.Error("Error occurs while decoding Stream Features")
    }

    value, err = decoder.GetNextElement()
    if _, ok := value.(*XMPPStreamEnd); !ok {
        t.Error("Error occurs while decoding Stream End")
    }
}
