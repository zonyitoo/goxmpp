package xmpp

import (
    "encoding/xml"
    "reflect"
    "testing"
    "time"
)

func Test_RPCParamValue(t *testing.T) {
    testxml := `
        <query xmlns='jabber:iq:rpc'>
            <methodCall>
                <methodName>example.getStateName</methodName>
                <params>
                    <param>
                        <value><i4>6</i4></value>
                    </param>
                    <param>
                        <value><string>Colorado</string></value>
                    </param>
                    <param>
                        <value><base64>c29tZSBkYXRhIHdpdGggACBhbmQg77u/</base64></value>
                    </param>
                    <param>
                        <value><double>11.2</double></value>
                    </param>
                    <param>
                        <value><boolean>1</boolean></value>
                    </param>
                    <param>
                        <value><dateTime.iso8601>20100101T15:04:05</dateTime.iso8601></value>
                    </param>
                </params>
            </methodCall>

            <methodResponse>
                <params>
                    <param>
                        <value>
                            <array>
                                <data>
                                    <value><int>10</int></value>
                                    <value><string>FUCK</string></value>
                                </data>
                            </array>
                        </value>
                    </param>

                    <param>
                        <value>
                            <struct>
                                <member>
                                    <name>Hello!!</name>
                                    <value><string>World</string></value>
                                </member>
                                <member>
                                    <name>Fuck</name>
                                    <value><int>11</int></value>
                                </member>
                            </struct>
                        </value>
                    </param>
                </params>
            </methodResponse>
        </query>
    `

    validval_call_time, _ := time.Parse(time_ISO8601_FORMAT, "20100101T15:04:05")

    validval_call := []interface{}{
        6,
        "Colorado",
        []byte("some data with \x00 and \ufeff"),
        11.2,
        true,
        validval_call_time,
    }

    validval_resp := []interface{}{
        []interface{}{
            10,
            "FUCK",
        },
        map[string]interface{}{
            "Hello!!": "World",
            "Fuck":    11,
        },
    }

    query := &XMPPStanzaIQRPCQuery{}
    if err := xml.Unmarshal([]byte(testxml), query); err != nil {
        t.Error(err)
    }

    for idx, val := range query.MethodCall.Params {
        validate_value(&val.Value, validval_call[idx], t)
    }

    for idx, val := range query.MethodResponse.Params {
        validate_value(&val.Value, validval_resp[idx], t)
    }

}

func validate_value(value *XMPPStanzaIQRPCParamValue, validval interface{}, t *testing.T) {
    val, err := value.Value()
    if err != nil {
        t.Error(err)
    } else {
        if !reflect.DeepEqual(val, validval) {
            t.Errorf("%s != %s", val, validval)
        }
    }
}

func Test_RPCParamValueSet(t *testing.T) {
    val := XMPPStanzaIQRPCParamValue{}
    val.SetValue(1)
    if ival, err := val.Value(); err != nil {
        t.Error(err)
    } else {
        if iv, ok := ival.(int); !ok || iv != 1 {
            t.Error("Error occurs while setting int value")
        }
    }

    val.SetValue("Hello")
    if ival, err := val.Value(); err != nil {
        t.Error(err)
    } else {
        if iv, ok := ival.(string); !ok || iv != "Hello" {
            t.Error("Error occurs while setting string value")
        }
    }

    val.SetValue([]byte("\x00\x01\x02\x03\x04"))
    if ival, err := val.Value(); err != nil {
        t.Error(err)
    } else {
        if iv, ok := ival.([]byte); !ok || !reflect.DeepEqual(iv, []byte("\x00\x01\x02\x03\x04")) {
            t.Error("Error occurs while setting base64 value")
        }
    }

    arr := []interface{}{"Fuck", "Hello", 1, 1.1}
    val.SetValue(arr)
    if ival, err := val.Value(); err != nil {
        t.Error(err)
    } else {
        if iv, ok := ival.([]interface{}); !ok || !reflect.DeepEqual(iv, arr) {
            t.Error("Error occurs while setting Array value")
        }
    }
}
