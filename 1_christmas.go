package main

import (
	"fmt"
	"math"
	"time"
)

func main() {

	//NOTE :: timezone will be used of where the program is being run.

	christmas, err := time.Parse("Jan 2 06", "Dec 25 21")
	if err != nil {
		panic("Could not parse the date, error : " + err.Error())
	}

	fmt.Printf("Milliseconds to christmas: %v\n", christmas.Sub(time.Now()).Milliseconds())

	//%.f will truncate the decimals, we dont need them for printing purpose.
	fmt.Printf("Seconds to christmas: %.f\n", christmas.Sub(time.Now()).Seconds())
	fmt.Printf("Minutes to christmas: %.f\n", christmas.Sub(time.Now()).Minutes())
	fmt.Printf("Hours to christmas: %.f\n", christmas.Sub(time.Now()).Hours())
	fmt.Printf("Days to christmas: %.f\n", math.Round(christmas.Sub(time.Now()).Hours()/24))
}
