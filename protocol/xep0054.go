package protocol

import (
    "encoding/xml"
)

const (
    XMLNS_XMPP_STANZA_IQ_VCARD_TEMP = "vcard-temp"
)

type XMPPStanzaIQVCard struct {
    XMLName       xml.Name                        `xml:"vcard-temp vCard"`
    Version       string                          `xml:"VERSION,omitempty"` // MUST be 2.0 if the document conforms to RFC2426
    FormattedName string                          `xml:"FN,omitempty"`      // Formatted or display name property.
    Name          *XMPPStanzaIQVCardName          `xml:",omitempty"`
    Nicknames     []string                        `xml:"NICKNAME,omitempty"`
    Photos        [][]byte                        `xml:"PHOTO,omitempty"` // Value is either a BASE64 encoded binary value or a URI to the external content.
    Birthdays     []string                        `xml:"BDAY,omitempty"`  // Value must be an ISO 8601 formatted date or date/time value.
    Addresses     []XMPPStanzaIQVCardAddress      `xml",omitempty"`
    Labels        []XMPPStanzaIQVCardAddressLabel `xml:",omitempty"`
    Tels          []XMPPStanzaIQVCardTel          `xml:",omitempty"`
    EMails        []XMPPStanzaIQVCardEMail        `xml:",omitempty"`
    JabberIDs     []string                        `xml:"JABBERID,omitempty"`
    Mailers       []string                        `xml:"MAILER,omitempty"`
    TimeZones     []string                        `xml:"TZ,omitempty"`
    Geos          []XMPPStanzaIQVCardGeo          `xml:",omitempty"`
    Titles        []string                        `xml:"TITLE,omitempty"`
    Roles         []string                        `xml:"ROLE,omitempty"`
    Logos         [][]byte                        `xml:"LOGO,omitempty"`
    AgentVCards   []XMPPStanzaIQVCardAgent        `xml:",omitempty"`
    Organizations []XMPPStanzaIQVCardOrganization `xml:",omitempty"`
    Categoriess   []XMPPStanzaIQVCardCategories   `xml:",omitempty"`
    Notes         []string                        `xml:"NOTE,omitempty"`
    ProdIDs       []string                        `xml:"PRODID,omitempty"`
    Revs          []string                        `xml:"REV,omitempty"`
    SortStrings   []string                        `xml:"SORT-STRING,omitempty"`
    Sounds        []string                        `xml:"SOUND,omitempty"`
    UIDs          []string                        `xml:"UID,omitempty"`
    URLs          []string                        `xml:"URL,omitempty"`
    Classes       []XMPPStanzaIQVCardClass        `xml:",omitempty"`
    Keys          []XMPPStanzaIQVCardKey          `xml:",omitempty"`
    Desc          string                          `xml:"DESC,omitempty"`
}

type XMPPStanzaIQVCardName struct {
    XMLName xml.Name `xml:"N"`
    Family  string   `xml:"FAMILY,omitempty"`
    Given   string   `xml:"GIVEN,omitempty"`
    Middle  string   `xml:"MIDDLE,omitempty"`
    Prefix  string   `xml:"PREFIX,omitempty"`
    Suffix  string   `xml:"SUFFIX,omitempty"`
}

type XMPPStanzaIQVCardAddress struct {
    XMLName    xml.Name                        `xml:"ADR"`
    Home       *XMPPStanzaIQVCardAddressHome   `xml:",omitempty"`
    Work       *XMPPStanzaIQVCardAddressWork   `xml:",omitempty"`
    Postal     *XMPPStanzaIQVCardAddressPostal `xml:",omitempty"`
    Parcel     *XMPPStanzaIQVCardAddressParcel `xml:",omitempty"`
    Dom        *XMPPStanzaIQVCardAddressDom    `xml:",omitempty"`
    Intl       *XMPPStanzaIQVCardAddressIntl   `xml:",omitempty"`
    Perf       *XMPPStanzaIQVCardAddressPref   `xml:",omitempty"`
    Pobox      string                          `xml:"POBOX,omitempty"`
    ExtAddress string                          `xml:"EXTADD,omitempty"`
    Street     string                          `xml:"STREET,omitempty"`
    Locality   string                          `xml:"LOCALITY,omitempty"`
    Region     string                          `xml:"REGION,omitempty"`
    PostalCode string                          `xml:"PCODE,omitempty"`
    Country    string                          `xml:"CTRY,omitempty"`
}

type XMPPStanzaIQVCardAddressHome struct {
    XMLName xml.Name `xml:"HOME"`
}

type XMPPStanzaIQVCardAddressWork struct {
    XMLName xml.Name `xml:"WORK"`
}

type XMPPStanzaIQVCardAddressPostal struct {
    XMLName xml.Name `xml:"POSTAL"`
}

type XMPPStanzaIQVCardAddressParcel struct {
    XMLName xml.Name `xml:"PARCEL"`
}

type XMPPStanzaIQVCardAddressDom struct {
    XMLName xml.Name `xml:"DOM"`
}

type XMPPStanzaIQVCardAddressIntl struct {
    XMLName xml.Name `xml:"INTL"`
}

type XMPPStanzaIQVCardAddressPref struct {
    XMLName xml.Name `xml:"PREF"`
}

type XMPPStanzaIQVCardAddressVoice struct {
    XMLName xml.Name `xml:"VOICE"`
}

type XMPPStanzaIQVCardAddressFax struct {
    XMLName xml.Name `xml:"FAX"`
}

type XMPPStanzaIQVCardAddressPager struct {
    XMLName xml.Name `xml:"PAGER"`
}

type XMPPStanzaIQVCardAddressMsg struct {
    XMLName xml.Name `xml:"MSG"`
}

type XMPPStanzaIQVCardAddressCell struct {
    XMLName xml.Name `xml:"CELL"`
}

type XMPPStanzaIQVCardAddressVideo struct {
    XMLName xml.Name `xml:"VIDEO"`
}

type XMPPStanzaIQVCardAddressBBS struct {
    XMLName xml.Name `xml:"BBS"`
}

type XMPPStanzaIQVCardAddressModem struct {
    XMLName xml.Name `xml:"MODEM"`
}

type XMPPStanzaIQVCardAddressISDN struct {
    XMLName xml.Name `xml:"ISDN"`
}

type XMPPStanzaIQVCardAddressPCS struct {
    XMLName xml.Name `xml:"PCS"`
}

type XMPPStanzaIQVCardAddressInternet struct {
    XMLName xml.Name `xml:"INTERNET"`
}

type XMPPStanzaIQVCardAddressX400 struct {
    XMLName xml.Name `xml:"X400"`
}

type XMPPStanzaIQVCardAddressLabel struct {
    XMLName xml.Name                        `xml:"LABEL"`
    Home    *XMPPStanzaIQVCardAddressHome   `xml:",omitempty"`
    Work    *XMPPStanzaIQVCardAddressWork   `xml:",omitempty"`
    Postal  *XMPPStanzaIQVCardAddressPostal `xml:",omitempty"`
    Parcel  *XMPPStanzaIQVCardAddressParcel `xml:",omitempty"`
    Dom     *XMPPStanzaIQVCardAddressDom    `xml:",omitempty"`
    Intl    *XMPPStanzaIQVCardAddressIntl   `xml:",omitempty"`
    Pref    *XMPPStanzaIQVCardAddressPref   `xml:",omitempty"`
    Line    []string                        `xml:"LINE"`
}

type XMPPStanzaIQVCardTel struct {
    XMLName xml.Name                       `xml:"TEL"`
    Home    *XMPPStanzaIQVCardAddressHome  `xml:",omitempty"`
    Work    *XMPPStanzaIQVCardAddressWork  `xml:",omitempty"`
    Voice   *XMPPStanzaIQVCardAddressVoice `xml:",omitempty"`
    Fax     *XMPPStanzaIQVCardAddressFax   `xml:",omitempty"`
    Pager   *XMPPStanzaIQVCardAddressPager `xml:",omitempty"`
    Msg     *XMPPStanzaIQVCardAddressMsg   `xml:",omitempty"`
    Cell    *XMPPStanzaIQVCardAddressCell  `xml:",omitempty"`
    Video   *XMPPStanzaIQVCardAddressVideo `xml:",omitempty"`
    BBS     *XMPPStanzaIQVCardAddressBBS   `xml:",omitempty"`
    Modem   *XMPPStanzaIQVCardAddressModem `xml:",omitempty"`
    ISDN    *XMPPStanzaIQVCardAddressISDN  `xml:",omitempty"`
    PCS     *XMPPStanzaIQVCardAddressPCS   `xml:",omitempty"`
    Pref    *XMPPStanzaIQVCardAddressPref  `xml:",omitempty"`
    Number  string                         `xml:"NUMBER"`
}

type XMPPStanzaIQVCardEMail struct {
    XMLName  xml.Name                          `xml:"EMAIL"`
    Home     *XMPPStanzaIQVCardAddressHome     `xml:",omitempty"`
    Work     *XMPPStanzaIQVCardAddressWork     `xml:",omitempty"`
    Internet *XMPPStanzaIQVCardAddressInternet `xml:",omitempty"`
    Pref     *XMPPStanzaIQVCardAddressPref     `xml:",omitempty"`
    X400     *XMPPStanzaIQVCardAddressX400     `xml:",omitempty"`
    UserID   string                            `xml:"USERID"`
}

type XMPPStanzaIQVCardGeo struct {
    XMLName   xml.Name `xml:"GEO"`
    Latitude  string   `xml:"LAT"`
    Longitude string   `xml:"LON"`
}

type XMPPStanzaIQVCardOrganization struct {
    XMLName xml.Name `xml:"ORG"`
    Name    string   `xml:"ORGNAME"`
    Unit    []string `xml:"ORGUNIT,omitempty"`
}

type XMPPStanzaIQVCardCategories struct {
    XMLName xml.Name `xml:"CATEGORIES"`
    Keyword []string `xml:"KEYWORD"`
}

type XMPPStanzaIQVCardClass struct {
    XMLName      xml.Name                            `xml:"CLASS"`
    Public       *XMPPStanzaIQVCardClassPublic       `xml:",omitempty"`
    Private      *XMPPStanzaIQVCardClassPrivate      `xml:",omitempty"`
    Confidential *XMPPStanzaIQVCardClassConfidential `xml:",omitempty"`
}

type XMPPStanzaIQVCardClassPublic struct {
    XMLName xml.Name `xml:"PUBLIC"`
}

type XMPPStanzaIQVCardClassPrivate struct {
    XMLName xml.Name `xml:"PRIVATE"`
}

type XMPPStanzaIQVCardClassConfidential struct {
    XMLName xml.Name `xml:"CONFIDENTIAL"`
}

type XMPPStanzaIQVCardKey struct {
    XMLName    xml.Name `xml:"KEY"`
    Type       string   `xml:"TYPE,omitempty"`
    Credential string   `xml:"CRED,omitempty"`
}

type XMPPStanzaIQVCardAgent struct {
    XMLName    xml.Name `xml:"AGENT"`
    ExtAddress string   `xml:",chardata,omitempty"`
    XMPPStanzaIQVCard
}
