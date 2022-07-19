package main

import (
	"fmt"
	"strconv"
)

type Duration struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

func (d *Duration) FromString(durationString string) error {
	var numberStr []rune

	isNumber := false

	for _, char := range durationString {
		switch {
		case (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z'):
			if !isNumber {
				return fmt.Errorf("invalid format")
			}

			isNumber = false

			number, _ := strconv.Atoi(string(numberStr))
			unit := string(char)
			switch unit {
			case "s":
				d.Second = number
			case "m":
				d.Minute = number
			case "h":
				d.Hour = number
			case "d":
				d.Day = number
			case "M":
				d.Month = number
			case "y":
				d.Year = number
			default:
				return fmt.Errorf("unknown time unit")
			}

			numberStr = []rune{}

		case char >= '0' && char <= '9':
			isNumber = true

			numberStr = append(numberStr, char)
		}
	}

	return nil
}

func (d Duration) GetSeconds() int {
	return d.Second +
		d.Minute*60 +
		d.Hour*60*60 +
		d.Day*60*60*24 +
		d.Month*60*60*24*30 +
		d.Year*60*60*24*365
}

func main() {
	durationString := "1h2m30s"

	duration := Duration{}

	err := duration.FromString(durationString)
	if err != nil {
		fmt.Println("error gan")
		return
	}

	fmt.Println(duration.GetSeconds())

	fmt.Println("Hello world gannnnn")
}
