package xmpp

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"strings"
)

type XMPPSASLServerMechanismHandler interface {
	Auth(data string) (string, interface{}, error)
	Response(data string) (string, interface{}, error)
	Abort() (interface{}, error)
}

type XMPPSASLClientMechanismHandler interface {
	Begin() (interface{}, error)
	Challenge(data string) (interface{}, error)
	Success(data string) (interface{}, error)
	Failure(data string) (interface{}, error)
	Abort() (interface{}, error)
}

type DefaultSASLPlainHandler struct {
	authencator func(name, pwd string) bool
}

func NewDefaultSASLPlainHandler(f func(name, pwd string) bool) (*DefaultSASLPlainHandler, error) {
	if f == nil {
		return nil, errors.New("Authenticator should not be null")
	}
	return &DefaultSASLPlainHandler{
		authencator: f,
	}, nil
}

func (h *DefaultSASLPlainHandler) Auth(fdata string) (string, interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(fdata)
	if err != nil {
		return "", &XMPPSASLFailure{
			MalformedRequest: &XMPPSASLErrorMalformedRequest{},
		}, err
	}
	parts := strings.Split(string(data), "\x00")
	if len(parts) != 3 {
		return "", &XMPPSASLFailure{
			MalformedRequest: &XMPPSASLErrorMalformedRequest{},
		}, err
	}
	log.Printf("%v", parts)

	if h.authencator != nil && h.authencator(parts[1], parts[2]) {
		return parts[1], &XMPPSASLSuccess{
			Data: "=",
		}, nil
	} else {
		return "", &XMPPSASLFailure{
			NotAuthorized: &XMPPSASLErrorNotAuthorized{},
		}, errors.New("Authentication Fail")
	}

	return "", nil, nil
}

func (h *DefaultSASLPlainHandler) Response(data string) (string, interface{}, error) {
	return "", nil, nil
}

func (h *DefaultSASLPlainHandler) Abort() (interface{}, error) {
	return &XMPPSASLFailure{
		Aborted: &XMPPSASLErrorAborted{},
	}, errors.New("Aborted")
}

type SASLDigestChallenge struct {
	Realm      string
	Nonce      string
	QOPOptions string
	Stale      bool
	Maxbuf     int64
	Charset    string
	Algorithm  string
	CipherOpts string
	AuthParams string
}

func (this *SASLDigestChallenge) LoadFromBase64String(data string) error {
	d, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	for _, t := range strings.Split(string(d), ",") {
		p := strings.SplitN(t, "=", 2)
		if len(p) != 2 {
			return errors.New("Malformed String")
		}

		key := strings.Trim(p[0], " \"")
		val := strings.Trim(p[1], " \"")

		switch key {
		case "realm":
			this.Realm = val
		case "nonce":
			this.Nonce = val
		case "qop":
			this.QOPOptions = val
		case "stale":
			this.Stale, err = strconv.ParseBool(val)
			if err != nil {
				return err
			}
		case "maxbuf":
			this.Maxbuf, err = strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
		case "charset":
			this.Charset = val
		case "algorithm":
			this.Algorithm = val
		case "cipher":
			this.CipherOpts = val
		case "token":
			this.AuthParams = val
		}
	}

	return nil
}

func (this *SASLDigestChallenge) ToBase64String() string {
	s := []string{}
	if this.Realm != "" {
		s = append(s, "realm=\""+this.Realm+"\"")
	}
	if this.Nonce != "" {
		s = append(s, "nonce=\""+this.Nonce+"\"")
	}
	if this.QOPOptions != "" {
		s = append(s, "qop=\""+this.QOPOptions+"\"")
	}
	if this.Charset != "" {
		s = append(s, "charset="+this.Charset)
	}
	if this.Algorithm != "" {
		s = append(s, "algorithm="+this.Algorithm)
	}
	if this.Stale != true {
		s = append(s, "stale="+strconv.FormatBool(this.Stale))
	}
	if this.Maxbuf != 0 {
		s = append(s, "maxbuf="+strconv.FormatInt(this.Maxbuf, 10))
	}
	if this.CipherOpts != "" {
		s = append(s, "cipher=\""+this.CipherOpts+"\"")
	}
	if this.AuthParams != "" {
		s = append(s, "token=\""+this.AuthParams+"\"")
	}

	joined := strings.Join(s, ",")
	return base64.StdEncoding.EncodeToString([]byte(joined))
}

type SASLDigestResponse struct {
	UserName    string
	Cnonce      string
	NounceCount int64
	QOPOptions  string
	DigestURI   string
	Response    string
	Cipher      string
	Authzid     string
}

func (this *SASLDigestResponse) LoadFromBase64String(str string) error {
	ddat, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	for _, t := range strings.Split(string(ddat), ",") {
		p := strings.SplitN(t, "=", 2)

		if len(p) != 2 {
			return errors.New("Malformed String")
		}

		key := strings.Trim(p[0], " \"")
		val := strings.Trim(p[1], " \"")

		switch key {
		case "username":
			this.UserName = val
		case "cnonce":
			this.Cnonce = val
		case "nc":
			this.NounceCount, err = strconv.ParseInt(val, 16, 64)
			if err != nil {
				return err
			}
		case "qop":
			this.QOPOptions = val
		case "digest-uri":
			this.DigestURI = val
		case "Response":
			this.Response = val
		case "cipher":
			this.Cipher = val
		case "authzid":
			this.Authzid = val
		}
	}

	return nil
}

type DefaultSASLDigestMD5Handler struct {
	authencator func(name, pwd string) bool
}

func NewDefaultSASLDigestMD5Handler(f func(name, pwd string) bool) (*DefaultSASLDigestMD5Handler, error) {
    if f == nil {
        return nil, errors.New("Null Authenticator")
    }

    return &DefaultSASLDigestMD5Handler{
        authencator: f,
    }, nil
}

func (h *DefaultSASLDigestMD5Handler) Auth(_ string) (string, interface{}, error) {
	// Construct Challenge
	cdat := &SASLDigestChallenge{
		Realm:      "abc",
		Nonce:      generate_random_id(),
		QOPOptions: "auth",
		Charset:    "utf-8",
		Algorithm:  "md5-sess",
	}

	challenge := &XMPPSASLChallenge{
		Data: cdat.ToBase64String(),
	}
	return "", challenge, nil
}

func (h *DefaultSASLDigestMD5Handler) Response(data string) (string, interface{}, error) {

    rdat := SASLDigestResponse{}
    if err := rdat.LoadFromBase64String(data); err != nil {
        return "", nil, err
    }

    if h.authencator(rdat.UserName, rdat.Response) {
        return rdat.UserName, &XMPPSASLSuccess{}, nil
    } else {
        return rdat.UserName, &XMPPSASLFailure{
            NotAuthorized: &XMPPSASLErrorNotAuthorized{},
        }, errors.New("Not Authorized")
    }
}

func (h *DefaultSASLDigestMD5Handler) Abort() (interface{}, error) {
	return &XMPPSASLFailure{
		Aborted: &XMPPSASLErrorAborted{},
	}, errors.New("Aborted")
}
