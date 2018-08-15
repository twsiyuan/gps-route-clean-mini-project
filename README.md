# GPS Route Clean Mini Project

## Problem

Given a series of points (latitude, longitude, timestamp) in CSV format for a journey from A-B, will disregard potentially erroneous points.

## Design

```
type GPSReader interface {
   Read() (*GPSPoint, error)
}
```

Use ```GPSReader``` interface to read points from the file, and filter potentially erroneous points by the specific algorithm. And easily extend for additional filters.
