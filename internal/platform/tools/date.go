package tools

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidDay   = fmt.Errorf("Invalid day")
	ErrInvalidMonth = fmt.Errorf("Invalid month")
	ErrInvalidYear  = fmt.Errorf("Invalid year")
	ErrInvalidDate  = fmt.Errorf("Invalid date")
)

// ParseDate parses a date in the format dd/mm/yyyy and returns an error if the date is invalid
func ParseDate(expiration string) (err error) {
	date := strings.Split(expiration, "/")

	// validate date format (dd/mm/yyyy)
	if len(date) != 3 {
		err = ErrInvalidDate
		return
	}

	// validate day
	day, _ := strconv.Atoi(date[0])
	if day < 1 || day > 31 {
		err = ErrInvalidDay

		return
	}

	// validate month
	month, _ := strconv.Atoi(date[1])
	if month < 1 || month > 12 {
		err = ErrInvalidMonth

		return
	}

	// validate year
	year, err := strconv.Atoi(date[2])
	if year < 1999 || year > 2099 {
		err = ErrInvalidYear

		return
	}

	return
}
