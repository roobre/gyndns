package gyndns

import (
	"github.com/miekg/dns"
	"fmt"
	"log"
)

const TTL = 16

func (g *GynDNS) runDNS(errChan chan error) {
	log.Println("Starting DNS server...")
	errChan <- dns.ListenAndServe(fmt.Sprintf("%s:%d", g.DNSAddress, g.DNSPort), "udp", g)
}

func (g *GynDNS) ServeDNS(rw dns.ResponseWriter, r *dns.Msg) {
	for _, q := range r.Question {
		if q.Qtype != dns.TypeA {
			log.Printf("Unsupported question type %d", q.Qtype)
		} else {
			log.Printf("Searching for hostname '%s'", q.Name)
			response := &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Id:            r.Id,
					Response:      true,
					Authoritative: true,
				},
			}

			response.Question = append(response.Question, q)

			g.lMutex.RLock()
			ip, found := g.leases[q.Name]
			g.lMutex.RUnlock()

			if found {
				response.Answer = append(response.Answer, &dns.A{
					Hdr: dns.RR_Header{
						Name:   q.Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    TTL,
					},
					A: ip,
				})
			} else {
				response.Rcode = dns.RcodeNameError
			}

			rw.WriteMsg(response)

			break
		}
	}
}
