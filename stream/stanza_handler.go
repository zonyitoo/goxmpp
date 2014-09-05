package stream

import "github.com/zonyitoo/goxmpp/protocol"

type StanzaHandler interface {
    HandleIQ(*protocol.XMPPStanzaIQ, Streamer) error
    HandleMessage(*protocol.XMPPStanzaMessage, Streamer) error
    HandlePresence(*protocol.XMPPStanzaPresence, Streamer) error
}
