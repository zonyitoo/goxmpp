package stream

import (
    "bytes"
    "github.com/zonyitoo/goxmpp/protocol"
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

    <starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>

    <proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>

    <stream:features>
        <mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>
            <mechanism>SCRAM-SHA-1-PLUS</mechanism>
            <mechanism>SCRAM-SHA-1</mechanism>
            <mechanism>PLAIN</mechanism>
        </mechanisms>
    </stream:features>

    <auth xmlns="urn:ietf:params:xml:ns:xmpp-sasl"
          mechanism="SCRAM-SHA-1">biwsbj1qdWxpZXQscj1vTXNUQUF3QUFBQU1BQUFBTlAwVEFBQUFBQUJQVTBBQQ==</auth>

    <challenge xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>cj1vTXNUQUF3QUFBQU1BQUFBTlAwVEFBQUFBQUJQVTBBQWUxMjQ2OTViLTY5YTktNGRlNi05YzMwLWI1MWIzODA4YzU5ZSxzPU5qaGtZVE0wTURndE5HWTBaaTAwTmpkbUxUa3hNbVV0TkRsbU5UTm1ORE5rTURNeixpPTQwOTY=</challenge>

    <response xmlns="urn:ietf:params:xml:ns:xmpp-sasl">Yz1iaXdzLHI9b01zVEFBd0FBQUFNQUFBQU5QMFRBQUFBQUFCUFUwQUFlMTI0Njk1Yi02OWE5LTRkZTYtOWMzMC1iNTFiMzgwOGM1OWUscD1VQTU3dE0vU3ZwQVRCa0gyRlhzMFdEWHZKWXc9</response>

    <success xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>dj1wTk5ERlZFUXh1WHhDb1NFaVc4R0VaKzFSU289</success>

    <stream:features>
        <bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/>
    </stream:features>

    <iq id='yhc13a95' type='set'>
        <bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'>
            <resource>balcony</resource>
        </bind>
    </iq>

    <iq id='yhc13a96' type='result'>
        <bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'>
            <jid>juliet@im.example.com/balcony</jid>
        </bind>
    </iq>

    <message from='romeo@example.net/orchard'
             id='ju2ba41c'
             to='juliet@im.example.com/balcony'
             type='chat'
             xml:lang='en'>
        <body>Neither, fair saint, if either thee dislike.</body>
    </message>

    <presence from='romeo@example.net/orchard' xml:lang='en'>
        <show>dnd</show>
        <status>Wooing Juliet</status>
    </presence>

</stream:stream>
`
)

func test_StreamHeader(value interface{}, t *testing.T) {
    if _, ok := value.(*protocol.XMPPStream); !ok {
        t.Fatal("Error occurs while decoding Stream header")
    }
}

func test_StreamFeatureTLS(value interface{}, t *testing.T) {
    if features, ok := value.(*protocol.XMPPStreamFeatures); !ok {
        t.Fatal("Error occurs while decoding Stream Features")
    } else {
        if features.StartTLS == nil {
            t.Error("StartTLS should not be nil")
        } else if features.StartTLS.Required == nil {
            t.Error("StartTLS.Required should not be nil")
        }
    }
}

func test_StartTLS(value interface{}, t *testing.T) {
    if _, ok := value.(*protocol.XMPPStartTLS); !ok {
        t.Fatal("Error occurs while decoding StartTLS")
    }
}

func test_TLSProceed(value interface{}, t *testing.T) {
    if _, ok := value.(*protocol.XMPPTLSProceed); !ok {
        t.Fatal("Error occurs while decoding TLSProceed")
    }
}

func test_FeatureSASLMechanism(value interface{}, t *testing.T) {
    if features, ok := value.(*protocol.XMPPStreamFeatures); !ok {
        t.Fatal("Error occurs while decoding Stream Features")
    } else {
        if features.SASLMechanisms == nil {
            t.Error("SASLMechanisms should not be nil")
        } else {
            valid := []string{"SCRAM-SHA-1-PLUS", "SCRAM-SHA-1", "PLAIN"}
            for idx := range features.SASLMechanisms.Mechanisms {
                if valid[idx] != features.SASLMechanisms.Mechanisms[idx] {
                    t.Errorf("SASLMechanisms not match, %+v", features.SASLMechanisms.Mechanisms)
                    break
                }
            }
        }
    }
}

func test_SASLAuth(value interface{}, t *testing.T) {
    if saslauth, ok := value.(*protocol.XMPPSASLAuth); !ok {
        t.Fatal("Error occurs while decoding SASLAuth")
    } else {
        if saslauth.Data != "biwsbj1qdWxpZXQscj1vTXNUQUF3QUFBQU1BQUFBTlAwVEFBQUFBQUJQVTBBQQ==" {
            t.Error("SASLAuth data not match, got ", saslauth.Data)
        } else if saslauth.Mechanism != "SCRAM-SHA-1" {
            t.Error("SASLAuth mechanism not match")
        }
    }
}

func test_SASLChallenge(value interface{}, t *testing.T) {
    if saslchallenge, ok := value.(*protocol.XMPPSASLChallenge); !ok {
        t.Fatal("Error occurs while decoding SASLChallenge")
    } else {
        if saslchallenge.Data != "cj1vTXNUQUF3QUFBQU1BQUFBTlAwVEFBQUFBQUJQVTBBQWUxMjQ2OTViLTY5YTktNGRlNi05YzMwLWI1MWIzODA4YzU5ZSxzPU5qaGtZVE0wTURndE5HWTBaaTAwTmpkbUxUa3hNbVV0TkRsbU5UTm1ORE5rTURNeixpPTQwOTY=" {
            t.Error("SASLChallenge data not match")
        }
    }
}

func test_SASLResponse(value interface{}, t *testing.T) {
    if saslresponse, ok := value.(*protocol.XMPPSASLResponse); !ok {
        t.Fatal("Error occurs while decoding SASLResponse")
    } else {
        if saslresponse.Data != "Yz1iaXdzLHI9b01zVEFBd0FBQUFNQUFBQU5QMFRBQUFBQUFCUFUwQUFlMTI0Njk1Yi02OWE5LTRkZTYtOWMzMC1iNTFiMzgwOGM1OWUscD1VQTU3dE0vU3ZwQVRCa0gyRlhzMFdEWHZKWXc9" {
            t.Error("SASLResponse data not match")
        }
    }
}

func test_SASLSucceed(value interface{}, t *testing.T) {
    if saslsuc, ok := value.(*protocol.XMPPSASLSuccess); !ok {
        t.Fatal("Error occurs while decoding SASLSuccess")
    } else {
        if saslsuc.Data != "dj1wTk5ERlZFUXh1WHhDb1NFaVc4R0VaKzFSU289" {
            t.Error("SASLSuccess data not match")
        }
    }
}

func test_FeatureBind(value interface{}, t *testing.T) {
    if fea, ok := value.(*protocol.XMPPStreamFeatures); !ok {
        t.Fatal("Error occurs while decoding StreamFeature")
    } else {
        if fea.Bind == nil {
            t.Error("SASLFeature.Bind should not be nil")
        }
    }
}

func test_IQBindReq(value interface{}, t *testing.T) {
    if iq, ok := value.(*protocol.XMPPStanzaIQ); !ok {
        t.Fatal("Error occurs while decoding IQ")
    } else {
        if iq.Bind == nil {
            t.Error("XMPPStanzaIQ.Bind should not be nil")
        } else if iq.Bind.Resource != "balcony" {
            t.Errorf("XMPPStanzaIQ.Bind.Resource not match %s != %s", iq.Bind.Resource, "balcony")
        } else if iq.Id != "yhc13a95" {
            t.Errorf("XMPPStanzaIQ.Id not match %s != %s", iq.Id, "yhc13a95")
        } else if iq.Type != "set" {
            t.Errorf("XMPPStanzaIQ.Type not match, %s != %s", iq.Type, "set")
        }
    }
}

func test_IQBindResp(value interface{}, t *testing.T) {
    if iq, ok := value.(*protocol.XMPPStanzaIQ); !ok {
        t.Fatal("Error occurs while decoding IQ")
    } else {
        if iq.Bind == nil {
            t.Error("XMPPStanzaIQ.Bind should not be nil")
        } else if iq.Bind.JID != "juliet@im.example.com/balcony" {
            t.Errorf("XMPPStanzaIQ.Bind.JID not match, %s != %s", iq.Bind.JID, "juliet@im.example.com/balcony")
        } else if iq.Id != "yhc13a96" {
            t.Errorf("XMPPStanzaIQ.Id not match, %s != %s", iq.Id, "yhc13a96")
        } else if iq.Type != "result" {
            t.Errorf("XMPPStanzaIQ.Type not match, %s != %s", iq.Type, "result")
        }
    }
}

func test_Message(value interface{}, t *testing.T) {
    if msg, ok := value.(*protocol.XMPPStanzaMessage); !ok {
        t.Fatal("Error occurs while decoding Message")
    } else {
        if msg.Body == nil {
            t.Error("XMPPStanzaMessage.Body should not be nil")
        } else if msg.Body.Data != "Neither, fair saint, if either thee dislike." {
            t.Error("XMPPStanzaMessage.Body.Data not match")
        } else if msg.From != "romeo@example.net/orchard" {
            t.Errorf("XMPPStanzaMessage.From not match, %s != %s", msg.From, "romeo@example.net/orchard")
        } else if msg.To != "juliet@im.example.com/balcony" {
            t.Errorf("XMPPStanzaMessage.To not match, %s != %s", msg.To, "juliet@im.example.com/balcony")
        } else if msg.Id != "ju2ba41c" {
            t.Errorf("XMPPStanzaMessage.Id not match, %s != %s", msg.Id, "ju2ba41c")
        } else if msg.Type != protocol.XMPP_STANZA_MESSAGE_TYPE_CHAT {
            t.Errorf("XMPPStanzaMessage.Type not match, %s != %s", msg.Type, protocol.XMPP_STANZA_MESSAGE_TYPE_CHAT)
        } else if msg.XMLLang != "en" {
            t.Errorf("XMPPStanzaMessage.XMLLang not match, %s != %s", msg.XMLLang, "en")
        }
    }
}

func test_Presence(value interface{}, t *testing.T) {
    if pres, ok := value.(*protocol.XMPPStanzaPresence); !ok {
        t.Fatal("Error occurs while decoding Presence")
    } else {
        if pres.From != "romeo@example.net/orchard" {
            t.Errorf("XMPPStanzaPresence.From not match, %s != %s", pres.From, "romeo@example.net/orchard")
        } else if pres.XMLLang != "en" {
            t.Errorf("XMPPStanzaPresence.XMLLang not match, %s != %s", pres.XMLLang, "en")
        } else if pres.Show != protocol.XMPP_STANZA_PRESENCE_SHOW_DND {
            t.Errorf("XMPPStanzaPresence.Show not match, %s != %s", pres.Show, protocol.XMPP_STANZA_PRESENCE_SHOW_DND)
        } else if pres.Status == nil {
            t.Errorf("XMPPStanzaPresence.Status should not be nil")
        } else if pres.Status.Data != "Wooing Juliet" {
            t.Errorf("XMPPStanzaPresence.Status not match, %s != %s", pres.Status.Data, "Wooing Juliet")
        }
    }
}

func Test_Decoder(t *testing.T) {
    r := bytes.NewBufferString(xmpp_stream_sample)
    decoder := NewDecoder(r)

    tests := []func(interface{}, *testing.T){
        test_StreamHeader,
        test_StreamFeatureTLS,
        test_StartTLS,
        test_TLSProceed,
        test_FeatureSASLMechanism,
        test_SASLAuth,
        test_SASLChallenge,
        test_SASLResponse,
        test_SASLSucceed,
        test_FeatureBind,
        test_IQBindReq,
        test_IQBindResp,
        test_Message,
        test_Presence,
    }

    for _, test_func := range tests {
        if value, err := decoder.GetNextElement(); err != nil {
            t.Fatal(err)
        } else {
            t.Logf("%+v", value)
            test_func(value, t)
        }
    }
}
