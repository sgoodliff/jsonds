package main

import (
	"log"
	"time"
)

func fetchDatapoints(from, to time.Time, maxDataPoints int) *[]row {

	log.Println(to)
	log.Println(from)
	log.Println(maxDataPoints)

	d := to.Sub(from)
	d.Round(time.Minute)
	log.Println(d.Minutes())

	/*
		g.m.Lock()
		defer g.m.Unlock()
		length := len(g.list)

		g.sort()


		// Stage 1: extract all data points within the given time range.
		pointsInRange := make([]row, 0, length)
		for i := 0; i < length; i++ {
			count := g.list[(i+g.head)%length] // wrap around
			if count.T.After(from) && count.T.Before(to) {
				pointsInRange = append(pointsInRange, row{count.N, count.T.UnixNano() / 1000000}) // need ms
			}
		}

		points := len(pointsInRange)

		if points <= maxDataPoints {
			return &pointsInRange
		}

		// Stage 2: if more data points than requested exist in the time range,
		// thin out the slice evenly
		rows := make([]row, maxDataPoints)
		ratio := float64(len(pointsInRange)) / float64(len(rows))
		for i := range rows {
			rows[i] = pointsInRange[int(float64(i)*ratio)]
		}
	*/
	rows := make([]row, maxDataPoints)
	return &rows
}
