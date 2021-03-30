package server

type Server interface {
	Serve() error
}

type Config struct {
	Protocol string
	Host     string
	Port     string
}
