package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"time"
)

var (
	errUnexpectedFormat = errors.New("unexpected format")
)

// GPSPoint is the struct that present GPS point
type GPSPoint struct {
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

// GPSReader is the interface that use read GPSPoints
//
// When Read encounters end-of-file condition, it returns io.EOF
type GPSReader interface {
	Read() (*GPSPoint, error)
}

// NewGPSReader returns a new Reader that reads GPS points in CSV format.
func NewGPSReader(r io.Reader) GPSReader {
	return &gpsReader{
		reader: csv.NewReader(r),
	}
}

// NewGPSReaderRaws returns a new Reader that reads GPS points from array/slice
func NewGPSReaderRaws(raws []GPSPoint) GPSReader {
	return &gpsReaderRaws{raws}
}

type gpsReader struct {
	reader *csv.Reader
}

func (r *gpsReader) Read() (*GPSPoint, error) {

	record, err := r.reader.Read()
	if err != nil {
		return nil, err
	}

	if len(record) != 3 {
		return nil, errUnexpectedFormat
	}

	point, err := r.parsePoint(record)
	if err != nil {
		return nil, errUnexpectedFormat
	}

	return point, nil
}

func (r gpsReader) parsePoint(record []string) (*GPSPoint, error) {
	latitude, err := strconv.ParseFloat(record[0], 64)
	if err != nil {
		return nil, err
	}

	longitude, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}

	timestamp, err := strconv.ParseInt(record[2], 10, 64)
	if err != nil {
		return nil, err
	}

	return &GPSPoint{
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Unix(timestamp, 0),
	}, nil
}

type gpsReaderRaws struct {
	raws []GPSPoint
}

func (r *gpsReaderRaws) Read() (*GPSPoint, error) {
	if len(r.raws) <= 0 {
		return nil, io.EOF
	}
	p := r.raws[0]
	r.raws = r.raws[1:]
	return &p, nil
}
