package reader

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFiles(t *testing.T) {
	files, err := ReadFiles("../sampledata")
	assert.NoError(t, err)

	require.Len(t, files, 2)

	require.Len(t, files[0][1], 3)
	assert.Equal(t, files[0][1][0], Ride{StartTime: 1577838495, FinishTime: 1577838783, Distance: 1.2})

	assert.Len(t, files[1][1], 1)
}
