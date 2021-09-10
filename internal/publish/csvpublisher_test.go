package publish

import (
	"encoding/csv"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/tilt"
	"github.com/google/uuid"
)

func TestSingleValues(t *testing.T) {
	path := "/tmp"
	filename := "single-test"
	fullPath := filepath.Join(path, filename+".csv")
	rmFile(fullPath, t)
	publisher := NewCsvPublisher(path, filename)

	dp := createDataPoint()
	publisher.Publish(dp)

	publisher.Close()

	writtenRecords := readCsv(fullPath, t)
	compDataPoints([]*tilt.TiltDataPoint{dp}, writtenRecords, t)

	publisher = NewCsvPublisher(path, filename)
	dp2 := createDataPoint()
	publisher.Publish(dp2)
	publisher.Close()

	writtenRecords = readCsv(fullPath, t)
	compDataPoints([]*tilt.TiltDataPoint{dp, dp2}, writtenRecords, t)

	rmFile(fullPath, t)
}

func TestMultipleValues(t *testing.T) {
	path := "/tmp"
	filename := "multi-test"
	fullPath := filepath.Join(path, filename+".csv")
	rmFile(fullPath, t)
	publisher := NewCsvPublisher(path, filename)

	dps := []*tilt.TiltDataPoint{
		createDataPoint(),
		createDataPoint(),
		createDataPoint(),
		createDataPoint(),
		createDataPoint(),
	}
	publisher.PublishAll(dps)

	publisher.Close()

	writtenRecords := readCsv(fullPath, t)
	compDataPoints(dps, writtenRecords, t)

	publisher = NewCsvPublisher(path, filename)

	dps2 := []*tilt.TiltDataPoint{
		createDataPoint(),
		createDataPoint(),
		createDataPoint(),
		createDataPoint(),
		createDataPoint(),
	}
	publisher.PublishAll(dps2)

	publisher.Close()

	writtenRecords = readCsv(fullPath, t)
	checkDps := append(dps, dps2...)
	compDataPoints(checkDps, writtenRecords, t)

	rmFile(fullPath, t)
}

func compDataPoints(dps []*tilt.TiltDataPoint, csvData [][]string, t *testing.T) {
	if len(dps) != len(csvData) {
		t.Fatalf("Expected datapoint size %d is different than CSV size %d", len(dps), len(csvData))
	}

	for i := 0; i < len(dps); i++ {
		currDp := dps[i]
		currDpAsCsv := currDp.StrVals()
		currCsv := csvData[i]
		for j := 0; j < len(currDpAsCsv); j++ {
			d := currDpAsCsv[j]
			c := currCsv[j]
			if d != c {
				t.Fatalf("Expected value %s does not equal actual %s", d, c)
			}
		}
	}
}

func rmFile(filePath string, t *testing.T) {
	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			t.Fatalf("Unable to remove file %s", filePath)
		}
	}
}

func readCsv(filePath string, t *testing.T) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Unalbe to open file %s", filePath)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		t.Fatalf("Unable to parse file as CSV for %s", filePath)
	}
	return records
}

func createDataPoint() *tilt.TiltDataPoint {
	return &tilt.TiltDataPoint{
		UUID:            uuid.NewString(),
		Timestamp:       time.Now(),
		Address:         uuid.NewString(),
		RSSI:            rand.Int(),
		TXPower:         rand.Int(),
		Color:           "RED",
		SpecificGravity: rand.Float32(),
		Temperature:     rand.Float32(),
	}
}
