package main

import "testing"

func TestVersion(t *testing.T) {
	if version == "" {
		t.Fatal("version should not be empty")
	}
}

func TestPrintBanner(t *testing.T) {
	// smoke test, confirm it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("printBanner panicked: %v", r)
		}
	}()
	printBanner()
}
