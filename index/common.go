package index

import (
	"rides/reader"
	"time"
)

// Index is an alternative for this java interface:
// interface AverageDistances extends Closeable {
//    /**
//     * Initializes the instance.
//     * @param dataDir Path to the data directory that contains the files
//     *                with CSV data files.
//     */
//    void init(Path dataDir);
//
//    /**
//     * Calculates an average of all trip_distance fields for each distinct
//     * passenger_count value for trips with tpep_pickup_datetime >= start
//     * and tpep_dropoff_datetime <= end.
//     * @return A map where key is passenger count and value is the average
//     *         trip distance for this passenger count.
//     */
//    Map<Integer, Double> getAverageDistances(LocalDateTime start, LocalDateTime end);
//}
type AverageDistances interface {
	Init(path string) error
	GetAverageDistances(start, end time.Time) map[int]float64
}

type Index interface {
	Init(files []reader.File)
	GetAverageDistances(start, end int64) map[int]float64
}

type AverageDistancesImpl struct {
	Index Index
}

func (m *AverageDistancesImpl) Init(path string) error {
	files, err := reader.ReadFiles(path)
	if err != nil {
		return err
	}
	m.Index.Init(files)
	return nil
}
func (m *AverageDistancesImpl) GetAverageDistances(start, end time.Time) map[int]float64 {
	return m.Index.GetAverageDistances(start.Unix(), end.Unix())
}
