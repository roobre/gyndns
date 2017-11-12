package gyndns

import (
	"net"
	"sync"
	"log"
)

type Config struct {
	HTTPAddress string
	HTTPPort    uint16

	DNSAddress string
	DNSPort    uint16
}

type Params struct {
	Config *Config
	Users  []User
}

type GynDNS struct {
	*Config

	users map[Username]User

	leases map[string]net.IP
	lMutex sync.RWMutex

	errChan chan error
}

type Username string

type User struct {
	Username Username
	Password string
	Names    []string
}

var defaultConfig = Config{
	HTTPAddress: "127.0.0.1",
	HTTPPort:    8000,
	DNSAddress:  "127.0.0.1",
	DNSPort:     5533,
}

func New(params *Params) *GynDNS {
	if params == nil {
		log.Fatal("Nil parametes supplied")
	}

	if params.Config == nil {
		params.Config = &defaultConfig
	}

	g := &GynDNS{
		Config:  params.Config,
		errChan: make(chan error),
		users:   make(map[Username]User),
		leases:  make(map[string]net.IP),
	}

	if len(params.Users) == 0 {
		log.Fatal("No users found in parameters file")
	}

	for _, u := range params.Users {
		g.users[u.Username] = u
	}

	return g
}

func (g *GynDNS) Run() {
	go g.runHTTP(g.errChan)
	go g.runDNS(g.errChan)

	log.Fatal(<-g.errChan)
}
