//Credit: KonradIT https://github.com/KonradIT/gopro-utils

package main

import (
	"flag"
	"fmt"
	"github.com/JuanIrache/gopro-utils/telemetry"//getting rid of some bugs
	"io"
	"os"
	"strconv"
)

func main() {
	inName := flag.String("i", "", "Required: telemetry file to read")
	outName := flag.String("o", "", "Output kml map")
	accuracyThreshold := flag.Int("a", 1000, "Optional: GPS accuracy threshold, defaults to 1000")
	fixThreshold := flag.Int("f", 3, "Optional: GPS fix state. Defaults to 0 (no fix), can be 2 (2D) or 3 (3D)")
	flag.Parse()
	if *inName == "" {
		flag.Usage()
		return
	}
	if *outName == "" {
		flag.Usage()
		return
	}
	/*
		<?xml version="1.0" encoding="UTF-8"?>
		<kml xmlns="http://earth.google.com/kml/2.0">
		<Document>
		<Placemark>
		<Point><coordinates>Longitude,Latitude,Altitude</coordinates></Point>
		</Placemark>

		[LOOP]
		<Placemark>
		<Point><coordinates>LON,LAT,ALT</coordinates></Point>
		</Placemark>
		[/LOOP]

		</Document>
		</kml>
	*/
	var gpsData = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<kml xmlns=\"http://earth.google.com/kml/2.0\">\n<Document>\n<Placemark>\n<Point><coordinates>Longitude,Latitude,Altitude</coordinates></Point>\n</Placemark>\n"
	gpsFile, err := os.Create(*outName)
	gpsFile.WriteString(gpsData)
	defer gpsFile.Close()

	telemFile, err := os.Open(*inName)
	if err != nil {
		fmt.Printf("Cannot access telemetry file %s.\n", *inName)
		os.Exit(1)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Cannot close file %s: %s", file.Name(), err)
			os.Exit(1)
		}
	}(telemFile)

	// currently processing sentence
	t := &telemetry.TELEM{}

	for {
		t, err = telemetry.Read(telemFile)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		if t == nil {
			break
		}
		/*
			<Placemark>
			<Point><coordinates>LON,LAT,ALT</coordinates></Point>
			</Placemark>
		*/
		if t.GpsAccuracy.Accuracy < uint16(*accuracyThreshold) && t.GpsFix.F >= uint32(*fixThreshold) {
			for i, _ := range t.Gps {
				var TempGpsData string
				TempGpsData = "<Placemark>\n<Point><coordinates>" + floattostr(t.Gps[i].Longitude) + "," + floattostr(t.Gps[i].Latitude) + "," + floattostr(t.Gps[i].Altitude) + "</coordinates></Point>" + "\n</Placemark>\n"
				gpsFile.WriteString(TempGpsData)
			}
		}

		t = &telemetry.TELEM{}
	}
	gpsFile.WriteString("</Document>\n</kml>")

}

func floattostr(input_num float64) string {

	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', -1, 64)
}
