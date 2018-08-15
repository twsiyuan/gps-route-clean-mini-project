package main

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"
)

func TestCSVExport(t *testing.T) {
	count := rand.Int63n(10000) + 1
	ps := make([]GPSPoint, 0, count)
	for i := int64(0); i < count; i++ {
		ps = append(ps, GPSPoint{
			rand.Float64(),
			rand.Float64(),
			time.Unix(rand.Int63(), 0),
		})
	}

	r := NewGPSReaderRaws(ps)
	buf := bytes.NewBuffer(nil)
	if err := ExportAsCSV(r, buf); err != nil {
		t.Fatalf("Export failed, err %v", err)
	}

	r = NewGPSReader(bytes.NewReader(buf.Bytes()))
	for i := int64(0); i <= count; i++ {
		p, err := r.Read()
		if err != nil {
			if err == io.EOF && i == count {
				break
			}
			t.Fatalf("Read failed, %v", err)
		}
		if *p != ps[i] {
			t.Fatalf("Data is invalid")
		}
	}
}

func TestKMLExport(t *testing.T) {
	// TODO:
}
