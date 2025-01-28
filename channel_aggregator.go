package main

import (
	"fmt"
	"time"
)

// Channel-based aggregator that reports the global average temperature periodically

// Report the averagage temperature across all `k` weatherstations every `averagePeriod`
// seconds by sending a `WeatherReport` struct to the `out` channel. The aggregator should
// terminate upon receiving a singnal on the `quit` channel.

// Note! To receive credit, channelAggregator must not use mutexes.

func channelAggregator(
	k int,
	averagePeriod float64,
	getWeatherData func(int, int) WeatherReport,
	out chan WeatherReport,
	quit chan struct{},
) {
	duration := time.Duration(averagePeriod * float64(time.Second))
	ticker := time.NewTicker(duration)

	tempData := make(chan WeatherReport, k) // could try setting buffer to 1

	batch := 0

	total := 0.0 // Shorthand declaration with := operator, infers float64 type

	count := 0

	completed := 0

	fmt.Println(averagePeriod)

	for i := 0; i < k; i++ { // i = station ID

		go func(i int, batch int) {
			//fmt.Println(i)
			//fmt.Printf("i: %d, batch: %d\n", i, batch) // Print the values of i and batch
			tempData <- getWeatherData(i, batch) // maybe this needs a -1
		}(i, batch)
	}
	for {
		select {
		case <-ticker.C: //newbatch   //could try making this the defaukt case for the first time it runs.
			average := (total / float64(count))

			out <- WeatherReport{Value: average, Id: -1, Batch: batch}
			batch++

			count = 0
			total = 0.0

			for i := 0; i < k; i++ { // i = station ID

				go func(i int, batch int) {
					completed++
					//fmt.Println(i)
					//fmt.Printf("i: %d, batch: %d\n", i, batch) // Print the values of i and batch
					tempData <- getWeatherData(i, batch) // maybe this needs a -1
				}(i, batch)

			}
		case report := <-tempData:
			if report.Batch == batch {
				total += report.Value
				count++

			}
		case <-quit:
			return
		}

	}

}

// case tempData <- getWeatherData(i, batch)
// 				if: getWeatherData.Batch = batch :
// 					total <- tempData.value
// 					count++

// package main

// import (
// 	"math"
// 	"time"
// )

// // Channel-based aggregator that reports the global average temperature periodically
// //
// // Report the averagage temperature across all `k` weatherstations every `averagePeriod`
// // seconds by sending a `WeatherReport` struct to the `out` channel. The aggregator should
// // terminate upon receiving a singnal on the `quit` channel.
// //
// // Note! To receive credit, channelAggregator must not use mutexes.
// func channelAggregator(
// 	k int,
// 	averagePeriod float64,
// 	getWeatherData func(int, int) WeatherReport,
// 	out chan WeatherReport,
// 	quit chan struct{},
// ) {

// 	duration := time.Duration(averagePeriod * float64(time.Millisecond))

// 	ticker := time.NewTicker(duration)   //Ticker used to start a new batch
// 	tempData := make(chan WeatherReport) //making a buffered channel that transmit and receives values of type WeatherReport with capacity k

// 	for batch := 0; ; batch++ {
// 		//starting batch loop that goes infinitely incremantaly until there is a break inside of loop

// 		count := 0
// 		totalTemp := 0.00

// 		for i := 0; i < k; i++ {
// 			//starting weather station loop that goes till 100 or it is broken
// 			go func(i, batch int) { //ananymous function that take i and batch ID
// 				select {
// 				case tempData <- getWeatherData(i, batch): //putting the weather data into the channel

// 				case <-time.After(time.Duration(averagePeriod) * time.Second): //times out if it takes longer than the average period

// 				}

// 			}(i, batch) //provided parameters for the go routine

// 			timer := time.NewTimer(duration)

// 		innerloop:
// 			for { //repeats until broken
// 				select {
// 				case report := <-tempData: //initializing report and puts tempData(channel) into it
// 					totalTemp += report.Value //adding value of report to totaltemp
// 					count++
// 					for i := 0; i < k; i++ {
// 						//starting weather station loop that goes till 100 or it is broken
// 						go func(i, batch int) { //ananymous function that take i and batch ID
// 							select {
// 							case tempData <- getWeatherData(i, batch): //putting the weather data into the channel

// 							case <-time.After(time.Duration(averagePeriod) * time.Second): //times out if it takes longer than the average period

// 							}

// 						}(i, batch) //provided parameters for the go routine
// 						if count == k { //break if all reports are in
// 							break innerloop
// 						}
// 					}
// 				case <-timer.C:
// 					break innerloop // run k loop // reset sum and num_successful_calls to 0
// 				case <-quit: // return if quitchannel is called
// 					return

// 				}
// 			}

// 			if count > 0 {
// 				out <- WeatherReport{Value: totalTemp / float64(count), Batch: batch}
// 			} else {
// 				out <- WeatherReport{Value: math.NaN(), Batch: batch}
// 			}
// 			<-ticker.C

// 		}
// 	}

// three case:
///signal weather channel(channel) check wehter bacth id = batch _> add temp tp sum add to count
/// case where ticker goes off  which resets all intitialized variable and startthe k loop
/// quit case
