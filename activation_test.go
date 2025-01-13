package goevo

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestActivationLoading(t *testing.T) {
	s := ""
	for _, a := range AllSingleActivations {
		s += a.String() + "\n"
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(AllSingleActivations); err != nil {
		t.Fatal(err)
	}
	var loaded []Activation
	if err := json.NewDecoder(buf).Decode(&loaded); err != nil {
		t.Fatal(err)
	}
	if len(loaded) != len(AllSingleActivations) {
		t.Fatalf("unmatching lengths: %v and %v", len(loaded), len(AllSingleActivations))
	}
	for i := range loaded {
		if loaded[i] != AllSingleActivations[i] {
			t.Fatalf("unmatching activations: %v and %v", loaded[i], AllSingleActivations[i])
		}
	}
}
