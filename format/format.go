package format

import (
	"fmt"
	"strconv"
	"time"
)

// Formats a given cent amount to Eurostring
func FormatAsPrice(cents int) string {
	result := ""
	posCents := cents
	if cents < 0 {
		posCents = -cents
	}
	if posCents%100 >= 10 {
		result = strconv.FormatInt(int64(posCents/100), 10) + "," + strconv.FormatInt(int64(posCents%100), 10) + "€"
	} else {
		result = strconv.FormatInt(int64(posCents/100), 10) + ",0" + strconv.FormatInt(int64(posCents%100), 10) + "€"
	}
	if cents < 0 {
		return "-" + result
	}
	return result
}

// Formats a timestamp so it's human readable
func FormatAsTime(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%02d.%02d.%d", day, month, year)
}
