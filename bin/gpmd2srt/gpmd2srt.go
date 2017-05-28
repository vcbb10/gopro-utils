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
	srtName := flag.String("o", "", "Required: output srt file")
	flag.Parse()
	
	if *inName == "" {
		flag.Usage()
		return
	}
	if *srtName == "" {
		flag.Usage()
		return
	}
	file, err := os.Create(*srtName)
    checkError("Cannot create " + *srtName + " file", err)
    defer file.Close()

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
	count := 0
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
			
	    	milliseconds := Round(float64(seconds*1000)+float64(((float64(1000)/float64(len(t.Accl)))*float64(i))), .5, 0) / 1000
			next_milliseconds := Round(float64(seconds*1000)+float64(((float64(1000)/float64(len(t.Accl)))*float64(i+1))), .5, 0) / 1000
			current := ""
			next := ""
			if int(milliseconds) != 0{
			count = count + 1
			current = Ms2hms(milliseconds)
			next = Ms2hms(next_milliseconds)
			text_str := "X: " + floattostr(t.Accl[i].X) + "\nY: " + floattostr(t.Accl[i].Y) + "\nZ: " + floattostr(t.Accl[i].Z)
			file.WriteString("\n" + strconv.Itoa(count) + "\n")
			file.WriteString(current + " --> " + next + "\n")
			file.WriteString(text_str + "\n")
			}
			
			
		}
		t = &telemetry.TELEM{}
		seconds++
	}
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
