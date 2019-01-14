package main

import (
	"log"
	"math/rand"
	"time"
)

// Count is a single time series data tuple, consisting of
// a floating-point value N and a timestamp T.
type Count struct {
	N float64
	T time.Time
	U int64
}

func fetchDatapoints(from, to time.Time, maxDataPoints int) *[]row {

	log.Println("INPUTS")
	log.Println(to)
	log.Println(from)
	log.Println(maxDataPoints)

	log.Println("TIMINGS")
	d := to.Sub(from)
	d.Round(time.Minute)
	log.Println(d.Minutes())
	points := int(d.Minutes())

	ir := IntRange{0, 20}
	rand.Seed(time.Now().UTC().UnixNano())

	log.Print("OUTPUT")
	// need something like  [622,1450754160000],  // Metric value as a float , unixtimestamp in milliseconds
	for i := 0; i < points; i++ {

		log.Println(ir.NextRandom())

	}

	rows := make([]row, maxDataPoints)
	log.Println("starting minute iteration")

	for x := from; x.Minute() <= to.Minute(); x = x.Add(time.Minute * 1) {
		c := Count{T: x, N: float64(ir.NextRandom()), U: x.Unix()}

		//log.Println(ir.NextRandom())
		//log.Println(int64(x.UnixNano()))
		log.Println(c)

	}

	log.Println("ending minute iteration")

	return &rows
}
