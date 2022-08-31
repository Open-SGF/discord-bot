package util

import (
	"testing"
	"time"
)

func TestCalculateNextBlastDate(t *testing.T) {
	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	now, err := time.Parse(layout, "2022-08-30 19:09:17.814691 -0500 CDT")
	if err != nil {
		t.Fatalf("%v", err)
	}

	expected, err := time.Parse(layout, "2022-09-05 10:00:00 -0500 CDT")
	if err != nil {
		t.Fatalf("%v", err)
	}

	actual := CalculateNextBlastDate(now, time.Monday, 10*time.Hour)
	if actual.String() != expected.String() {
		t.Fatalf("expected %s got %s", expected, actual)
	}
}
