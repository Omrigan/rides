package index

import (
	"log"
	"rides/reader"
	"sort"
	"sync"
)

func iterate(rides []reader.Ride, tsEnd int64, start, end int) (float64, float64) {
	var sum float64
	var cnt float64
	for i := start; i < end; i++ {
		if rides[i].FinishTime <= tsEnd {
			sum += rides[i].Distance
			cnt += 1
		}
	}
	return sum, cnt
}

type Binsearch struct {
	metaFile    reader.File
	concurrency int
}

func NewBinsearch(concurrency int) *Binsearch {
	return &Binsearch{
		concurrency: concurrency,
	}
}

func (m *Binsearch) iterateConc(rides []reader.Ride, tsEnd int64, start, end int) (float64, float64) {
	batchSize := (end - start) / (m.concurrency - 1)
	var wg sync.WaitGroup
	wg.Add(m.concurrency)
	sums := make([]float64, m.concurrency)
	cnts := make([]float64, m.concurrency)
	for i := 0; i < m.concurrency; i++ {
		i := i
		go func() {
			left := start + i*batchSize
			right := left + batchSize
			if right > end {
				right = end
			}
			sum, cnt := iterate(rides, tsEnd, left, right)
			sums[i] = sum
			cnts[i] = cnt
			wg.Done()
		}()
	}
	wg.Wait()
	var sum, cnt float64
	for i := 0; i < m.concurrency; i++ {
		sum += sums[i]
		cnt += cnts[i]
	}
	return sum, cnt
}

func (m *Binsearch) avgDist(rides []reader.Ride, start, end int64) (float64, bool) {
	leftIdx := sort.Search(len(rides), func(i int) bool { return start <= rides[i].StartTime })
	rightIdx := sort.Search(len(rides), func(i int) bool { return end < rides[i].StartTime })

	if rightIdx == 0 {
		return 0, false
	}
	if leftIdx == rightIdx {
		return 0, false
	}

	var sum, cnt float64
	if m.concurrency <= 1 || rightIdx-leftIdx <= m.concurrency {
		sum, cnt = iterate(rides, end, leftIdx, rightIdx)
	} else {
		sum, cnt = m.iterateConc(rides, end, leftIdx, rightIdx)
	}

	if cnt == 0 {
		return 0, false
	}
	return sum / cnt, true
}

func (m *Binsearch) Init(files []reader.File) {
	m.metaFile = reader.File{}

	for _, file := range files {
		for cnt, rides := range file {
			m.metaFile[cnt] = append(m.metaFile[cnt], rides...)
		}
	}
	var cntRides int
	for _, rides := range m.metaFile {
		// Never trust the data!
		sort.Slice(rides, func(i, j int) bool {
			return rides[i].StartTime < rides[j].StartTime
		})
		cntRides += len(rides)
	}
	log.Printf("built index with %d rides\n", cntRides)
}

func (m *Binsearch) GetAverageDistances(start, end int64) map[int]float64 {
	result := make(map[int]float64)
	for cnt, rides := range m.metaFile {
		if avg, ok := m.avgDist(rides, start, end); ok {
			result[cnt] = avg
		}
	}
	return result
}
