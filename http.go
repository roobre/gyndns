package gyndns

import (
	"net/http"
	"log"
	"net"
	"fmt"
	"strings"
)

const HOSTNAME_KEY = "hostname"
const IP_KEY = "myip"
const IP_HEADER = "X-Real-IP"

func (g *GynDNS) runHTTP(errChan chan error) {
	addr := fmt.Sprintf("%s:%d", g.HTTPAddress, g.HTTPPort)
	log.Printf("Starting HTTP server at %s...", addr)

	errChan <- http.ListenAndServe(addr, g)
}

func (g *GynDNS) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		http.Error(rw, "No credentials were found", 401)
		return
	}

	user, ok := g.users[Username(username)]
	if !ok {
		http.Error(rw, "Username "+username+" is not registered", 403)
		return
	}
	if user.Password != password {
		http.Error(rw, "Mismatching password for user "+username, 403)
		return
	}

	r.ParseForm()

	hostname := r.Form.Get(HOSTNAME_KEY)
	if hostname == "" {
		http.Error(rw, "Missing '"+HOSTNAME_KEY+"' parameter", 400)
		return
	}

	ip := net.ParseIP(r.Form.Get(IP_KEY))

	if ip == nil {
		ip = net.ParseIP(r.Header.Get(IP_HEADER))
	}

	if ip == nil {
		ip = net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
	}

	if ip == nil {
		http.Error(rw, "Missing '"+IP_KEY+"' parameter ("+r.Form.Get(IP_KEY)+")", 400)
		return
	}

	var found bool
	for _, name := range user.Names {
		if name == hostname {
			found = true
			break
		}
	}

	if !found {
		http.Error(rw, "User "+username+" is not allowed to update "+hostname, 403)
		return
	}

	if !strings.HasSuffix(hostname, ".") {
		hostname = hostname + "."
	}

	log.Printf("Updating %s to IP %s via request from %s", hostname, ip.String(), username)

	g.lMutex.Lock()
	g.leases[hostname] = ip
	g.lMutex.Unlock()

	log.Print(g.leases)
}
