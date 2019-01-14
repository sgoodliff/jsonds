package main

import (
	"time"
)

func (s *server) seed(max int) {
	s.RLock()
	defer s.RUnlock()

	expansion := 20 * time.Minute
	n := time.Now().Add(-(expansion * time.Duration(max)))
	for i := 0; i < max; i++ {
		t := n.Add(time.Duration(i+1) * expansion)
		s.events = append(s.events, annResp(t, i))
		s.i++
	}
}

func (s *server) generate(period time.Duration) {
	t := time.NewTicker(period)
	for {
		select {
		case <-t.C:
			n := time.Now()
			s.RLock()
			s.events = append(s.events, annResp(n, s.i))
			s.i++
			s.RUnlock()
		case <-s.ctx.Done():
			return
		}
	}
}
