package main

import (
	"flag"
	"fmt"
	"github.com/stilldavid/gopro-utils/telemetry"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	inName := flag.String("i", "", "Required: telemetry file to read")
	outName := flag.String("o", "", "Required: output marker file")
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
		Gets the top X highest altitude/speed
	*/
	var markerData = ""
	markerFile, err := os.Create(*outName)
	markerFile.WriteString(markerData)
	defer markerFile.Close()

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
	seconds := -1
	Speed := []float64{}
	Altitude := []float64{}
	Epochtime := []string{}
	BiggestSpeed, BiggestAltitude := 0.0, 0.0
	BiggestSpeedTime, BiggestAltitudeTime := "", ""
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
			[MARKER START]\t[MARKER END]\t[TEXT]
		*/
		t.FillTimes(t.Time.Time)
		for i, _ := range t.Gps {
			Speed = append(Speed, t.Gps[i].Speed)
			if t.Gps[i].Speed > BiggestSpeed {
				BiggestSpeed = t.Gps[i].Speed
				BiggestSpeedTime = int64tostr(t.Gps[i].TS)
			}
			Altitude = append(Altitude, t.Gps[i].Altitude)
			if t.Gps[i].Altitude > BiggestAltitude {
				BiggestAltitude = t.Gps[i].Altitude
				BiggestAltitudeTime = int64tostr(t.Gps[i].TS)
			}
			Epochtime = append(Epochtime, int64tostr(t.Gps[i].TS))

		}
		t = &telemetry.TELEM{}
		seconds++
	}
	var speed, altitude float64 = 0, 0
	for _, value := range Speed {
		speed += value
	}
	for _, value := range Altitude {
		altitude += value
	}
	markerFile.WriteString(getDifferenceBetweenDates(getUTCTimeFromUnix(Epochtime[0]), getUTCTimeFromUnix(BiggestAltitudeTime)) + "\t" + getDifferenceBetweenDates(getUTCTimeFromUnix(Epochtime[0]), getUTCTimeFromUnix(BiggestAltitudeTime)) + "\t" + floattostr(getBiggestFromSlice(Altitude)) + " m")
	markerFile.WriteString("\n" + getDifferenceBetweenDates(getUTCTimeFromUnix(Epochtime[0]), getUTCTimeFromUnix(BiggestSpeedTime)) + "\t" + getDifferenceBetweenDates(getUTCTimeFromUnix(Epochtime[0]), getUTCTimeFromUnix(BiggestSpeedTime)) + "\t" + floattostr(getBiggestFromSlice(Speed)) + " m/s")
}

func floattostr(input_num float64) string {

	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', -1, 64)
}

func getDifferenceBetweenDates(date_1 string, date_2 string) string {
	var year, _ = strconv.Atoi(date_1[0:4])
	var month, _ = strconv.Atoi(date_1[5:7])
	var day, _ = strconv.Atoi(date_1[8:10])
	var hour, _ = strconv.Atoi(date_1[11:13])
	var minute, _ = strconv.Atoi(date_1[13:16])
	var second, _ = strconv.Atoi(date_1[17:19])
	var decimals, _ = strconv.Atoi(date_1[31:35])
	start := time.Date(
		year, time.Month(month), day, hour, minute, second, decimals, time.UTC)
	var year_2, _ = strconv.Atoi(date_2[0:4])
	var month_2, _ = strconv.Atoi(date_2[5:7])
	var day_2, _ = strconv.Atoi(date_2[8:10])
	var hour_2, _ = strconv.Atoi(date_2[11:13])
	var minute_2, _ = strconv.Atoi(date_2[13:16])
	var second_2, _ = strconv.Atoi(date_2[17:19])
	var decimals_2, _ = strconv.Atoi(date_2[31:35])
	big := time.Date(
		year_2, time.Month(month_2), day_2, hour_2, minute_2, second_2, decimals_2, time.UTC)
	diff := start.Sub(big)
	final := strconv.Itoa(int(diff.Seconds()))
	return final
}
func getUTCTimeFromUnix(timestamp string) string {
	i, err := strconv.ParseInt(timestamp[0:10], 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0).String()
	return tm + "-" + timestamp[11:16]
}
func int64tostr(input_num int64) string {

	// to convert a float number to a string
	return strconv.FormatInt(input_num, 10)
}
func getBiggestFromSlice(slice []float64) float64 {
	var n, biggest float64
	for _, v := range slice {
		if v > n {
			n = v
			biggest = n
		}
	}
	return biggest
}
