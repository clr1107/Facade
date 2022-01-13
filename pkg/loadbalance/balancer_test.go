package loadbalance

import (
	"testing"
)

type NamedBalancerMedium struct {
	name int
}

func TestRoundRobin_Get(t *testing.T) {
	m := make([]BalancerMedium, 0)
	for i := 0; i < 3; i++ {
		m = append(m, NamedBalancerMedium{name: i})
	}

	r := NewRoundRobinPool(m)
	for i := 0; i < 4; i++ {
		cast, ok := r.Get().(NamedBalancerMedium)
		if !ok {
			t.Errorf("Round robin didn't return a NamedBalancerMedium")
		}

		if cast.name % 3 != i % 3 {
			t.Errorf("Wrong balancer returned from round robin, wanted %v got %v", i % 3, cast.name % 3)
		}
	}
}
