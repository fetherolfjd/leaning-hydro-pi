package publish

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/tilt"
	"github.com/google/logger"
)

type CsvPublisher struct {
	csvFile   *os.File
	csvWriter *csv.Writer
}

func (p *CsvPublisher) Publish(tdp *tilt.TiltDataPoint) {
	logger.Infof("Writing a record to file %s", p.csvFile.Name())
	if writeErr := p.csvWriter.Write(tdp.StrVals()); writeErr != nil {
		logger.Errorf("Error writing to file: %v", writeErr)
	}
	p.csvWriter.Flush()
}

func (p *CsvPublisher) PublishAll(tdps []*tilt.TiltDataPoint) {
	logger.Infof("Writing %d records to file %s", len(tdps), p.csvFile.Name())
	data := make([][]string, len(tdps))
	for _, tdp := range tdps {
		data = append(data, tdp.StrVals())
	}
	p.csvWriter.WriteAll(data)
	p.csvWriter.Flush()
}

func (p *CsvPublisher) Close() error {
	logger.Infof("Closing publisher for file %s", p.csvFile.Name())
	p.csvWriter.Flush()
	cErr := p.csvFile.Close()
	if cErr != nil {
		return fmt.Errorf("error closing CSV file %s; %v", p.csvFile.Name(), cErr)
	}
	return nil
}

func NewCsvPublisher(destDir string, fileName string) DatapointPublisher {
	path := filepath.Join(destDir, fileName+".csv")
	logger.Infof("Attempting to open CSV file %x", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatalf("Error opening file %s; error: %v", path, err)
	}
	return &CsvPublisher{
		csvFile:   f,
		csvWriter: csv.NewWriter(f),
	}
}
