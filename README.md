This project simulates a system that aggregates temperature data from multiple weather stations using Go's concurrency model with goroutines and channels. The main purpose of the system is to compute the global average temperature periodically and report it. 

The program uses a channel-based aggregator that collects weather data from `k` weather stations and reports the average temperature every `averagePeriod` seconds. The aggregator terminates upon receiving a signal on the `quit` channel.

## Features
- **Concurrency**: Leverages Go's goroutines and channels to aggregate data concurrently.
- **No Mutexes**: The aggregator is designed to avoid using mutexes, relying on channels to manage synchronization.
- **Periodic Reporting**: The global average temperature is reported at specified intervals.
- **Fault Tolerance**: Handles timeouts and ensures data consistency using channels.

## Code Overview
- **`channelAggregator` function**: This is the main function that aggregates weather data from the stations and periodically reports the global average temperature. 
- **Weather Stations**: Each station runs in its own goroutine and sends weather data to the aggregator.
- **`WeatherReport` struct**: The structure used to represent the weather data, which includes temperature values, batch numbers, and other relevant information.
