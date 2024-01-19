package testdata

import "testing"

func TestS(t *testing.T) {
	var s = new(S)
	s2 := s.DeepCopy()
	t.Log(s2)
}
