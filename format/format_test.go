package format

import (
	"testing"
	"time"
)

func TestFormatAsTime(t *testing.T) {
	want := "20.07.2020"
	timestamp := time.Date(2020, 07, 20, 12, 0, 0, 0, time.UTC)
	if got := AsTime(timestamp); got != want {
		t.Errorf("format.AsTime() returned %q, wanted %q", got, want)
	}
}

func TestFormatAsPrice(t *testing.T) {
	want := "1,99€"
	if got := AsPrice(199); got != want {
		t.Errorf("format.AsPrice() returned %q, wanted %q", got, want)
	}
}

func TestFromPrice(t *testing.T) {
	want := 4210
	if got, _ := FromPrice("42,10"); got != want {
		t.Errorf("format.FromPrice() returned %d, wanted %d", got, want)
	}
	if _, err := FromPrice("1,5"); err == nil {
		t.Errorf("format.FromPrice() returned no error when %q was given as argument.", "1,5")
	}
}
