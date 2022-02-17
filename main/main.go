package main

import (
	"log"
	"rides/index"
	"time"
)

func main() {
	avgDists := index.AverageDistancesImpl{Index: index.NewBinsearch(8)}
	err := avgDists.Init("data")
	if err != nil {
		log.Fatal(err)
	}
	result := avgDists.GetAverageDistances(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC))
	log.Println("result is", result)

}
