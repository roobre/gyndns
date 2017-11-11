package gyndns

import (
	"net"
	"sync"
	"encoding/json"
	"log"
)

type Config struct {
	HTTPAddress string
	HTTPPort    uint16

	DNSAddress string
	DNSPort    uint16
}

type GynDNS struct {
	Config

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
	HTTPAddress: "0.0.0.0",
	HTTPPort:    8000,
	DNSAddress:  "0.0.0.0",
	DNSPort:     5533,
}

func New(config *Config, usersFile []byte) *GynDNS {
	if config == nil {
		config = &defaultConfig
	}

	g := &GynDNS{
		Config:  *config,
		errChan: make(chan error),
		users:   make(map[Username]User),
		leases:  make(map[string]net.IP),
	}

	var users []User
	err := json.Unmarshal(usersFile, &users)
	if err != nil {
		panic(err)
	}

	for _, u := range users {
		g.users[u.Username] = u
	}

	return g
}

func (g *GynDNS) Run() {
	go g.runHTTP(g.errChan)
	go g.runDNS(g.errChan)

	log.Fatal(<-g.errChan)
}
