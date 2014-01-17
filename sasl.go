package xmpp

type XMPPSASLServerMechanismHandler interface {
	Auth(data string) (string, interface{}, error)
	Response(data string) (interface{}, error)
	Abort() (interface{}, error)
}

type XMPPSASLClientMechanismHandler interface {
	Begin() (interface{}, error)
	Challenge(data string) (interface{}, error)
	Success(data string) (interface{}, error)
	Failure(data string) (interface{}, error)
	Abort() (interface{}, error)
}
