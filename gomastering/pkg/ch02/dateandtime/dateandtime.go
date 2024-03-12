package dateandtime

import (
	"fmt"
	"os"
	"time"
)

func WorkingWithDates() {
	os.Args = append(os.Args, "12 March 2024 15:04")

	fmt.Println(time.Now().Format("01-02-2006 15:04"))
	fmt.Println(1)

	start := time.Now()

	if len(os.Args) == 1 {
		fmt.Println("Usage: dateandtime parse_string")
		return
	}

	dateStr := os.Args[1]
	if d, err := time.Parse("02 January 2006", dateStr); err == nil {
		fmt.Println("Full date:", d)
		fmt.Println("Time:", d.Day(), d.Month(), d.Year())
		fmt.Println(2)
	}

	if d, err := time.Parse("02 January 2006 15:04", dateStr); err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Date: ", d.Day(), d.Month(), d.Year())
		fmt.Println("Time: ", d.Hour(), d.Minute())
		fmt.Println(3)
	}

	if d, err := time.Parse("02-01-2006 15:04", dateStr); err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Date:", d.Day(), d.Month(), d.Year())
		fmt.Println("Time:", d.Hour(), d.Minute())
		fmt.Println(4)
	}

	if d, err := time.Parse("02/01/2006 1504", dateStr); err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Date:", d.Day(), d.Month(), d.Year())
		fmt.Println("Time:", d.Hour(), d.Minute())
		fmt.Println(5)
	}

	if d, err := time.Parse("15:04", dateStr); err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Time: ", d.Hour(), d.Minute())
		fmt.Println(6)
	}

	t := time.Now().Unix()
	fmt.Println("Epoch time:", t)
	fmt.Println(7)

	d := time.Unix(t, 0)

	fmt.Println("Date:", d.Day(), d.Month(), d.Year())
	fmt.Printf("Time: %d:%d\n", d.Hour(), d.Minute())
	fmt.Println(8)

	duration := time.Since(start)
	fmt.Println("Execution time:", duration)
	fmt.Println(9)
}

func LocationParsing() {
	os.Args = append(os.Args, "12 March 2024 17:09 UTC")
	if len(os.Args) == 1 {
		fmt.Println("Need input: _02 January 2006 15:04 MST_ format")
		return
	}

	input := os.Args[1]
	d, err := time.Parse("02 January 2006 15:04 MST", input)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Local
	loc, _ := time.LoadLocation("Local")
	fmt.Println("Current location:", d.In(loc))

	//NY
	loc, _ = time.LoadLocation("America/New_York")
	fmt.Println("New York Time:", d.In(loc))

	// London
	loc, _ = time.LoadLocation("Europe/London")
	fmt.Println("London Time:", d.In(loc))

	// Tokyo
	loc, _ = time.LoadLocation("Asia/Tokyo")
	fmt.Println("Tokyo Time:", d.In(loc))

}
