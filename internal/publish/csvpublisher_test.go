package publish

import (
	"encoding/csv"
	"os"
	"path"
	"testing"
	"time"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCsvPublisher(t *testing.T) {

	tempDir := t.TempDir()
	csvFile := path.Join(tempDir, "testfile.csv")

	r := model.TiltReading{
		UUID:            "test",
		SpecificGravity: 1,
		Temperature:     2,
		RSSI:            3,
		TXPower:         4,
		Address:         "testaddr",
		Color:           "testcolor",
		Timestamp:       time.Now(),
	}
	strSli := stringify(r)

	csvPub, err := NewCsvPublisher(csvFile, hclog.NewNullLogger())
	require.NoError(t, err)

	err = csvPub.Publish(r)
	require.NoError(t, err)

	err = csvPub.PublishAll([]model.TiltReading{r, r})
	require.NoError(t, err)

	err = csvPub.Close()
	require.NoError(t, err)

	csvF, err := os.Open(csvFile)
	defer csvF.Close()
	require.NoError(t, err)
	csvRead := csv.NewReader(csvF)
	recs, err := csvRead.ReadAll()
	require.NoError(t, err)
	assert.Len(t, recs, 3)
	for _, rec := range recs {
		assert.Equal(t, strSli, rec)
	}
}
