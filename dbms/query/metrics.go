// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "fmt"

type metrics struct {
	fixcost  Cost
	varcost  Cost
	costself Cost
	frac     float64
	ngets    int32
	nsels    int32
	nlooks   int32
	tget     uint64
	tgetself uint64
}

func (m *metrics) String() string {
	return fmt.Sprintf("metrics{fixcost: %v varcost: %v costself: %v frac: %.2f ngets: %d nsels: %d nlooks: %d tget: %d tgetself: %d}",
		m.fixcost, m.varcost, m.costself, m.frac, m.ngets, m.nsels, m.nlooks, m.tget, m.tgetself)
}

func (m *metrics) setCost(frac float64, fixcost, varcost Cost) {
	m.frac = frac
	m.fixcost = fixcost
	m.varcost = varcost
}

func CalcSelf(q0 Query) { // recursive
	m := q0.Metrics()
	if m.tgetself != 0 {
		return // already calculated
	}
	switch q := q0.(type) {
	case q2i:
		m1 := q.Source().Metrics()
		m2 := q.Source2().Metrics()
		m.tgetself = m.tget - (m1.tget + m2.tget)
		m.costself = (m.fixcost + m.varcost) -
			(m1.fixcost + m1.varcost + m2.fixcost + m2.varcost)
		CalcSelf(q.Source())
		CalcSelf(q.Source2())
	case q1i:
		sm := q.Source().Metrics()
		m.tgetself = m.tget - sm.tget
		m.costself = (m.fixcost + m.varcost) - (sm.fixcost + sm.varcost)
		CalcSelf(q.Source())
	default:
		m.tgetself = q0.Metrics().tget
		m.costself = q0.Metrics().fixcost + q0.Metrics().varcost
	}
}
