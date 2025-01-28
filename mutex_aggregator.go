package main

// Mutex-based aggregator that reports the global average temperature periodically
//
// Report the averagage temperature across all `k` weatherstations every `averagePeriod`
// seconds by sending a `WeatherReport` struct to the `out` channel. The aggregator should
// terminate upon receiving a singnal on the `quit` channel.
//
// Note! To receive credit, mutexAggregator must implement a mutex based solution.

import (
	"math"
	"sync"
	"time"
)

func mutexAggregator(
	k int,
	averagePeriod float64,
	getWeatherData func(int, int) WeatherReport,
	out chan WeatherReport,
	quit chan struct{},
) {

	duration := time.Duration(averagePeriod * float64(time.Second))
	ticker := time.NewTicker(duration)

	var mutex sync.Mutex
	globalBatch := 0
	total := 0.0
	count := 0

	for i := 0; i < k; i++ {
		go func(j int, batch int) {
			report := getWeatherData(j, batch) // t+1 second
			mutex.Lock()
			if report.Batch == globalBatch {
				total += report.Value
				count++
			}
			mutex.Unlock()
		}(i, globalBatch)
	}

	for {
		select {
		case <-ticker.C:

			mutex.Lock()
			out <- WeatherReport{Value: (total / float64(count)), Id: -1, Batch: globalBatch}

			globalBatch++
			total = 0.0
			count = 0
			mutex.Unlock()

			for i := 0; i < k; i++ {
				go func(j int, batch int) {

					report := getWeatherData(j, batch)
					mutex.Lock()
					if report.Batch == globalBatch {
						if !math.IsNaN(report.Value) {
							total += report.Value
							count++
						}
					}
					mutex.Unlock()

				}(i, globalBatch)
			}

		case <-quit:
			return

		}
	}
}
