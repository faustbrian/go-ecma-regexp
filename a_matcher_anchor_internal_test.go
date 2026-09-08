package ecmascript

import "testing"

func TestEndAnchorRequiresTheFinalInputPosition(t *testing.T) {
	executor := executor{
		input: &inputView{units: []uint16{'a'}},
	}

	if executor.atEnd(0, Flags{}) {
		t.Fatal("end anchor accepted a non-final input position")
	}
	if !executor.atEnd(1, Flags{}) {
		t.Fatal("end anchor rejected the final input position")
	}
}
