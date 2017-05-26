// Extracts stats from GoPro Metadata
/* Essential to have:
	- Mean (! (vf-vi)/(tf-ti) ) acceleration on x, y and z DONE
	- Distance travelled
	- Mean speed DONE
	- Mean altitude DONE
	- Peak altitude DONE
	- Peak speed DONE
	- Peak acceleration
	- G-Force
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
    "log"
	"math"
	"strings"
	"github.com/stilldavid/gopro-utils/telemetry"
)

func main() {

	inName := flag.String("i", "", "Required: telemetry file to read")
	flag.Parse()
	
	if *inName == "" {
		flag.Usage()
		return
	}

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
	//Arrays
	
	Gyroscope_X := []float64{}
	Gyroscope_Y := []float64{}
	Gyroscope_Z := []float64{}
	Speed := []float64{}
	Altitude := []float64{}
	
	Accel_X := []float64{}
	Accel_Y := []float64{}
	Accel_Z := []float64{}
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

		for i, _ := range t.Accl {
	    	//milliseconds := float64(seconds*1000)+float64(((float64(1000)/float64(len(t.Accl)))*float64(i)))
			Accel_X = append(Accel_X, t.Accl[i].X)
			Accel_Y = append(Accel_Y, t.Accl[i].Y)
			Accel_Z = append(Accel_Z, t.Accl[i].Z)
		}
		
	    for i, _ := range t.Gyro {
	    	//milliseconds := float64(seconds*1000)+float64(((float64(1000)/float64(len(t.Gyro)))*float64(i)))
			Gyroscope_X = append(Gyroscope_X, t.Gyro[i].X)
			Gyroscope_Y = append(Gyroscope_Y, t.Gyro[i].Y)
			Gyroscope_Z = append(Gyroscope_Z, t.Gyro[i].Z)
		}	
		
		for i, _ := range t.Gps {
			//gpsCsv = append(gpsCsv, []string{floattostr(t.Gps[i].Latitude),floattostr(t.Gps[i].Longitude),floattostr(t.Gps[i].Altitude),floattostr(t.Gps[i].Speed),floattostr(t.Gps[i].Speed3D),int64tostr(t.Gps[i].TS)})
			Speed = append(Speed, t.Gps[i].Speed)
			Altitude = append(Altitude, t.Gps[i].Altitude)
		}
		
			
		
		t = &telemetry.TELEM{}
		seconds++
	}
	var total_accel_x, total_accel_y, total_accel_z, speed, altitude float64 = 0, 0, 0, 0, 0
	for _, value:= range Accel_X {
		total_accel_x += value
	}
	for _, value:= range Accel_Y {
		total_accel_y += value
	}
	for _, value:= range Accel_Z {
		total_accel_z += value
	}
	for _, value:= range Speed {
		speed += value
	}
	for _, value:= range Altitude {
		altitude += value
	}
	fmt.Println("Data from " + *inName)
	fmt.Println("Average acceleration on X axis: " + floattostr(Round(total_accel_x/float64(len(Accel_X)), .5, 3)))
	fmt.Println("Average acceleration on Y axis: " + floattostr(Round(total_accel_y/float64(len(Accel_Y)), .5, 3)))
	fmt.Println("Average acceleration on Z axis: " + floattostr(Round(total_accel_z/float64(len(Accel_Z)), .5, 3)))
	fmt.Println("===============================")
	fmt.Println("Average speed: " + floattostr(Round(speed/float64(len(Speed)), .5, 3)) + " m/s")
	fmt.Println("Average altitude: " + floattostr(Round(altitude/float64(len(Altitude)), .5, 3)) + " meters")
	fmt.Println("Peak altitude: " + floattostr(getBiggestFromSlice(Altitude)) + " meters")
	fmt.Println("Peak speed: " + floattostr(getBiggestFromSlice(Speed)) + " m/s")
}

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
func Round(val float64, roundOn float64, places int ) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func Ms2hms(ms float64) string {
	xstr := floattostr(ms)
	i := strings.Index(xstr, ".")
	decimal := xstr[i+1:]
	if len(xstr[i+1:]) < 3{
		decimal = xstr[i+1:] + "0"
	}
	x := ms
	seconds := int(ms)
	x = x / 60
	minutes := int(math.Mod(x, 60))
	x = x / 60
	hours := int(math.Mod(x, 24))
	s := ""
	m := ""
	h := ""
	if len(strconv.Itoa(hours)) == 1{
		h = "0"
	} 
	if len(strconv.Itoa(minutes)) == 1{
		m = "0"
	}
	if len(strconv.Itoa(seconds)) == 1{
		s = "0"
	}
	return h + strconv.Itoa(hours) + ":" + m + strconv.Itoa(minutes) + ":" + s + strconv.Itoa(seconds) + "." + decimal
}
func getBiggestFromSlice(slice []float64) float64 {
	var n, biggest float64
	for _,v:=range slice {
    if v>n {
      n = v
      biggest = n
    }
	}
	return biggest
}
func getSmallestFromSlice(slice []float64) float64 {
	var n, smallest float64
	for _,v:=range slice {
    if v<n {
      n = v
      smallest = n
    }
	}
	return smallest
}