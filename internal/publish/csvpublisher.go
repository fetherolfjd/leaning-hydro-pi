package publish

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/hashicorp/go-hclog"
)

type CsvPublisher struct {
	csvFile   *os.File
	csvWriter *csv.Writer
	logger    hclog.Logger
}

func stringify(r model.TiltReading) []string {
	return []string{
		fmt.Sprintf("%d", r.Timestamp.UnixMilli()),
		r.Color,
		r.UUID,
		r.Address,
		fmt.Sprintf("%1.4f", r.SpecificGravity),
		fmt.Sprintf("%.1f", r.Temperature),
		fmt.Sprintf("%d", r.RSSI),
		fmt.Sprintf("%d", r.TXPower),
	}
}

func (p *CsvPublisher) Publish(r model.TiltReading) error {
	p.logger.Debug("writing a record to file %s", p.csvFile.Name())
	if err := p.csvWriter.Write(stringify(r)); err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}
	p.csvWriter.Flush()
	if err := p.csvWriter.Error(); err != nil {
		return fmt.Errorf("unable to flush to file: %w", err)
	}
	return nil
}

func (p *CsvPublisher) PublishAll(rs []model.TiltReading) error {
	p.logger.Debug("writing %d records to file %s", len(rs), p.csvFile.Name())
	data := make([][]string, 0, len(rs))
	for _, r := range rs {
		data = append(data, stringify(r))
	}
	if err := p.csvWriter.WriteAll(data); err != nil {
		return fmt.Errorf("unable to write multiple records to file: %w", err)
	}
	p.csvWriter.Flush()
	if err := p.csvWriter.Error(); err != nil {
		return fmt.Errorf("unable to flush multiple records to file: %w", err)
	}
	return nil
}

func (p *CsvPublisher) Close() error {
	p.logger.Debug("Closing publisher for file %s", p.csvFile.Name())
	p.csvWriter.Flush()
	if err := p.csvFile.Close(); err != nil {
		return fmt.Errorf("error closing CSV file %s; %v", p.csvFile.Name(), err)
	}
	return nil
}

func NewCsvPublisher(filePath string, logger hclog.Logger) (*CsvPublisher, error) {
	logger.Debug("Attempting to open CSV file %s", filePath)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open CSV file: %s; %w", filePath, err)
	}
	return &CsvPublisher{
		csvFile:   f,
		csvWriter: csv.NewWriter(f),
		logger:    logger,
	}, nil
}
