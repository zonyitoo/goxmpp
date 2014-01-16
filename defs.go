package xmpp

import (
	"encoding/xml"
	"fmt"
)

const (
	XMLNS_JABBER_CLIENT    = "jabber:client"
	XMLNS_JABBER_IQ_ROSTER = "jabber:iq:roster"

	XMLNS_STREAM       = "http://etherx.jabber.org/streams"
	XMLNS_XMPP_TLS     = "urn:ietf:params:xml:ns:xmpp-tls"
	XMLNS_XMPP_SASL    = "urn:ietf:params:xml:ns:xmpp-sasl"
	XMLNS_XMPP_BIND    = "urn:ietf:params:xml:ns:xmpp-bind"
	XMLNS_XMPP_STANZAS = "urn:ietf:params:xml:ns:xmpp-stanzas"

	HANDLER_STREAM = "stream:stream"

	HANDLER_SASL_AUTH      = XMLNS_XMPP_SASL + ":auth"
	HANDLER_SASL_CHALLENGE = XMLNS_XMPP_SASL + ":challenge"
	HANDLER_SASL_RESPONSE  = XMLNS_XMPP_SASL + ":response"
	HANDLER_SASL_SUCCESS   = XMLNS_XMPP_SASL + ":success"
	HANDLER_SASL_FAILURE   = XMLNS_XMPP_SASL + ":failure"

	HANDLER_CLIENT_IQ = ":iq"
)

type XMPPStream struct {
	XMLName xml.Name `xml:"jabber:client stream:stream"`
	From    string   `xml:"from,attr"`
	To      string   `xml:"to,attr"`
	Id      string   `"xml:id,attr,omitempty"`
	Version string   `xml:"version,attr"`
	XmlLang string   `xml:"xml:lang,attr"`
}

type XMPPStreamFeatures struct {
	XMLName        xml.Name            `xml:"stream:features"`
	StartTLS       *XMPPStartTLS       `xml:",omitempty"`
	TLSProceed     *XMPPTLSProceed     `xml:",omitempty"`
	TLSFailure     *XMPPTLSFailure     `xml:",omitempty"`
	SASLMechanisms *XMPPSASLMechanisms `xml:",omitempty"`
	Bind           *XMPPBind           `xml:",omitempty"`
}

type XMPPStartTLS struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required bool     `xml:"required,omitempty"`
}

type XMPPTLSProceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

type XMPPTLSFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

type XMPPSASLMechanisms struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanisms []string `xml:"mechanism"`
}

type XMPPSASLAuth struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Data      string   `xml:",chardata"`
}

type XMPPSASLChallenge struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Data    string   `xml:",chardata"`
}

type XMPPSASLResponse struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
	Data    string   `xml:",chardata"`
}

type XMPPSASLSuccess struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
	Data    string   `xml:",chardata"`
}

type XMPPSASLFailure struct {
	XMLName       xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
	NotAuthorized bool     `xml:"not-authorized,omitempty"`
}

type XMPPBind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource,omitempty"`
	Jid      string   `xml:"jid,omitempty"`
}

type XMPPClientIQ struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	From    string   `xml:"from,attr"`
	To      string   `xml:"to,attr"`
	Id      string   `xml:"id,attr"`
	Type    string   `xml:"type,attr"`

	Error XMPPClientError
	Bind  XMPPBind
}

type XMPPClientError struct {
	XMLName xml.Name `xml:"jabber:client error"`
	Type    string   `xml:"type,attr"`
}

type XMPPClientMessage struct {
	XMLName xml.Name `xml:"jabber:client message"`
	From    string   `xml:"from,attr"`
	To      string   `xml:"from,attr"`
	Id      string   `xml:"id,attr"`
	Type    string   `xml:"type,attr"`
	XmlLang string   `xml:"xml:lang,attr"`
	Body    string   `xml:"body"`
}

type XMPPStreamError struct {
	XMLName   xml.Name    `xml:"stream:error"`
	ErrorType interface{} `xml:",omitempty"`
}

type XMPPStreamErrorDescriptiveText struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams text"`
	XmlLang string   `xml:"xml:lang,attr"`
	Text    string   `xml:",chardata"`
}

type XMPPStreamErrorNotWellFormed struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams not-well-formed"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorInvalidNamespace struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-namespace"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorHostUnknown struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams host-unknown"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorBadFormat struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams bad-format"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorBadNamespacePrefix struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams bad-namespace-prefix"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorConflict struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams conflict"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorConnectionTimeout struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams connection-timeout"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorHostGone struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams host-gone"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorImproperAddressing struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams improper-addressing"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorInternalServerError struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams internal-server-error"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorInvalidFrom struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-from"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorInvalidXML struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-xml"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorNotAuthorized struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams not-authorized"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorPolicyViolation struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams policy-violation"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorRemoteConnectionFailed struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams remote-connection-failed"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorReset struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams reset"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorResourceConstraint struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams resource-constraint"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamRestrictedXML struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams restricted-xml"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorSeeOtherHost struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams see-other-host"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorSystemShutdown struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams system-shutdown"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorUndefinedCondition struct {
	XMLName  xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams undefined-condition"`
	Text     *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
	AppError interface{}
}

type XMPPStreamErrorUnsupportedEncoding struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-encoding"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorUnsupportedFeature struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-feature"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorUnsupportedStanzaType struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-stanza-type"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorUnsupportedVersion struct {
	XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-version"`
	Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

const stream_response_begin_fmt = `<stream:stream
       from='%s'
       to='%s'
       version='%s'
       xml:lang='%s'
       id='%s'
       xmlns='%s'
       xmlns:stream='%s'>
`

const stream_end_fmt = `</stream:stream>`

func make_stream_begin(s *XMPPStream) string {
	return fmt.Sprintf(stream_response_begin_fmt,
		s.From, s.To, s.Version,
		s.XmlLang, s.Id, XMLNS_JABBER_CLIENT,
		XMLNS_STREAM)
}
