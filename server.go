package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type server struct {
	sync.RWMutex

	ctx    context.Context
	i      int
	events []AnnotationResponse
}

func newServer() *server {
	return &server{ctx: context.Background()}
}

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

// root exists so that jsonds can be successfully added as a Grafana Data Source.
//
// If this exists then Grafana emits this when adding the datasource:
//
//		Success
// 		Data source is working
//
// otherwise it emits "Unknown error"
func (s *server) root(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v: %v", r.URL.Path, r.Method)
	fmt.Fprintf(w, "ok\n")
}

func (s *server) search(w http.ResponseWriter, r *http.Request) {

	var targets []string

	resp, err := json.Marshal(targets)
	if err != nil {

		http.Error(w, "cannot marshal targets response", http.StatusBadRequest)
	}
	w.Write(resp)

}
func (s *server) query(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Sorry not yet", http.StatusNotImplemented)
}

func (s *server) annotations(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v: %v", r.URL.Path, r.Method)
	switch r.Method {
	case http.MethodOptions:
	case http.MethodPost:
		ar := AnnotationsReq{}
		if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
			http.Error(w, fmt.Sprintf("json decode failure: %v", err), http.StatusBadRequest)
			return
		}

		evs := s.filterEvents(ar.Annotation, ar.Range.From, ar.Range.To)
		if err := json.NewEncoder(w).Encode(evs); err != nil {
			log.Printf("json enc: %+v", err)
		}
	default:
		http.Error(w, "bad method; supported OPTIONS, POST", http.StatusBadRequest)
		return
	}
}

func (s *server) filterEvents(a Annotation, from, to time.Time) []AnnotationResponse {
	events := []AnnotationResponse{}
	for _, event := range s.events {
		event.Annotation = a
		event.Annotation.ShowLine = true
		if event.Time > from.Unix()*1000 && event.Time < to.Unix()*1000 {
			events = append(events, event)
		}
	}
	return events
}

// annResp isn't required; it just codifies a standard AnnotationResponse
// between the seed and generate funcs.
func annResp(t time.Time, i int) AnnotationResponse {
	return AnnotationResponse{
		// Grafana expects unix milliseconds:
		// https://github.com/grafana/simple-json-datasource#annotation-api
		Time: t.Unix() * 1000,

		Title: fmt.Sprintf("event %04d", i),
		Text:  fmt.Sprintf("text about the event %04d", i),
		Tags:  "atag btag ctag",
	}
}
