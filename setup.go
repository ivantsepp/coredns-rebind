package rebind

import (
	"net"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("rebind", setup) }

func setup(c *caddy.Controller) error {
	r, err := parse(c)

	if err != nil {
		return plugin.Error("rebind", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		r.Next = next
		return r
	})

	return nil
}

func parse(c *caddy.Controller) (Rebind, error) {
	r := Rebind{
		strategy: firstThenSecondStrategy,
	}
	found := false

	for c.Next() {

		if found {
			return r, plugin.ErrOnce
		}
		found = true

		args := c.RemainingArgs()
		r.zones = plugin.OriginsFromArgsOrServerBlock(args, c.ServerBlockKeys)

		for c.NextBlock() {
			switch c.Val() {
			case "first_ip":
				remaining := c.RemainingArgs()
				if len(remaining) < 1 {
					return r, c.Errf("first_ip needs an IP address")
				}
				r.firstIP = net.ParseIP(remaining[0])
			case "second_ip":
				remaining := c.RemainingArgs()
				if len(remaining) < 1 {
					return r, c.Errf("second_ip needs an IP address")
				}
				r.secondIP = net.ParseIP(remaining[0])
			case "strategy":
				remaining := c.RemainingArgs()
				if len(remaining) < 1 {
					return r, c.Errf("strategy needs to be one of [first_then_second, random, round_robin]")
				}
				switch remaining[0] {
				case firstThenSecondStrategy:
					r.strategy = firstThenSecondStrategy
				case randomStrategy:
					r.strategy = randomStrategy
				case roundRobinStrategy:
					r.strategy = roundRobinStrategy
				default:
					return r, c.Errf("unknown strategy: %s", remaining[0])
				}
			default:
				return r, c.Errf("invalid option %q", c.Val())
			}
		}

		if r.firstIP == nil || r.secondIP == nil {
			return r, c.Err("rebind needs firstIP and secondIP addresses")
		}

		r.visited = make(map[string]int)
	}

	return r, nil
}
