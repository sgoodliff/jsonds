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
	length := int(d.Minutes())

	ir := IntRange{0, 20}
	rand.Seed(time.Now().UTC().UnixNano())

	log.Print("OUTPUT")

	pointsInRange := make([]row, 0, length)

	log.Println("starting minute iteration")

	// 1 minute resolution
	x := from
	//for x := from; x.Minute() <= to.Minute(); x = x.Add(time.Minute * 1) {
	for i := 0; i <= length; i++ {

		c := Count{T: x, N: float64(ir.NextRandom()), U: x.UnixNano() / 1000000}

		pointsInRange = append(pointsInRange, row{c.N, c.U})

		log.Println(c)
		x = x.Add(time.Minute * 1)
	}

	log.Println("ending minute iteration")

	return &pointsInRange
}
