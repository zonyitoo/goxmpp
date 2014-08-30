package server

type Server interface {
    Accept() Client
    Serve()
}
