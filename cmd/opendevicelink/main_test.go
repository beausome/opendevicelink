package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDemoCompletesAndDisclosesSimulation(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	if err := run([]string{"demo"}, &out, &errOut); err != nil {
		t.Fatalf("run demo: %v", err)
	}
	for _, want := range []string{"Result: PASS", "simulation only", "lease released"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("output %q does not contain %q", out.String(), want)
		}
	}
}

func TestUnknownCommandFails(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	if err := run([]string{"missing"}, &out, &errOut); err == nil {
		t.Fatal("unknown command unexpectedly succeeded")
	}
}
