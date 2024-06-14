package rebind

import (
	"net"
	"testing"

	"github.com/coredns/caddy"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
		zones     []string
		firstIP   net.IP
		secondIP  net.IP
		strategy  string
	}{
		// oks
		{`rebind {
			first_ip 1.2.3.4
			second_ip 0.0.0.0
		}`, false, nil, net.ParseIP("1.2.3.4"), net.ParseIP("0.0.0.0"), firstThenSecondStrategy},
		{`rebind {
			first_ip 1.2.3.4
			second_ip 0.0.0.0
			strategy first_then_second
		}`, false, nil, net.ParseIP("1.2.3.4"), net.ParseIP("0.0.0.0"), firstThenSecondStrategy},
		{`rebind {
			first_ip 1.2.3.4
			second_ip 0.0.0.0
			strategy round_robin
		}`, false, nil, net.ParseIP("1.2.3.4"), net.ParseIP("0.0.0.0"), roundRobinStrategy},
		{`rebind {
			first_ip 1.2.3.4
			second_ip 0.0.0.0
			strategy random
		}`, false, nil, net.ParseIP("1.2.3.4"), net.ParseIP("0.0.0.0"), randomStrategy},

		// fails
		{`rebind`, true, nil, nil, nil, ""},
		{`rebind {
			first_ip notanip
			second_ip 0.0.0.0
		}`, true, nil, nil, nil, ""},
		{`rebind {
			first_ip 1.2.3.4
			second_ip 0.0.0.0
			strategy notastrategy
		}`, true, nil, nil, nil, ""},
	}
	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		r, err := parse(c)
		if test.shouldErr && err == nil {
			t.Errorf("Test %v: Expected error but found nil", i)
			continue
		} else if !test.shouldErr && err != nil {
			t.Errorf("Test %v: Expected no error but found error: %v", i, err)
			continue
		}

		if test.shouldErr {
			continue
		}

		if len(r.zones) != len(test.zones) {
			t.Errorf("Test %d expected %v, got %v", i, test.zones, r.zones)
		}
		for j, name := range test.zones {
			if r.zones[j] != name {
				t.Errorf("Test %d expected %v for %d th zone, got %v", i, name, j, r.zones[j])
			}
		}
		if !r.firstIP.Equal(test.firstIP) {
			t.Errorf("Test %d expected %v, got %v", i, test.firstIP, r.firstIP)
		}

		if !r.secondIP.Equal(test.secondIP) {
			t.Errorf("Test %d expected %v, got %v", i, test.secondIP, r.secondIP)
		}

		if r.strategy != test.strategy {
			t.Errorf("Test %d expected %v, got %v", i, test.strategy, r.strategy)
		}

	}
}
