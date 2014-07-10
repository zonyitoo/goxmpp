package xmpp

type xmppClientStreams map[string]Stream

type ServerDispatcher struct {
    clients map[BareJID]xmppClientStreams
}

func NewServerDispatcher() *ServerDispatcher {
    return &ServerDispatcher{
        clients: make(map[BareJID]xmppClientStreams),
    }
}

func (sd *ServerDispatcher) RegisterClient(jid *JID, s Stream) {
    if client, ok := sd.clients[jid.BareJID]; !ok {
        cs := make(xmppClientStreams)
        cs[jid.Resource] = s
        sd.clients[jid.BareJID] = cs
    } else {
        client[jid.Resource] = s
    }
}

func (sd *ServerDispatcher) UnregisterClient(jid *JID) {
    if client, ok := sd.clients[jid.BareJID]; ok {
        delete(client, jid.Resource)

        if len(client) == 0 {
            delete(sd.clients, jid.BareJID)
        }
    }
}
