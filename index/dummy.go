package index

import (
	"rides/reader"
)

type Dummy struct {
	files []reader.File
}

func (m *Dummy) Init(files []reader.File) {
	m.files = files
}

func (m *Dummy) GetAverageDistances(start, end int64) map[int]float64 {
	result := make(map[int]float64)
	count := make(map[int]int)
	for _, f := range m.files {
		for cnt, rides := range f {
			for _, ride := range rides {
				if start <= ride.StartTime && ride.FinishTime <= end {
					result[cnt] += ride.Distance
					count[cnt] += 1
				}
			}
		}
	}
	for cnt, dist := range result {
		result[cnt] = dist / float64(count[cnt])
	}
	return result
}
