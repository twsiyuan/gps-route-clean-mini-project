package main

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"
)

func testRead(t *testing.T, r GPSReader, latitude, longitude float64, timestamp int64) {
	if p, err := r.Read(); err != nil {
		t.Error(err)
	} else if v := p.Latitude; v != latitude {
		t.Errorf("Latitude is invalid, read %f", v)
	} else if v := p.Longitude; v != longitude {
		t.Errorf("Longitude is invalid, read %f", v)
	} else if v := p.Timestamp.Unix(); v != timestamp {
		t.Errorf("Timestamp is invalid, read %d", v)
	}
}

func TestRead(t *testing.T) {
	br := bytes.NewReader(([]byte)(`0,1,2
0.1,0.2,3

-0.1,-0.2,4
`))

	r := NewGPSReader(br)
	testRead(t, r, 0, 1, 2)
	testRead(t, r, 0.1, 0.2, 3)
	testRead(t, r, -0.1, -0.2, 4)

	if _, err := r.Read(); err != io.EOF {
		t.Errorf("Expected io.EOF, got: %v", err)
	}
	if _, err := r.Read(); err != io.EOF {
		t.Errorf("Expected io.EOF, got: %v", err)
	}
}

func TestReadError(t *testing.T) {
	br := bytes.NewReader(([]byte)(`0,0,0,`))
	r := NewGPSReader(br)
	if _, err := r.Read(); err == nil {
		t.Errorf("Expected error")
	}

	br = bytes.NewReader(([]byte)(`0`))
	r = NewGPSReader(br)
	if _, err := r.Read(); err == nil {
		t.Errorf("Expected error")
	}

	br = bytes.NewReader(([]byte)(`1,1,a`))
	r = NewGPSReader(br)
	if _, err := r.Read(); err == nil {
		t.Errorf("Expected error")
	}

	br = bytes.NewReader(([]byte)(`1,a,1`))
	r = NewGPSReader(br)
	if _, err := r.Read(); err == nil {
		t.Errorf("Expected error")
	}

	br = bytes.NewReader(([]byte)(`a,1,1`))
	r = NewGPSReader(br)
	if _, err := r.Read(); err == nil {
		t.Errorf("Expected error")
	}

	br = bytes.NewReader(([]byte)(`0,1,1.1`))
	r = NewGPSReader(br)
	if _, err := r.Read(); err == nil {
		t.Errorf("Expected error")
	}
}

func TestReadRaws(t *testing.T) {
	const count = 30
	ps := make([]GPSPoint, 0, count)
	for i := 0; i < count; i++ {
		ps = append(ps, GPSPoint{
			rand.Float64(),
			rand.Float64(),
			time.Unix(rand.Int63(), 0),
		})
	}

	r := NewGPSReaderRaws(ps)
	for i := 0; i <= count; i++ {
		p, err := r.Read()
		if err != nil {
			if err == io.EOF && i == count {
				break
			}
			t.Fatal(err)
		}
		if *p != ps[i] {
			t.Fatalf("Data is invalid")
		}
	}
}
