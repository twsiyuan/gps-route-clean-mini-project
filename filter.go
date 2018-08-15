package main

import (
	"io"
	"math"
)

// NewSpeedFilter return reader that performs speed-limit algorithm
func NewSpeedFilter(r GPSReader, limitMph float64) GPSReader {
	return &speedFilter{
		r:        r,
		limitMph: limitMph,
	}
}

// NewNoiseFilter return reader that performs de-noise algorithm
func NewNoiseFilter(r GPSReader) GPSReader {
	return &noiseFilter{r: r}
}

func geoPoint(p *GPSPoint) Vector2 {
	// SEE: https://stackoverflow.com/questions/14344207/how-to-convert-distancemiles-to-degrees
	return Vector2{
		X: p.Longitude * math.Cos(p.Latitude) * 69.17101972,
		Y: p.Latitude * 68.70747695,
	}
}

func geoVelocity(p1, p2 *GPSPoint) Vector2 {
	duration := geoDuration(p1, p2)
	geo1 := geoPoint(p1)
	geo2 := geoPoint(p2)
	return geo2.Sub(geo1).Multiply(float64(1) / duration)
}

func geoDuration(p1, p2 *GPSPoint) float64 {
	return p2.Timestamp.Sub(p1.Timestamp).Hours()
}

type speedFilter struct {
	r        GPSReader
	previous *GPSPoint
	limitMph float64
}

func (f *speedFilter) Read() (*GPSPoint, error) {
	for {
		p, err := f.r.Read()
		if err != nil {
			return nil, err
		}

		q, passed := f.process(p)
		if passed {
			return q, nil
		}
	}
}

func (f *speedFilter) updatePrevious(p *GPSPoint) {
	pp := *p
	f.previous = &pp
}

func (f *speedFilter) process(current *GPSPoint) (result *GPSPoint, ok bool) {
	defer func() {
		if ok {
			f.updatePrevious(current)
		}
	}()

	previous := f.previous
	if previous == nil {
		return current, true
	}

	speedMph := geoVelocity(previous, current).Magnitude()
	if speedMph <= 0.001 {
		// No movement
		return nil, false
	} else if speedMph > f.limitMph {
		// Too fast
		return nil, false
	}
	return current, true
}

type noiseFilter struct {
	r      GPSReader
	points []*GPSPoint
}

func (f *noiseFilter) Read() (*GPSPoint, error) {
	// Keep streaming possible here
	points := f.points
	for len(points) < 3 {
		p, err := f.r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	if n := len(points); n <= 0 {
		return nil, io.EOF
	} else if n >= 3 {
		p1 := geoPoint(points[n-3])
		p2 := geoPoint(points[n-2])
		p3 := geoPoint(points[n-1])

		dir1 := p2.Sub(p1).Normalize()
		dir2 := p3.Sub(p2).Normalize()

		// Use directions of p3p2 and p2p1,
		// If the angle of this two vectors is over 165 degrees,
		// means it might be the noise pattern (from observation),
		// So remove those points p3 and p2
		if c := dir2.Dot(dir1); c <= -0.97 {
			points = points[:n-2]
		}
	}

	p := points[0]
	f.points = points[1:]
	return p, nil
}
