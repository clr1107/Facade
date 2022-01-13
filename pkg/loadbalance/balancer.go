package loadbalance

import (
	"sync/atomic"
)

// BalancerMedium is a medium for balancing, e.g. a http proxy, or another
// network interface to use.
type BalancerMedium interface {

}

// Balancer is an interface for algorithms that can balance loads around
// various media (`BalancerMedium`)
// Must be thread safe.
type Balancer interface {
	Get() BalancerMedium
}

// RoundRobin implements the Round Robin algorithm to return a new BalancerMedium.
type RoundRobin struct {
	balancers []BalancerMedium
	index     *uint32
}

func NewRoundRobinPool(balancers []BalancerMedium) RoundRobin {
	return RoundRobin{
		balancers: balancers,
		index:     new(uint32),
	}
}

// Get retrieves, using the balancing algorithm, the next BalancerMedium.
func (p *RoundRobin) Get() BalancerMedium {
	return p.balancers[(int(atomic.AddUint32(p.index, 1)) - 1) % len(p.balancers)]
}
