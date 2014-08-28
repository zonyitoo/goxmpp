package protocol

import (
    "bytes"
    "encoding/base64"
    "encoding/xml"
    "errors"
    "fmt"
    "strconv"
    "strings"
    "time"
)

const (
    XMLNS_JABBER_IQ_RPC = "jabber:iq:rpc"

    time_ISO8601_FORMAT = "20060102T15:04:05"
)

type XMPPStanzaIQRPCQuery struct {
    XMLName        xml.Name                       `xml:"jabber:iq:rpc query"`
    MethodCall     *XMPPStanzaIQRPCMethodCall     `xml:""`
    MethodResponse *XMPPStanzaIQRPCMethodResponse `xml:""`
}

type XMPPStanzaIQRPCMethodCall struct {
    XMLName    xml.Name               `xml:"methodCall"`
    MethodName string                 `xml:"methodName"`
    Params     []XMPPStanzaIQRPCParam `xml:"params>param,omitempty"`
}

type XMPPStanzaIQRPCMethodResponse struct {
    XMLName xml.Name                  `xml:"methodResponse"`
    Params  []XMPPStanzaIQRPCParam    `xml:"params>param,omitempty"`
    Fault   XMPPStanzaIQRPCParamValue `xml:"fault>value"`
}

type XMPPStanzaIQRPCParam struct {
    XMLName xml.Name                  `xml:"param"`
    Value   XMPPStanzaIQRPCParamValue `xml:"value"`
}

type XMPPStanzaIQRPCParamValue struct {
    XMLName xml.Name `xml:"value"`
    Value_  string   `xml:",innerxml"`
}

func (val *XMPPStanzaIQRPCParamValue) Value() (interface{}, error) {
    decoder := xml.NewDecoder(bytes.NewBufferString(val.Value_))
    token, err := decoder.Token()
    if err != nil {
        return nil, err
    }

    return val.parse_value(decoder, token)
}

func btoi(b bool) int {
    if b == false {
        return 0
    }
    return 1
}

func value_to_xml(v interface{}) (string, error) {
    switch t := v.(type) {
    case int, int8, int64, int32, *int, *int8, *int64, *int32:
        return fmt.Sprintf("<int>%d</int>", t), nil
    case string:
        buf := bytes.NewBuffer(make([]byte, 0))
        err := xml.EscapeText(buf, []byte(t))
        if err != nil {
            return "", err
        }
        return fmt.Sprintf("<string>%s</string>", buf.String()), nil
    case []byte:
        return fmt.Sprintf("<base64>%s</base64>", base64.StdEncoding.EncodeToString(t)), nil
    case float32, float64, *float32, *float64:
        return fmt.Sprintf("<double>%f</double>", t), nil
    case bool:
        return fmt.Sprintf("<boolean>%d</boolean>", btoi(t)), nil
    case time.Time:
        buf := bytes.NewBuffer(make([]byte, 0))
        err := xml.EscapeText(buf, []byte(t.Format(time_ISO8601_FORMAT)))
        if err != nil {
            return "", err
        }
        return fmt.Sprintf("<dateTime.iso8601>%s</dateTime.iso8601>", buf.String()), nil
    case *time.Time:
        buf := bytes.NewBuffer(make([]byte, 0))
        err := xml.EscapeText(buf, []byte(t.Format(time_ISO8601_FORMAT)))
        if err != nil {
            return "", err
        }
        return fmt.Sprintf("<dateTime.iso8601>%s</dateTime.iso8601>", buf.String()), nil
    case []interface{}:
        tmp_str_arr := []string{"<array><data>"}
        for _, v := range t {
            subval, err := value_to_xml(v)
            if err != nil {
                return "", err
            }
            tmp_str_arr = append(tmp_str_arr, fmt.Sprintf("<value>%s</value>", subval))
        }
        tmp_str_arr = append(tmp_str_arr, "</data></array>")
        return strings.Join(tmp_str_arr, ""), nil
    case map[string]interface{}:
        tmp_str_arr := []string{"<struct>"}
        for name, item := range t {
            tmp_str_arr = append(tmp_str_arr, "<member>")
            tmp_str_arr = append(tmp_str_arr, fmt.Sprintf("<name>%s</name>", name))
            subval, err := value_to_xml(item)
            if err != nil {
                return "", err
            }
            tmp_str_arr = append(tmp_str_arr, fmt.Sprintf("<value>%s</value>", subval))
            tmp_str_arr = append(tmp_str_arr, "</member>")
        }
        tmp_str_arr = append(tmp_str_arr, "</struct>")
        return strings.Join(tmp_str_arr, ""), nil
    default:
        return "", errors.New("Unrecognized type")
    }
}

func (val *XMPPStanzaIQRPCParamValue) SetValue(v interface{}) error {
    pv, err := value_to_xml(v)
    if err != nil {
        return err
    }
    val.Value_ = pv
    return nil
}

func (val *XMPPStanzaIQRPCParamValue) parse_value(decoder *xml.Decoder, token xml.Token) (interface{}, error) {
    for {
        if se, ok := token.(xml.StartElement); ok {
            switch se.Name.Local {
            case "i1", "i2", "i4", "i8", "int":
                return val.parse_int(decoder, se)
            case "string":
                return val.parse_string(decoder, se)
            case "base64":
                return val.parse_base64(decoder, se)
            case "double":
                return val.parse_double(decoder, se)
            case "boolean":
                return val.parse_bool(decoder, se)
            case "dateTime.iso8601":
                return val.parse_datetime(decoder, se)
            case "array":
                return val.parse_array(decoder, se)
            case "struct":
                return val.parse_struct(decoder, se)
            default:
                return nil, errors.New(fmt.Sprintf("Unrecognized value type `%s`", se.Name))
            }
            break
        } else if _, ok := token.(xml.EndElement); ok {
            return nil, errors.New("Unexpected EndElement")
        }

        _token, err := decoder.Token()
        if err != nil {
            return nil, err
        }
        token = _token
    }

    return nil, errors.New("Unexpected element")
}

func (val *XMPPStanzaIQRPCParamValue) parse_int(decoder *xml.Decoder, se xml.StartElement) (interface{}, error) {
    if t, err := decoder.Token(); err != nil {
        return nil, err
    } else {
        if cdata, ok := t.(xml.CharData); !ok {
            return nil, errors.New(fmt.Sprintf("Unexpected element. Expecting Chardata, got %+v", t))
        } else {
            if xval, err := strconv.Atoi(string(cdata)); err != nil {
                return nil, err
            } else {
                return xval, nil
            }
        }
    }
}

func (val *XMPPStanzaIQRPCParamValue) parse_string(decoder *xml.Decoder, se xml.StartElement) (interface{}, error) {
    if t, err := decoder.Token(); err != nil {
        return nil, err
    } else {
        if cdata, ok := t.(xml.CharData); !ok {
            return nil, errors.New(fmt.Sprintf("Unexpected element. Expecting Chardata, got %+v", t))
        } else {
            return string(cdata), nil
        }
    }
}

func (val *XMPPStanzaIQRPCParamValue) parse_base64(decoder *xml.Decoder, se xml.StartElement) (interface{}, error) {
    if t, err := decoder.Token(); err != nil {
        return nil, err
    } else {
        if cdata, ok := t.(xml.CharData); !ok {
            return nil, errors.New(fmt.Sprintf("Unexpected element. Expecting Chardata, got %+v", t))
        } else {
            return base64.StdEncoding.DecodeString(string(cdata))
        }
    }
}

func (val *XMPPStanzaIQRPCParamValue) parse_double(decoder *xml.Decoder, se xml.StartElement) (interface{}, error) {
    if t, err := decoder.Token(); err != nil {
        return nil, err
    } else {
        if cdata, ok := t.(xml.CharData); !ok {
            return nil, errors.New(fmt.Sprintf("Unexpected element. Expecting Chardata, got %+v", t))
        } else {
            if xval, err := strconv.ParseFloat(string(cdata), 64); err != nil {
                return nil, err
            } else {
                return xval, nil
            }
        }
    }

}

func (val *XMPPStanzaIQRPCParamValue) parse_bool(decoder *xml.Decoder, se xml.StartElement) (interface{}, error) {
    if t, err := decoder.Token(); err != nil {
        return nil, err
    } else {
        if cdata, ok := t.(xml.CharData); !ok {
            return nil, errors.New(fmt.Sprintf("Unexpected element. Expecting Chardata, got %+v", t))
        } else {
            if xval, err := strconv.ParseBool(string(cdata)); err != nil {
                return nil, err
            } else {
                return xval, nil
            }
        }
    }
}

func (val *XMPPStanzaIQRPCParamValue) parse_datetime(decoder *xml.Decoder, se xml.StartElement) (interface{}, error) {
    if t, err := decoder.Token(); err != nil {
        return nil, err
    } else {
        if cdata, ok := t.(xml.CharData); !ok {
            return nil, errors.New(fmt.Sprintf("Unexpected element. Expecting Chardata, got %+v", t))
        } else {
            if xval, err := time.Parse(time_ISO8601_FORMAT, string(cdata)); err != nil {
                return nil, err
            } else {
                return xval, nil
            }
        }
    }
}

func (val *XMPPStanzaIQRPCParamValue) parse_array(decoder *xml.Decoder, se xml.StartElement) (interface{}, error) {
    type XMPPStanzaIQRPCParamValueArray struct {
        XMLName xml.Name                    `xml:"array"`
        Values  []XMPPStanzaIQRPCParamValue `xml:"data>value"`
    }

    arr := &XMPPStanzaIQRPCParamValueArray{}
    if err := decoder.DecodeElement(arr, &se); err != nil {
        return nil, err
    }

    var _arr []interface{}
    for _, v := range arr.Values {
        subval, err := v.Value()
        if err != nil {
            return nil, err
        }
        _arr = append(_arr, subval)
    }

    return _arr, nil

}

func (val *XMPPStanzaIQRPCParamValue) parse_struct(decoder *xml.Decoder, se xml.StartElement) (interface{}, error) {
    type XMPPStanzaIQRPCParamValueStruct struct {
        XMLName xml.Name `xml:"struct"`
        Members []struct {
            XMLName xml.Name                  `xml:"member"`
            Name    string                    `xml:"name"`
            Value   XMPPStanzaIQRPCParamValue `xml:""`
        }   `xml:"member"`
    }
    mmp := &XMPPStanzaIQRPCParamValueStruct{}
    if err := decoder.DecodeElement(mmp, &se); err != nil {
        return nil, err
    }
    _map := make(map[string]interface{})

    for _, mem := range mmp.Members {
        v, err := mem.Value.Value()
        if err != nil {
            return nil, err
        }
        _map[mem.Name] = v
    }
    return _map, nil
}
