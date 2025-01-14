package goevo

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestActivationLoading(t *testing.T) {
	s := ""
	for _, a := range AllVectorActivations {
		s += a.String() + "\n"
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(AllVectorActivations); err != nil {
		t.Fatal(err)
	}
	var loaded []Activation
	if err := json.NewDecoder(buf).Decode(&loaded); err != nil {
		t.Fatal(err)
	}
	if len(loaded) != len(AllVectorActivations) {
		t.Fatalf("unmatching lengths: %v and %v", len(loaded), len(AllVectorActivations))
	}
	for i := range loaded {
		if loaded[i] != AllVectorActivations[i] {
			t.Fatalf("unmatching activations: %v and %v", loaded[i], AllVectorActivations[i])
		}
	}
}
