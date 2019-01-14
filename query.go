package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Anything to sort out /query can go here

// query is a `/query` request from Grafana.
//
// All JSON-related structs were generated from the JSON examples
// of the "SimpleJson" data source documentation
// using [JSON-to-Go](https://mholt.github.io/json-to-go/),
// with a little tweaking afterwards.
type query struct {
	PanelID int `json:"panelId"`
	Range   struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
		Raw  struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"raw"`
	} `json:"range"`
	RangeRaw struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"rangeRaw"`
	Interval   string `json:"interval"`
	IntervalMs int    `json:"intervalMs"`
	Targets    []struct {
		Target string `json:"target"`
		RefID  string `json:"refId"`
		Type   string `json:"type"`
	} `json:"targets"`
	Format        string `json:"format"`
	MaxDataPoints int    `json:"maxDataPoints"`
}

// row is used in timeseriesResponse and tableResponse.
// Grafana's JSON contains weird arrays with mixed types!
type row []interface{}

// column is used in tableResponse.
type column struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// timeseriesResponse is the response to a `/query` request
// if "Type" is set to "timeserie".
// It sends time series data back to Grafana.
type timeseriesResponse struct {
	Target     string `json:"target"`
	Datapoints []row  `json:"datapoints"`
}

// tableResponse is the response to send when "Type" is "table".
type tableResponse struct {
	Columns []column `json:"columns"`
	Rows    []row    `json:"rows"`
	Type    string   `json:"type"`
}

func (s *server) query(w http.ResponseWriter, r *http.Request) {

	if (*r).Method == "OPTIONS" {
		return
	}

	//http.Error(w, "Sorry not yet", http.StatusNotImplemented)
	log.Println("Doing a query")
	var q bytes.Buffer

	_, err := q.ReadFrom(r.Body)
	if err != nil {
		writeError(w, err, "Cannot read request body")
		return
	}

	query := &query{}
	err = json.Unmarshal(q.Bytes(), query)
	if err != nil {
		writeError(w, err, "cannot unmarshal request body")
		return
	}

	log.Println(query.Targets[0].Type)

	switch query.Targets[0].Type {
	case "timeserie":
		log.Println("TimeSeries Query")
		s.sendTimeseries(w, query)
	case "table":
		http.Error(w, "Sorry not yet", http.StatusNotImplemented)

	default:
		http.Error(w, "Fall Through", http.StatusNotImplemented)
	}
	log.Println("Leaving Query")
}

// sendTimeseries creates and writes a JSON response to a request for time series data.
func (s *server) sendTimeseries(w http.ResponseWriter, q *query) {

	response := []timeseriesResponse{}
	// SJG GO FETCH STUFF

	for _, t := range q.Targets {
		target := t.Target
		/*
			metric, err := s.metrics.Get(target)
			if err != nil {
				writeError(w, err, "Cannot get metric for target "+target)
				return
			}
		*/
		response = append(response, timeseriesResponse{
			Target:     target,
			Datapoints: *(fetchDatapoints(q.Range.From, q.Range.To, q.MaxDataPoints)),
		})
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		writeError(w, err, "cannot marshal timeseries response")
	}

	w.Write(jsonResp)

}
