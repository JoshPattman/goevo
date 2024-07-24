package goevo

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestActivationLoading(t *testing.T) {
	s := ""
	for _, a := range AllActivations {
		s += a.String() + "\n"
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(AllActivations); err != nil {
		t.Fatal(err)
	}
	var loaded []Activation
	if err := json.NewDecoder(buf).Decode(&loaded); err != nil {
		t.Fatal(err)
	}
	if len(loaded) != len(AllActivations) {
		t.Fatalf("unmatching lengths: %v and %v", len(loaded), len(AllActivations))
	}
	for i := range loaded {
		if loaded[i] != AllActivations[i] {
			t.Fatalf("unmatching activations: %v and %v", loaded[i], AllActivations[i])
		}
	}
}
