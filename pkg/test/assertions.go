package test

import (
	"regexp"
	"testing"
)

func AssertTimeFormat(t testing.TB, time string) {
	t.Helper()
	r := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)

	if !r.MatchString(time) {
		t.Fatalf("time format is not valid: %s", time)
	}
}
