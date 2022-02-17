package index

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"rides/reader"
	"testing"
	"time"
)

func TestIndexes(t *testing.T) {
	files := []reader.File{
		{
			1: []reader.Ride{
				{StartTime: 1, FinishTime: 2, Distance: 1},
				{StartTime: 1, FinishTime: 3, Distance: 2},
				{StartTime: 2, FinishTime: 3, Distance: 6},
			},
			2: []reader.Ride{
				{StartTime: 1, FinishTime: 2, Distance: 10},
				{StartTime: 1, FinishTime: 3, Distance: 20},
				{StartTime: 2, FinishTime: 3, Distance: 60},
			},
		},
	}
	indexes := []Index{
		&Dummy{},
		NewBinsearch(0),
		NewBinsearch(1),
		NewBinsearch(2),
	}
	for _, index := range indexes {
		index.Init(files)

		assert.Equal(t, map[int]float64{1: 1, 2: 10}, index.GetAverageDistances(1, 2), index)
		assert.Equal(t, map[int]float64{1: 3, 2: 30}, index.GetAverageDistances(1, 3), index)
		assert.Equal(t, map[int]float64{1: 6, 2: 60}, index.GetAverageDistances(2, 3), index)

		assert.Equal(t, map[int]float64{}, index.GetAverageDistances(4, 4), index)
	}
}

func randomRides(maxTime int64) []reader.Ride {
	var rides []reader.Ride
	for i := 0; i < int(maxTime); i++ {
		t1 := rand.Int63n(maxTime)
		t2 := rand.Int63n(maxTime)
		if t2 < t1 {
			t2, t1 = t1, t2
		}

		rides = append(rides, reader.Ride{
			StartTime:  t1,
			FinishTime: t2,
			Distance:   rand.Float64(),
		})
	}
	return rides
}

func randomFiles(maxTime int64) []reader.File {
	var files []reader.File
	for i := 0; i < 10; i++ {
		file := reader.File{}
		for j := 0; j < 10; j++ {
			file[j] = randomRides(maxTime)
		}
		files = append(files, file)
	}
	return files
}

func stressTest(t *testing.T, index1, index2 Index) {
	maxTime := int64(1000)
	files := randomFiles(maxTime)
	index1.Init(files)
	index2.Init(files)
	for i := 0; i < 1000; i++ {
		t1 := rand.Int63n(maxTime)
		t2 := rand.Int63n(maxTime)
		if t2 < t1 {
			t2, t1 = t1, t2
		}

		ans1 := index1.GetAverageDistances(t1, t2)
		ans2 := index2.GetAverageDistances(t1, t2)
		assert.Equal(t, len(ans1), len(ans2))
		for k, v := range ans1 {
			assert.InDelta(t, v, ans2[k], 1e-9, "t1: %d, t2: %d, k: %d", t1, t2, k)
		}
	}
}

func TestStress(t *testing.T) {
	stressTest(t, &Dummy{}, &Binsearch{})
	stressTest(t, &Dummy{}, NewBinsearch(2))
}

func randomBenchmark(name string, maxTime int64, index Index) {
	files := randomFiles(maxTime)
	index.Init(files)
	start := time.Now()
	for i := 0; i < 1000; i++ {
		t1 := rand.Int63n(maxTime)
		t2 := rand.Int63n(maxTime)
		if t2 < t1 {
			t2, t1 = t1, t2
		}

		index.GetAverageDistances(t1, t2)
	}
	dur := time.Since(start)
	fmt.Printf("%s took %s with maxTime=%d\n", name, dur, maxTime)
}

func TestSpeedSmall(t *testing.T) {
	maxTime := int64(1000)
	randomBenchmark("Dummy", maxTime, &Dummy{})
	randomBenchmark("Binsearch-1", maxTime, NewBinsearch(1))
	randomBenchmark("Binsearch-4", maxTime, NewBinsearch(4))
	randomBenchmark("Binsearch-8", maxTime, NewBinsearch(8))
}

func TestSpeedMedium(t *testing.T) {
	maxTime := int64(10000)
	randomBenchmark("Dummy", maxTime, &Dummy{})
	randomBenchmark("Binsearch-1", maxTime, NewBinsearch(1))
	randomBenchmark("Binsearch-4", maxTime, NewBinsearch(4))
	randomBenchmark("Binsearch-8", maxTime, NewBinsearch(8))
}
