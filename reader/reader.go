package reader

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const tsLayout = "2006-01-02 15:04:05"

type Ride struct {
	StartTime  int64
	FinishTime int64
	Distance   float64
}

type File = map[int][]Ride

func readCSV(r *csv.Reader) (File, error) {
	f := File{}

	// VendorID,tpep_pickup_datetime,tpep_dropoff_datetime,passenger_count,trip_distance
	// 1,2020-01-01 00:28:15,2020-01-01 00:33:03,1,1.20
	// Skip header
	_, err := r.Read()

	if err != nil {
		return nil, err
	}

	for i := 0; ; i++ {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if record[0] == "" {
			// VendroID is empty, skip
			continue
		}

		startTime, err := time.Parse(tsLayout, record[1])
		if err != nil {
			return nil, err
		}

		finishTime, err := time.Parse(tsLayout, record[2])
		if err != nil {
			return nil, err
		}

		passengerCount, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		tripDistance, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}

		f[passengerCount] = append(f[passengerCount], Ride{
			StartTime:  startTime.Unix(),
			FinishTime: finishTime.Unix(),
			Distance:   tripDistance,
		})
	}
	return f, nil
}

func ReadFiles(path string) ([]File, error) {
	var result []File
	entries, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		f, err := os.Open(filepath.Join(path, entry.Name()))
		if err != nil {
			return nil, err
		}
		log.Printf("reading file %s", entry.Name())
		csvF := csv.NewReader(f)
		parsedFile, err := readCSV(csvF)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedFile)
	}
	return result, err
}
