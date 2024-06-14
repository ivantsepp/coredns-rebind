package rebind

import (
	"context"
	"math/rand"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

const (
	firstThenSecondStrategy = "first_then_second"
	randomStrategy          = "random"
	roundRobinStrategy      = "round_robin"
)

// Rebind is a plugin that rebinds a query to a different IP address based on the strategy.
type Rebind struct {
	Next     plugin.Handler
	firstIP  net.IP
	secondIP net.IP
	zones    []string
	visited  map[string]int
	strategy string
}

func (a Rebind) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	state := request.Request{W: w, Req: r}
	qname := state.Name()
	answers := []dns.RR{}

	// If the zone does not match one of ours, just pass it on.
	if plugin.Zones(a.zones).Matches(qname) == "" {
		return plugin.NextOrFailure(a.Name(), a.Next, ctx, w, r)
	}

	ipAnswer := a.getIPAnswer(qname)

	answers = append(answers, &dns.A{
		Hdr: dns.RR_Header{
			Name:   qname,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		A: ipAnswer,
	})

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

func (a Rebind) Name() string { return "rebind" }

func (a Rebind) getIPAnswer(qname string) net.IP {
	switch a.strategy {
	case randomStrategy:
		return a.getIPAnswerRandom()
	case roundRobinStrategy:
		return a.getIPAnswerRoundRobin(qname)
	case firstThenSecondStrategy:
		return a.getIPAnswerFirstThenSecond(qname)
	}

	return a.firstIP
}

func (a Rebind) getIPAnswerRandom() net.IP {
	possibleAnswers := []net.IP{a.firstIP, a.secondIP}
	return possibleAnswers[rand.Intn(len(possibleAnswers))]
}

func (a Rebind) getIPAnswerRoundRobin(qname string) net.IP {
	_, ok := a.visited[qname]
	if !ok {
		a.visited[qname] = 0
	}

	switch a.visited[qname] {
	case 0:
		a.visited[qname] = 1
		return a.firstIP
	case 1:
		a.visited[qname] = 2
		return a.secondIP
	case 2:
		a.visited[qname] = 1
		return a.firstIP
	}

	return a.firstIP
}

func (a Rebind) getIPAnswerFirstThenSecond(qname string) net.IP {
	_, ok := a.visited[qname]
	if !ok {
		a.visited[qname] = 1
		return a.firstIP
	} else {
		return a.secondIP
	}
}
