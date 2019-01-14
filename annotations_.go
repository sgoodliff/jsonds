package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// AnnotationsReq encodes the information provided by Grafana in its requests.
type AnnotationsReq struct {
	Range      Range      `json:"range"`
	Annotation Annotation `json:"annotation"`
}

// Range specifies the time range the request is valid for.
type Range struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// Annotation is the object passed by Grafana when it fetches annotations.
//
// http://docs.grafana.org/plugins/developing/datasources/#annotation-query
type Annotation struct {
	// Name must match in the request and response
	Name string `json:"name"`

	Datasource string `json:"datasource"`
	IconColor  string `json:"iconColor"`
	Enable     bool   `json:"enable"`
	ShowLine   bool   `json:"showLine"`
	Query      string `json:"query"`
}

// AnnotationResponse contains all the information needed to render an
// annotation event.
//
// https://github.com/grafana/simple-json-datasource#annotation-api
type AnnotationResponse struct {
	// The original annotation sent from Grafana.
	Annotation Annotation `json:"annotation"`
	// Time since UNIX Epoch in milliseconds. (required)
	Time int64 `json:"time"`
	// The title for the annotation tooltip. (required)
	Title string `json:"title"`
	// Tags for the annotation. (optional)
	Tags string `json:"tags"`
	// Text for the annotation. (optional)
	Text string `json:"text"`
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
