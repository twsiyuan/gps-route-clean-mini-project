package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"strconv"
)

// ExportAsCSV is function that exports all gps point to w in CSV format
func ExportAsCSV(r GPSReader, w io.Writer) error {
	csvw := csv.NewWriter(w)
	for true {
		p, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		record := []string{
			strconv.FormatFloat(p.Latitude, 'g', -1, 64),
			strconv.FormatFloat(p.Longitude, 'g', -1, 64),
			strconv.FormatInt(p.Timestamp.Unix(), 10),
		}
		csvw.Write(record)
	}
	csvw.Flush()
	return nil
}

// ExportAsKML is function that exports all gps point to w in KML format
//
// TODO: Use xml library
func ExportAsKML(r GPSReader, w io.Writer) error {
	writer := bufio.NewWriter(w)
	writer.WriteString(`
<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <Placemark>
      <LineString>
        <tessellate>1</tessellate>
		<coordinates>`)

	for true {
		p, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		writer.WriteString("\n")
		writer.WriteString(strconv.FormatFloat(p.Longitude, 'g', -1, 64))
		writer.WriteString(",")
		writer.WriteString(strconv.FormatFloat(p.Latitude, 'g', -1, 64))
		writer.WriteString(",0")
	}
	writer.WriteString(`
	</coordinates>
      </LineString>
    </Placemark>
  </Document>
</kml>`)

	return writer.Flush()
}
