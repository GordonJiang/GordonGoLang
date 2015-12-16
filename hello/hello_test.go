package main

import (
	"testing"
)

func TestStringSliceEqual(t *testing.T) {
	t.Log(StringSliceEqual([]string{""}, []string{"ab"}))
	t.Log("In Test")
}
