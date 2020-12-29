package format

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AsPrice formats a given cent amount to Eurostring
func AsPrice(cents int) string {
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

// FromPrice takes a price string (e. g. 42,10) and return the cent value (e. g. 4210)
func FromPrice(price string) (int, error) {
	re := regexp.MustCompile(`^(\d)+(,(\d){2})?$`)
	if !re.MatchString(price) {
		return 0, errors.New("price has not got the right format")
	}

	splittedPrice := strings.Split(price, ",")
	euros, euroErr := strconv.ParseInt(splittedPrice[0], 10, 64)
	var cents int64
	var centsErr error
	if len(splittedPrice) > 1 {
		cents, centsErr = strconv.ParseInt(splittedPrice[1], 10, 64)
	}

	if euroErr != nil || centsErr != nil {
		return 0, errors.New("parsing price failed")
	}
	return int(100*euros + cents), nil
}

// AsTime formats a timestamp so it's human readable
func AsTime(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%02d.%02d.%d", day, month, year)
}
