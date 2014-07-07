package xmpp

import (
    "encoding/xml"
    "fmt"
)

const (
    XMLNS_JABBER_CLIENT    = "jabber:client"
    XMLNS_JABBER_IQ_ROSTER = "jabber:iq:roster"
    XMLNS_JABBER_SERVER    = "jabber:server"

    XMLNS_STREAM       = "http://etherx.jabber.org/streams"
    XMLNS_XMPP_TLS     = "urn:ietf:params:xml:ns:xmpp-tls"
    XMLNS_XMPP_SASL    = "urn:ietf:params:xml:ns:xmpp-sasl"
    XMLNS_XMPP_BIND    = "urn:ietf:params:xml:ns:xmpp-bind"
    XMLNS_XMPP_STANZAS = "urn:ietf:params:xml:ns:xmpp-stanzas"
)

var TAG_STREAM xml.Name = xml.Name{Space: XMLNS_STREAM, Local: "stream"}
var TAG_STREAM_FEATURES xml.Name = xml.Name{Space: XMLNS_STREAM, Local: "features"}
var TAG_TLS_START xml.Name = xml.Name{Space: XMLNS_XMPP_TLS, Local: "starttls"}
var TAG_TLS_PROCEED xml.Name = xml.Name{Space: XMLNS_XMPP_TLS, Local: "proceed"}
var TAG_TLS_FAILURE xml.Name = xml.Name{Space: XMLNS_XMPP_TLS, Local: "failure"}
var TAG_TLS_ABORT xml.Name = xml.Name{Space: XMLNS_XMPP_TLS, Local: "abort"}
var TAG_SASL_AUTH xml.Name = xml.Name{Space: XMLNS_XMPP_SASL, Local: "auth"}
var TAG_SASL_CHALLENGE xml.Name = xml.Name{Space: XMLNS_XMPP_SASL, Local: "challenge"}
var TAG_SASL_RESPONSE xml.Name = xml.Name{Space: XMLNS_XMPP_SASL, Local: "response"}
var TAG_SASL_SUCCESS xml.Name = xml.Name{Space: XMLNS_XMPP_SASL, Local: "success"}
var TAG_SASL_FAILURE xml.Name = xml.Name{Space: XMLNS_XMPP_SASL, Local: "failure"}
var TAG_CLIENT_IQ xml.Name = xml.Name{Space: XMLNS_JABBER_CLIENT, Local: "iq"}
var TAG_CLIENT_PRESENCE xml.Name = xml.Name{Space: XMLNS_JABBER_CLIENT, Local: "presence"}
var TAG_CLIENT_MESSAGE xml.Name = xml.Name{Space: XMLNS_JABBER_CLIENT, Local: "message"}

// RFC6120 Section 4
type XMPPStream struct {
    XMLName xml.Name `xml:"http://etherx.jabber.org/streams stream"`
    From    string   `xml:"from,attr"`
    To      string   `xml:"to,attr"`
    Id      string   `xml:"id,attr,omitempty"`
    Version string   `xml:"version,attr"`
    XmlLang string   `xml:"xml:lang,attr"`
}

type XMPPStreamEnd struct {
    XMLName xml.Name `xml:"http://etherx.jabber.org/streams stream"`
}

type XMPPStreamFeatures struct {
    XMLName        xml.Name            `xml:"http://etherx.jabber.org/streams features"`
    StartTLS       *XMPPStartTLS       `xml:",omitempty"`
    SASLMechanisms *XMPPSASLMechanisms `xml:",omitempty"`
    Bind           *XMPPBind           `xml:",omitempty"`
}

type XMPPRequired struct {
    XMLName xml.Name `xml:"required"`
}

// RFC6120 Section 5.4.2.1
//
// In order to begin the STARTTLS negotiation, the initiating entity issues the
// STARTTLS command (i.e., a <starttls/> element qualified by the
// 'urn:ietf:params:xml:ns:xmpp-tls' namespace) to instruct the receiving entity
// that it wishes to begin a STARTTLS negotiation to secure the stream.
type XMPPStartTLS struct {
    XMLName  xml.Name      `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
    Required *XMPPRequired `xml:"required,omitempty"`
}

// RFC6120 Section 5.4.2.3
//
// If the proceed case occurs, the receiving entity MUST return a <proceed/>
// element qualified by the 'urn:ietf:params:xml:ns:xmpp-tls' namespace.
type XMPPTLSProceed struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

// RFC6120 Section 5.4.2.2
//
// If the failure case occurs, the receiving entity MUST return a <failure/> element
// qualified by the 'urn:ietf:params:xml:ns:xmpp-tls' namespace, close the XML stream,
// and terminate the underlying TCP connection.
type XMPPTLSFailure struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

type XMPPTLSAbort struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls abort"`
}

type XMPPSASLMechanisms struct {
    XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
    Mechanisms []string `xml:"mechanism"`
}

// RFC6120 Section 6.4.2
//
// In order to begin the SASL negotiation, the initiating entity sends an <auth/>
// element qualified by the 'urn:ietf:params:xml:ns:xmpp-sasl' namespace and includes
// an appropriate value for the 'mechanism' attribute, thus starting the handshake for
// that particular authentication mechanism. This element MAY contain XML character
// data (in SASL terminology, the "initial response") if the mechanism supports or
// requires it. If the initiating entity needs to send a zero-length initial response,
// it MUST transmit the response as a single equals sign character ("="), which
// indicates that the response is present but contains no data.
type XMPPSASLAuth struct {
    XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
    Mechanism string   `xml:"mechanism,attr"`
    Data      string   `xml:",chardata"`
}

// RFC6120 Section 6.4.3
//
// If necessary, the receiving entity challenges the initiating entity by sending a
// <challenge/> element qualified by the 'urn:ietf:params:xml:ns:xmpp-sasl' namespace;
// this element MAY contain XML character data (which MUST be generated in accordance
// with the definition of the SASL mechanism chosen by the initiating entity).
//
// The initiating entity responds to the challenge by sending a <response/> element
// qualified by the 'urn:ietf:params:xml:ns:xmpp-sasl' namespace; this element MAY
// contain XML character data (which MUST be generated in accordance with the definition
// of the SASL mechanism chosen by the initiating entity).
//
// If necessary, the receiving entity sends more challenges and the initiating entity
// sends more responses.
//
// This series of challenge/response pairs continues until one of three things happens:
//
// * The initiating entity aborts the handshake for this authentication mechanism.
// * The receiving entity reports failure of the handshake.
// * The receiving entity reports success of the handshake.
// * These scenarios are described in the following sections.
type XMPPSASLChallenge struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
    Data    string   `xml:",chardata"`
}

type XMPPSASLResponse struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
    Data    string   `xml:",chardata"`
}

// RFC6120 Section 6.4.6
//
// Before considering the SASL handshake to be a success, if the initiating entity provided
// a 'from' attribute on an initial stream header whose confidentiality and integrity were
// protected via TLS or an equivalent security layer (such as the SASL GSSAPI mechanism)
// then the receiving entity SHOULD correlate the authentication identity resulting from
// the SASL negotiation with that 'from' address; if the two identities do not match then
// the receiving entity SHOULD terminate the connection attempt (however, the receiving
// entity might have legitimate reasons not to terminate the connection attempt, for example,
// because it has overridden a connecting client's address to correct the JID format or
// assign a JID based on information presented in an end-user certificate).
//
// The receiving entity reports success of the handshake by sending a <success/> element
// qualified by the 'urn:ietf:params:xml:ns:xmpp-sasl' namespace; this element MAY contain
// XML character data (in SASL terminology, "additional data with success") if the chosen
// SASL mechanism supports or requires it. If the receiving entity needs to send additional
// data of zero length, it MUST transmit the data as a single equals sign character ("=").
type XMPPSASLSuccess struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
    Data    string   `xml:",chardata"`
}

// RFC6120 Section 6.4.4
//
// The initiating entity aborts the handshake for this authentication mechanism by sending
// an <abort/> element qualified by the 'urn:ietf:params:xml:ns:xmpp-sasl' namespace.
type XMPPSASLAbort struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl abort"`
}

// RFC6120 Section 6.4.5
//
// The receiving entity reports failure of the handshake for this authentication mechanism
// by sending a <failure/> element qualified by the 'urn:ietf:params:xml:ns:xmpp-sasl'
// namespace (the particular cause of failure MUST be communicated in an appropriate child
// element of the <failure/> element as defined under Section 6.5).
type XMPPSASLFailure struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`

    Aborted              *XMPPSASLErrorAborted              `xml:",omitempty"`
    AccountDisabled      *XMPPSASLErrorAccountDisabled      `xml:",omitempty"`
    CredentialsExpired   *XMPPSASLErrorCredentialsExpired   `xml:",omitempty"`
    EncryptionRequired   *XMPPSASLErrorEncryptionRequired   `xml:",omitempty"`
    IncorrectEncoding    *XMPPSASLErrorIncorrectEncoding    `xml:",omitempty"`
    InvalidAuthzid       *XMPPSASLErrorInvalidAuthzid       `xml:",omitempty"`
    InvalidMechanism     *XMPPSASLErrorInvalidMechanism     `xml:",omitempty"`
    MalformedRequest     *XMPPSASLErrorMalformedRequest     `xml:",omitempty"`
    MechanismTooWeak     *XMPPSASLErrorMechanismTooWeak     `xml:",omitempty"`
    NotAuthorized        *XMPPSASLErrorNotAuthorized        `xml:",omitempty"`
    TemporaryAuthFailure *XMPPSASLErrorTemporaryAuthFailure `xml:",omitempty"`
}

// RFC6120 Section 6.5.1
//
// The receiving entity acknowledges that the authentication handshake has been aborted
// by the initiating entity; sent in reply to the <abort/> element.
type XMPPSASLErrorAborted struct {
    XMLName xml.Name                      `xml:"aborted"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.2
//
// The account of the initiating entity has been temporarily disabled; sent in reply to
// an <auth/> element (with or without initial response data) or a <response/> element.
type XMPPSASLErrorAccountDisabled struct {
    XMLName xml.Name                      `xml:"account-disabled"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.3
//
// The authentication failed because the initiating entity provided credentials that
// have expired; sent in reply to a <response/> element or an <auth/> element with
// initial response data.
type XMPPSASLErrorCredentialsExpired struct {
    XMLName xml.Name                      `xml:"credentials-expired"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.4
//
// The mechanism requested by the initiating entity cannot be used unless the
// confidentiality and integrity of the underlying stream are protected (typically
// via TLS); sent in reply to an <auth/> element (with or without initial response data).
type XMPPSASLErrorEncryptionRequired struct {
    XMLName xml.Name                      `xml:"encryption-required"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.5
//
// The data provided by the initiating entity could not be processed because the base
// 64 encoding is incorrect (e.g., because the encoding does not adhere to the definition
// in Section 4 of [BASE64]); sent in reply to a <response/> element or an <auth/> element
// with initial response data.
type XMPPSASLErrorIncorrectEncoding struct {
    XMLName xml.Name                      `xml:"incorrect-encoding"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.6
//
// The authzid provided by the initiating entity is invalid, either because it is
// incorrectly formatted or because the initiating entity does not have permissions to
// authorize that ID; sent in reply to a <response/> element or an <auth/> element with
// initial response data.
type XMPPSASLErrorInvalidAuthzid struct {
    XMLName xml.Name                      `xml:"invalid-authzid"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.7
//
// The initiating entity did not specify a mechanism, or requested a mechanism that
// is not supported by the receiving entity; sent in reply to an <auth/> element.
type XMPPSASLErrorInvalidMechanism struct {
    XMLName xml.Name                      `xml:"invalid-mechanism"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.8
//
// The request is malformed (e.g., the <auth/> element includes initial response data
// but the mechanism does not allow that, or the data sent violates the syntax for
// the specified SASL mechanism); sent in reply to an <abort/>, <auth/>, <challenge/>,
// or <response/> element.
type XMPPSASLErrorMalformedRequest struct {
    XMLName xml.Name                      `xml:"malformed-request"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.9
//
// The mechanism requested by the initiating entity is weaker than server policy
// permits for that initiating entity; sent in reply to an <auth/> element (with or
// without initial response data).
type XMPPSASLErrorMechanismTooWeak struct {
    XMLName xml.Name                      `xml:"mechanism-too-weak"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.10
//
// The authentication failed because the initiating entity did not provide proper
// credentials, or because some generic authentication failure has occurred but the
// receiving entity does not wish to disclose specific information about the cause of
// the failure; sent in reply to a <response/> element or an <auth/> element with
// initial response data.
type XMPPSASLErrorNotAuthorized struct {
    XMLName xml.Name                      `xml:"not-authorized"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 6.5.11
//
// The authentication failed because of a temporary error condition within the
// receiving entity, and it is advisable for the initiating entity to try again later;
// sent in reply to an <auth/> element or a <response/> element.
type XMPPSASLErrorTemporaryAuthFailure struct {
    XMLName xml.Name                      `xml:"temporary-auth-failure"`
    Text    *XMPPSASLErrorDescriptiveText `xml:",omitempty"`
}

type XMPPSASLErrorDescriptiveText struct {
    XMLName xml.Name `xml:"text"`
    XmlLang string   `xml:"xml:lang,attr"`
    Text    string   `xml:",chardata"`
}

// RFC6120 Section 7
type XMPPBind struct {
    XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
    Resource string   `xml:"resource,omitempty"`
    Jid      string   `xml:"jid,omitempty"`
}

// RFC6120 Section 8.2.3
//
// Info/Query, or IQ, is a "request-response" mechanism, similar in some ways to the
// Hypertext Transfer Protocol [HTTP]. The semantics of IQ enable an entity to make a
// request of, and receive a response from, another entity. The data content of the
// request and response is defined by the schema or other structural definition associated
// with the XML namespace that qualifies the direct child element of the IQ element
// (see Section 8.4), and the interaction is tracked by the requesting entity through
// use of the 'id' attribute. Thus, IQ interactions follow a common pattern of structured
// data exchange such as get/result or set/result (although an error can be returned in
// reply to a request if appropriate)
type XMPPClientIQ struct {
    XMLName xml.Name `xml:"jabber:client iq"`
    From    string   `xml:"from,attr,omitempty"`
    To      string   `xml:"to,attr,omitempty"`
    Id      string   `xml:"id,attr"`
    Type    string   `xml:"type,attr"`

    Error  *XMPPStanzaError    `xml:",omitempty"`
    Bind   *XMPPBind           `xml:",omitempty"`
    Roster *XMPPClientIQRoster `xml:",omitempty"`
}

type XMPPClientIQRoster struct {
    XMLName xml.Name `xml:"jabber:iq:roster query"`
}

// RFC6120 Section 8.2.1
//
// The <message/> stanza is a "push" mechanism whereby one entity pushes information to
// another entity, similar to the communications that occur in a system such as email.
// All message stanzas will possess a 'to' attribute that specifies the intended recipient
// of the message (see Section 8.1.1 and Section 10.3), unless the message is being sent
// to the bare JID of a connected client's account. Upon receiving a message stanza with a
// 'to' address, a server SHOULD attempt to route or deliver it to the intended recipient
// (see Section 10 for general routing and delivery rules related to XML stanzas).
type XMPPClientMessage struct {
    XMLName xml.Name `xml:"jabber:client message"`
    From    string   `xml:"from,attr"`
    To      string   `xml:"to,attr"`
    Type    string   `xml:"type,attr"`
    XmlLang string   `xml:"xml:lang,attr"`
    Body    string   `xml:"body"`

    Error *XMPPStanzaError `xml:",omitempty"`
}

// RFC6120 Section 8.2.2
//
// The <presence/> stanza is a specialized "broadcast" or "publish-subscribe" mechanism,
// whereby multiple entities receive information (in this case, network availability information)
// about an entity to which they have subscribed. In general, a publishing client SHOULD send a
// presence stanza with no 'to' attribute, in which case the server to which the client is
// connected will broadcast that stanza to all subscribed entities. However, a publishing client
// MAY also send a presence stanza with a 'to' attribute, in which case the server will route or
// deliver that stanza to the intended recipient. Although the <presence/> stanza is most often
// used by XMPP clients, it can also be used by servers, add-on services, and any other kind of
// XMPP entity. See Section 10 for general routing and delivery rules related to XML stanzas,
// and [XMPP‑IM] for rules specific to presence applications.
type XMPPClientPresence struct {
    XMLName xml.Name `xml:"jabber:client presence"`
    XmlLang string   `xml:"xml:lang,attr"`
    Show    string   `xml:"show",omitempty`
    Status  string   `xml:"status",omitempty`

    Error *XMPPStanzaError `xml:",omitempty"`
}

// RFC6120 Section 4.9
type XMPPStreamError struct {
    XMLName xml.Name `xml:"stream:error"`

    NotWellFormed          *XMPPStreamErrorNotWellFormed          `xml:",omitempty"`
    InvalidNamespace       *XMPPStreamErrorInvalidNamespace       `xml:",omitempty"`
    HostUnknown            *XMPPStreamErrorHostUnknown            `xml:",omitempty"`
    BadFormat              *XMPPStreamErrorBadFormat              `xml:",omitempty"`
    BadNamespacePrefix     *XMPPStreamErrorBadNamespacePrefix     `xml:",omitempty"`
    Conflict               *XMPPStreamErrorConflict               `xml:",omitempty"`
    ConnectionTimeout      *XMPPStreamErrorConnectionTimeout      `xml:",omitempty"`
    HostGone               *XMPPStreamErrorHostGone               `xml:",omitempty"`
    ImproperAddressing     *XMPPStreamErrorImproperAddressing     `xml:",omitempty"`
    InternalServerError    *XMPPStreamErrorInternalServerError    `xml:",omitempty"`
    InvalidFrom            *XMPPStreamErrorInvalidFrom            `xml:",omitempty"`
    InvalidXML             *XMPPStreamErrorInvalidXML             `xml:",omitempty"`
    NotAuthorized          *XMPPStreamErrorNotAuthorized          `xml:",omitempty"`
    PolicyViolation        *XMPPStreamErrorPolicyViolation        `xml:",omitempty"`
    RemoteConnectionFailed *XMPPStreamErrorRemoteConnectionFailed `xml:",omitempty"`
    Reset                  *XMPPStreamErrorReset                  `xml:",omitempty"`
    ResourceConstraint     *XMPPStreamErrorResourceConstraint     `xml:",omitempty"`
    RestrictedXML          *XMPPStreamErrorRestrictedXML          `xml:",omitempty"`
    SeeOtherHost           *XMPPStreamErrorSeeOtherHost           `xml:",omitempty"`
    SystemShutdown         *XMPPStreamErrorSystemShutdown         `xml:",omitempty"`
    UndefinedCondition     *XMPPStreamErrorUndefinedCondition     `xml:",omitempty"`
    UnsupportedEncoding    *XMPPStreamErrorUnsupportedEncoding    `xml:",omitempty"`
    UnsupportedFeature     *XMPPStreamErrorUnsupportedFeature     `xml:",omitempty"`
    UnsupportedStanzaType  *XMPPStreamErrorUnsupportedStanzaType  `xml:",omitempty"`
    UnsupportedVersion     *XMPPStreamErrorUnsupportedVersion     `xml:",omitempty"`

    // RFC6120 Section 4.9.4
    //
    // As noted, an application MAY provide application-specific stream error information
    // by including a properly namespaced child in the error element. The application-specific
    // element SHOULD supplement or further qualify a defined element. Thus, the <error/>
    // element will contain two or three child elements.
    ApplicationSpecificConditions string `xml:",omitempty"`
}

type XMPPStreamErrorDescriptiveText struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams text"`
    XmlLang string   `xml:"xml:lang,attr"`
    Text    string   `xml:",chardata"`
}

// RFC6120 Section 4.9.3.1
// The entity has sent XML that cannot be processed.
//
//    C: <message>
//         <body>No closing tag!
//       </message>
//
//    S: <stream:error>
//         <bad-format
//             xmlns='urn:ietf:params:xml:ns:xmpp-streams'/>
//       </stream:error>
//       </stream:stream>
//
// RFC6120 Section 4.9.3.13
//
// The initiating entity has sent XML that violates the well-formedness rules of [XML] or [XML‑NAMES].
type XMPPStreamErrorNotWellFormed struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams not-well-formed"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.10
//
// The stream namespace name is something other than "http://etherx.jabber.org/streams"
// (see Section 11.2) or the content namespace declared as the default namespace is not
// supported (e.g., something other than "jabber:client" or "jabber:server").
type XMPPStreamErrorInvalidNamespace struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-namespace"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.6
//
// The value of the 'to' attribute provided in the initial stream header does not correspond
// to an FQDN that is serviced by the receiving entity.
type XMPPStreamErrorHostUnknown struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams host-unknown"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStreamErrorBadFormat struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams bad-format"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.2
//
// The entity has sent a namespace prefix that is unsupported,
// or has sent no namespace prefix on an element that needs such a prefix
type XMPPStreamErrorBadNamespacePrefix struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams bad-namespace-prefix"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.3
//
// The server either (1) is closing the existing stream for this entity
// because a new stream has been initiated that conflicts with the existing stream,
// or (2) is refusing a new stream for this entity because allowing the new stream
// would conflict with an existing stream
type XMPPStreamErrorConflict struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams conflict"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.4
//
// One party is closing the stream because it has reason to believe that the other
// party has permanently lost the ability to communicate over the stream. The lack of
// ability to communicate can be discovered using various methods, such as whitespace
// keepalives as specified under Section 4.4, XMPP-level pings as defined in [XEP‑0199],
// and XMPP Stream Management as defined in [XEP‑0198].
type XMPPStreamErrorConnectionTimeout struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams connection-timeout"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.5
//
// The value of the 'to' attribute provided in the initial stream header corresponds to
// an FQDN that is no longer serviced by the receiving entity.
type XMPPStreamErrorHostGone struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams host-gone"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.7
//
// A stanza sent between two servers lacks a 'to' or 'from' attribute,
// the 'from' or 'to' attribute has no value, or the value violates the rules for XMPP
// addresses [XMPP‑ADDR].
type XMPPStreamErrorImproperAddressing struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams improper-addressing"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.8
//
// The server has experienced a misconfiguration or other internal error that prevents
// it from servicing the stream.
type XMPPStreamErrorInternalServerError struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams internal-server-error"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.9
//
// The data provided in a 'from' attribute does not match an authorized JID or validated
// domain as negotiated (1) between two servers using SASL or Server Dialback, or (2)
// between a client and a server via SASL authentication and resource binding.
type XMPPStreamErrorInvalidFrom struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-from"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.11
//
// The entity has sent invalid XML over the stream to a server that performs validation
// (see Section 11.4).
type XMPPStreamErrorInvalidXML struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-xml"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.12
//
// The entity has attempted to send XML stanzas or other outbound data before the stream
// has been authenticated, or otherwise is not authorized to perform an action related
// to stream negotiation; the receiving entity MUST NOT process the offending data before
// sending the stream error.
type XMPPStreamErrorNotAuthorized struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams not-authorized"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.14
//
// The entity has violated some local service policy (e.g., a stanza exceeds a configured
// size limit); the server MAY choose to specify the policy in the <text/> element or in
// an application-specific condition element.
type XMPPStreamErrorPolicyViolation struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams policy-violation"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.15
//
// The server is unable to properly connect to a remote entity that is needed for
// authentication or authorization (e.g., in certain scenarios related to Server Dialback
// [XEP‑0220]); this condition is not to be used when the cause of the error is within the
// administrative domain of the XMPP service provider, in which case the
// <internal-server-error/> condition is more appropriate.
type XMPPStreamErrorRemoteConnectionFailed struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams remote-connection-failed"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.16
//
// The server is closing the stream because it has new (typically security-critical)
// features to offer, because the keys or certificates used to establish a secure context
// for the stream have expired or have been revoked during the life of the stream
// (Section 13.7.2.3), because the TLS sequence number has wrapped (Section 5.3.5), etc.
// The reset applies to the stream and to any security context established for that stream
// (e.g., via TLS and SASL), which means that encryption and authentication need to be
// negotiated again for the new stream (e.g., TLS session resumption cannot be used).
type XMPPStreamErrorReset struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams reset"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.17
//
// The server lacks the system resources necessary to service the stream.
type XMPPStreamErrorResourceConstraint struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams resource-constraint"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.18
//
// The entity has attempted to send restricted XML features such as a comment, processing
// instruction, DTD subset, or XML entity reference (see Section 11.1).
type XMPPStreamErrorRestrictedXML struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams restricted-xml"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.19
//
// The server will not provide service to the initiating entity but is redirecting traffic
// to another host under the administrative control of the same service provider. The XML
// character data of the <see-other-host/> element returned by the server MUST specify
// the alternate FQDN or IP address at which to connect, which MUST be a valid domainpart
// or a domainpart plus port number (separated by the ':' character in the form "domainpart:port").
// If the domainpart is the same as the source domain, derived domain, or resolved IPv4 or
// IPv6 address to which the initiating entity originally connected (differing only by the
// port number), then the initiating entity SHOULD simply attempt to reconnect at that address.
// (The format of an IPv6 address MUST follow [IPv6‑ADDR], which includes the enclosing the
// IPv6 address in square brackets '[' and ']' as originally defined by [URI].) Otherwise,
// the initiating entity MUST resolve the FQDN specified in the <see-other-host/> element as
// described under Section 3.2.
type XMPPStreamErrorSeeOtherHost struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams see-other-host"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.20
//
// The server is being shut down and all active streams are being closed.
type XMPPStreamErrorSystemShutdown struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams system-shutdown"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.21
//
// The error condition is not one of those defined by the other conditions in this list;
// this error condition SHOULD NOT be used except in conjunction with an application-specific
// condition.
type XMPPStreamErrorUndefinedCondition struct {
    XMLName  xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams undefined-condition"`
    Text     *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
    AppError interface{}
}

// RFC6120 Section 4.9.3.22
//
// The initiating entity has encoded the stream in an encoding that is not supported by the
// server (see Section 11.6) or has otherwise improperly encoded the stream (e.g., by violating
// the rules of the [UTF‑8] encoding).
type XMPPStreamErrorUnsupportedEncoding struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-encoding"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.23
//
// The receiving entity has advertised a mandatory-to-negotiate stream feature that the initiating
// entity does not support, and has offered no other mandatory-to-negotiate feature alongside the
// unsupported feature.
type XMPPStreamErrorUnsupportedFeature struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-feature"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.24
//
// The initiating entity has sent a first-level child of the stream that is not supported
// by the server, either because the receiving entity does not understand the namespace
// or because the receiving entity does not understand the element name for the applicable
// namespace (which might be the content namespace declared as the default namespace).
type XMPPStreamErrorUnsupportedStanzaType struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-stanza-type"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 4.9.3.25
//
// The 'version' attribute provided by the initiating entity in the stream header
// specifies a version of XMPP that is not supported by the server.
type XMPPStreamErrorUnsupportedVersion struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-version"`
    Text    *XMPPStreamErrorDescriptiveText `xml:",omitempty"`
}

type XMPPStanzaErrorDescriptiveText struct {
    XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas text"`
    XmlLang string   `xml:"xml:lang,attr"`
}

// RFC6120 Section 8.3
//
// Stanza-related errors are handled in a manner similar to stream errors. Unlike stream errors,
// stanza errors are recoverable; therefore, they do not result in termination of the XML stream
// and underlying TCP connection. Instead, the entity that discovers the error condition returns
// an error stanza, which is a stanza that:
//
// * is of the same kind (message, presence, or IQ) as the generated stanza that triggered the error
// * has a 'type' attribute set to a value of "error"
// * typically swaps the 'from' and 'to' addresses of the generated stanza
// * mirrors the 'id' attribute (if any) of the generated stanza that triggered the error
// * contains an <error/> child element that specifies the error condition and therefore provides
//   a hint regarding actions that the sender might be able to take in an effort to remedy the
//   error (however, it is not always possible to remedy the error)
type XMPPStanzaError struct {
    XMLName xml.Name `xml:"error"`
    Type    string   `xml:"type,attr"`

    BadRequest            *XMPPStanzaErrorBadRequest            `xml:",omitempty"`
    Conflict              *XMPPStanzaErrorConflict              `xml:",omitempty"`
    FeatureNotImplemented *XMPPStanzaErrorFeatureNotImplemented `xml:",omitempty"`
    Forbidden             *XMPPStanzaErrorForbidden             `xml:",omitempty"`
    Gone                  *XMPPStanzaErrorGone                  `xml:",omitempty"`
    InternalServerError   *XMPPStanzaErrorInternalServerError   `xml:",omitempty"`
    ItemNotFound          *XMPPStanzaErrorItemNotFound          `xml:",omitempty"`
    JIDMalformed          *XMPPStanzaErrorJIDMalformed          `xml:",omitempty"`
    NotAcceptable         *XMPPStanzaErrorNotAcceptable         `xml:",omitempty"`
    NotAllowed            *XMPPStanzaErrorNotAllowed            `xml:",omitempty"`
    NotAuthorized         *XMPPStanzaErrorNotAuthorized         `xml:",omitempty"`
    PolicyViolation       *XMPPStanzaErrorPolicyVoilation       `xml:",omitempty"`
    RecipientUnavailable  *XMPPStanzaErrorRecipientUnavailable  `xml:",omitempty"`
    Redirect              *XMPPStanzaErrorRedirect              `xml:",omitempty"`
    RegistrationRequired  *XMPPStanzaErrorRegistrationRequired  `xml:",omitempty"`
    RemoteServerTimeout   *XMPPStanzaErrorRemoteServerTimeout   `xml:",omitempty"`
    RemoteServerNotFound  *XMPPStanzaErrorRemoteServerNotFound  `xml:",omitempty"`
    ResourceConstraint    *XMPPStanzaErrorResourceConstraint    `xml:",omitempty"`
    ServiceUnavailable    *XMPPStanzaErrorServiceUnavailable    `xml:",omitempty"`
    SubscriptionRequired  *XMPPStanzaErrorSubscriptionRequired  `xml:",omitempty"`
    UndefinedCondition    *XMPPStanzaErrorUndefinedCondition    `xml:",omitempty"`
    UnexpectedRequest     *XMPPStanzaErrorUnexpectedRequest     `xml:",omitempty"`

    ApplicationSpecificConditions string `xml:",any,innerxml,omitempty"`
}

// RFC6120 Section 8.3.3.1
//
// The sender has sent a stanza containing XML that does not conform to the appropriate schema
// or that cannot be processed (e.g., an IQ stanza that includes an unrecognized value of the
// 'type' attribute, or an element that is qualified by a recognized namespace but that violates
// the defined syntax for the element); the associated error type SHOULD be "modify".
type XMPPStanzaErrorBadRequest struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas bad-request"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.2
//
// Access cannot be granted because an existing resource exists with the same name or address;
// the associated error type SHOULD be "cancel".
type XMPPStanzaErrorConflict struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas conflict"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.3
//
// The feature represented in the XML stanza is not implemented by the intended recipient or
// an intermediate server and therefore the stanza cannot be processed (e.g., the entity
// understands the namespace but does not recognize the element name); the associated error
// type SHOULD be "cancel" or "modify".
type XMPPStanzaErrorFeatureNotImplemented struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas feature-not-implemented"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.4
//
// The requesting entity does not possess the necessary permissions to perform an action
// that only certain authorized roles or individuals are allowed to complete (i.e., it
// typically relates to authorization rather than authentication); the associated error
// type SHOULD be "auth".
type XMPPStanzaErrorForbidden struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas forbidden"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.5
//
// The recipient or server can no longer be contacted at this address, typically on a
// permanent basis (as opposed to the <redirect/> error condition, which is used for
// temporary addressing failures); the associated error type SHOULD be "cancel" and the
// error stanza SHOULD include a new address (if available) as the XML character data of
// the <gone/> element (which MUST be a Uniform Resource Identifier [URI] or
// Internationalized Resource Identifier [IRI] at which the entity can be contacted,
// typically an XMPP IRI as specified in [XMPP‑URI]).
type XMPPStanzaErrorGone struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas gone"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.6
//
// The server has experienced a misconfiguration or other internal error that prevents
// it from processing the stanza; the associated error type SHOULD be "cancel".
type XMPPStanzaErrorInternalServerError struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas internal-server-error"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.7
//
// The addressed JID or item requested cannot be found; the associated error type
// SHOULD be "cancel".
type XMPPStanzaErrorItemNotFound struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas item-not-found"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.8
//
// The sending entity has provided (e.g., during resource binding) or communicated
// (e.g., in the 'to' address of a stanza) an XMPP address or aspect thereof that
// violates the rules defined in [XMPP‑ADDR]; the associated error type SHOULD be
// "modify".
type XMPPStanzaErrorJIDMalformed struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas jid-malformed"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.9
//
// The recipient or server understands the request but cannot process it because
// the request does not meet criteria defined by the recipient or server (e.g.,
// a request to subscribe to information that does not simultaneously include
// configuration parameters needed by the recipient); the associated error type
// SHOULD be "modify".
type XMPPStanzaErrorNotAcceptable struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-acceptable"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.10
//
// The recipient or server does not allow any entity to perform the action (e.g.,
// sending to entities at a blacklisted domain); the associated error type SHOULD
// be "cancel".
type XMPPStanzaErrorNotAllowed struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-allowed"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.11
//
// The sender needs to provide credentials before being allowed to perform the
// action, or has provided improper credentials (the name "not-authorized", which
// was borrowed from the "401 Unauthorized" error of [HTTP], might lead the reader
// to think that this condition relates to authorization, but instead it is
// typically used in relation to authentication); the associated error type SHOULD
// be "auth".
type XMPPStanzaErrorNotAuthorized struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-authorized"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.12
// The entity has violated some local service policy (e.g., a message contains words
// that are prohibited by the service) and the server MAY choose to specify the
// policy in the <text/> element or in an application-specific condition element;
// the associated error type SHOULD be "modify" or "wait" depending on the policy
// being violated.
type XMPPStanzaErrorPolicyVoilation struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas policy-violation"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.13
//
// The intended recipient is temporarily unavailable, undergoing maintenance, etc.;
// the associated error type SHOULD be "wait".
type XMPPStanzaErrorRecipientUnavailable struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas recipient-unavailable"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.14
//
// The recipient or server is redirecting requests for this information to another
// entity, typically in a temporary fashion (as opposed to the <gone/> error condition,
// which is used for permanent addressing failures); the associated error type SHOULD
// be "modify" and the error stanza SHOULD contain the alternate address in the XML
// character data of the <redirect/> element (which MUST be a URI or IRI with which
// the sender can communicate, typically an XMPP IRI as specified in [XMPP‑URI]).
type XMPPStanzaErrorRedirect struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas redirect"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.15
//
// The requesting entity is not authorized to access the requested service because
// prior registration is necessary (examples of prior registration include members-only
// rooms in XMPP multi-user chat [XEP‑0045] and gateways to non-XMPP instant messaging
// services, which traditionally required registration in order to use the gateway
// [XEP‑0100]); the associated error type SHOULD be "auth".
type XMPPStanzaErrorRegistrationRequired struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas registration-required"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.16
//
// A remote server or service specified as part or all of the JID of the intended
// recipient does not exist or cannot be resolved (e.g., there is no _xmpp-server._tcp
// DNS SRV record, the A or AAAA fallback resolution fails, or A/AAAA lookups succeed
// but there is no response on the IANA-registered port 5269); the associated error
// type SHOULD be "cancel".
type XMPPStanzaErrorRemoteServerNotFound struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas remote-server-not-found"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.17
//
// A remote server or service specified as part or all of the JID of the intended
// recipient (or needed to fulfill a request) was resolved but communications could
// not be established within a reasonable amount of time (e.g., an XML stream cannot
// be established at the resolved IP address and port, or an XML stream can be
// established but stream negotiation fails because of problems with TLS, SASL,
// Server Dialback, etc.); the associated error type SHOULD be "wait" (unless the error
// is of a more permanent nature, e.g., the remote server is found but it cannot be
// authenticated or it violates security policies).
type XMPPStanzaErrorRemoteServerTimeout struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas remote-server-timeout"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.18
//
// The server or recipient is busy or lacks the system resources necessary to service
// the request; the associated error type SHOULD be "wait".
type XMPPStanzaErrorResourceConstraint struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas resource-constraint"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.19
//
// The server or recipient does not currently provide the requested service; the
// associated error type SHOULD be "cancel".
type XMPPStanzaErrorServiceUnavailable struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas service-unavailable"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.20
//
// The requesting entity is not authorized to access the requested service because
// a prior subscription is necessary (examples of prior subscription include
// authorization to receive presence information as defined in [XMPP‑IM] and opt-in
// data feeds for XMPP publish-subscribe as defined in [XEP‑0060]); the associated
// error type SHOULD be "auth".
type XMPPStanzaErrorSubscriptionRequired struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas subscription-required"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
}

// RFC6120 Section 8.3.3.21
//
// The error condition is not one of those defined by the other conditions in this
// list; any error type can be associated with this condition, and it SHOULD NOT
// be used except in conjunction with an application-specific condition.
type XMPPStanzaErrorUndefinedCondition struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas undefined-condition"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
    Error   string                          `xml:",innerxml,omitempty"`
}

// RFC6120 Section 8.3.3.22
//
// The recipient or server understood the request but was not expecting it at this
// time (e.g., the request was out of order); the associated error type SHOULD be
// "wait" or "modify".
type XMPPStanzaErrorUnexpectedRequest struct {
    XMLName xml.Name                        `xml:"urn:ietf:params:xml:ns:xmpp-stanzas unexpected-request"`
    Text    *XMPPStanzaErrorDescriptiveText `xml:",omitempty"`
    Error   string                          `xml:",innerxml,omitempty"`
}

type XMPPCustom struct {
    XMLName xml.Name
    XML     string `xml:",any"`
}

const stream_response_begin_fmt = `<stream:stream from='%s' to='%s' version='%s' xml:lang='%s' id='%s' xmlns='%s' xmlns:stream='%s'>`

func GenXMPPStreamHeader(s *XMPPStream) string {
    return fmt.Sprintf(stream_response_begin_fmt,
        s.From,
        s.To,
        s.Version,
        s.XmlLang,
        s.Id,
        XMLNS_JABBER_CLIENT,
        XMLNS_STREAM)
}

const stream_end_fmt = `</stream:stream>`
