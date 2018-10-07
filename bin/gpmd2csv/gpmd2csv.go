package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"github.com/JuanIrache/gopro-utils/telemetry"	//linking to my own repository while the main one is behind. Not sure if this is a good practice

	//////used for csv
	"strconv"
    "log"
    "encoding/csv"

)

func main() {

	inName := flag.String("i", "", "Required: telemetry file to read")
	outName := flag.String("o", "", "Output csv files")
	flag.Parse()

	if *inName == "" {
		flag.Usage()
		return
	}

	///////////////////////////////////////////////////////////////////////////////////////////csv
	nameOut := string(*inName);
	if *outName != "" {
		nameOut = string(*outName);
	}
	////////////////////accelerometer
	var acclCsv = [][]string{{"Milliseconds","AcclX","AcclY","AcclZ"}}
	acclFile, err := os.Create(nameOut[:len(nameOut)-4]+"-accl.csv")
    checkError("Cannot create accl.csv file", err)
    defer acclFile.Close()
    acclWriter := csv.NewWriter(acclFile)
    /////////////////////gyroscope
    var gyroCsv = [][]string{{"Milliseconds","GyroX","GyroY","GyroZ"}}
	gyroFile, err := os.Create(nameOut[:len(nameOut)-4]+"-gyro.csv")
    checkError("Cannot create gyro.csv file", err)
    defer gyroFile.Close()
    gyroWriter := csv.NewWriter(gyroFile)
    //////////////////////temperature
    var tempCsv = [][]string{{"Milliseconds","Temp"}}
	tempFile, err := os.Create(nameOut[:len(nameOut)-4]+"-temp.csv")
    checkError("Cannot create temp.csv file", err)
    defer tempFile.Close()
    tempWriter := csv.NewWriter(tempFile)
    ///////////////////////Uncomment for Gps
    var gpsCsv = [][]string{{"Milliseconds","Latitude","Longitude","Altitude","Speed","Speed3D","TS"}}
	gpsFile, err := os.Create(nameOut[:len(nameOut)-4]+"-gps.csv")
    checkError("Cannot create gps.csv file", err)
    defer gpsFile.Close()
    gpsWriter := csv.NewWriter(gpsFile)
    //////////////////////////////////////////////////////////////////////////////////////////////

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
	t_prev := &telemetry.TELEM{}

	seconds := -1
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

		// first full, guess it's about a second
		if t_prev.IsZero() {
			*t_prev = *t
			t.Clear()
			continue
		}

		// process until t.Time
		t_prev.FillTimes(t.Time.Time)

		// this is pretty useless and info overload: change it to pick a field you want
		// or mangle it to your wishes into JSON/CSV/format of choice
		// fmt.Println(t_prev)

		///////////////////////////////////////////////////////////////////Modified to save CSV
		/////////////////////Accelerometer
	    for i, _ := range t_prev.Accl {
	    	milliseconds := float64(seconds*1000)+float64(((float64(1000)/float64(len(t_prev.Accl)))*float64(i)))
			acclCsv = append(acclCsv, []string{floattostr(milliseconds),floattostr(t_prev.Accl[i].X),floattostr(t_prev.Accl[i].Y),floattostr(t_prev.Accl[i].Z)})
		}
		/////////////////////Gyroscope
	    for i, _ := range t_prev.Gyro {
	    	milliseconds := float64(seconds*1000)+float64(((float64(1000)/float64(len(t_prev.Gyro)))*float64(i)))
			gyroCsv = append(gyroCsv, []string{floattostr(milliseconds),floattostr(t_prev.Gyro[i].X),floattostr(t_prev.Gyro[i].Y),floattostr(t_prev.Gyro[i].Z)})
		}
		////////////////////Temperature
		milliseconds := seconds*1000
		tempCsv = append(tempCsv, []string{strconv.Itoa(milliseconds),floattostr(float64(t_prev.Temp.Temp))})
		////////////////////Uncomment for Gps
		for i, _ := range t_prev.Gps {
			fmt.Println(t_prev.Gps[i].TS)
			len := len(t_prev.Gps)
			milliseconds := float64(seconds*1000)+float64(((float64(1000)/float64(len))*float64(i)))
			gpsCsv = append(gpsCsv, []string{floattostr(milliseconds),floattostr(t_prev.Gps[i].Latitude),floattostr(t_prev.Gps[i].Longitude),floattostr(t_prev.Gps[i].Altitude),floattostr(t_prev.Gps[i].Speed),floattostr(t_prev.Gps[i].Speed3D),int64tostr(t_prev.Gps[i].TS)})
		}
	    //////////////////////////////////////////////////////////////////////////////////
		
		*t_prev = *t
		t = &telemetry.TELEM{}
		seconds++
	}
	/////////////////////////////////////////////////////////////////////////////////////for csv
	///////////////accelerometer
	for _, value := range acclCsv {
        err := acclWriter.Write(value)
        checkError("Cannot write to accl.csv file", err)
    }
    defer acclWriter.Flush()
    ///////////////gyroscope
    for _, value := range gyroCsv {
        err := gyroWriter.Write(value)
        checkError("Cannot write to gyro.csv file", err)
    }
    defer gyroWriter.Flush()
    /////////////temperature
    for _, value := range tempCsv {
        err := tempWriter.Write(value)
        checkError("Cannot write to temp.csv file", err)
    }
    defer tempWriter.Flush()
    /////////////Uncomment for Gps
    for _, value := range gpsCsv {
        err := gpsWriter.Write(value)
        checkError("Cannot write to gps.csv file", err)
    }
    defer gpsWriter.Flush()
    /////////////////////////////////////////////////////////////////////////////////////
}


///////////for csv

func floattostr(input_num float64) string {

        // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', -1, 64)
}



func int64tostr(input_num int64) string {

        // to convert a float number to a string
    return strconv.FormatInt(input_num, 10)
}

 func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}

