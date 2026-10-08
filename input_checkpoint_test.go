package ecmascript

import "testing"

func TestInputContextCheckpointBoundaries(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		processed uint64
		want      bool
	}{
		{0, true},
		{1, false},
		{255, false},
		{256, true},
		{257, false},
		{511, false},
		{512, true},
		{513, false},
	} {
		if got := inputContextCheckpoint(test.processed); got != test.want {
			t.Errorf("inputContextCheckpoint(%d) = %t; want %t", test.processed, got, test.want)
		}
	}
}
