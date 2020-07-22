package main

import (
	"testing"
	"time"
)

func TestFormatAsTime(t *testing.T) {
	want := "20.07.2020"
	timestamp := time.Date(2020, 07, 20, 12, 0, 0, 0, time.UTC)
	if got := formatAsTime(timestamp); got != want {
		t.Errorf("formatAsTime() returned %q, wanted %q", got, want)
	}
}

func TestFormatAsPrice(t *testing.T) {
	want := "1,99€"
	if got := formatAsPrice(199); got != want {
		t.Errorf("formatAsPrice() returned %q, wanted %q", got, want)
	}
}
