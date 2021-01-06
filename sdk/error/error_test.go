package error

import "testing"

func TestErrorEqual(t *testing.T) {
	t0 := NewError("abc ")
	t1 := NewError("abc")
	t2 := NewError("abc")
	if t0 == t1 {
		t.Fatal()
	}
	if t1 != t2 {
		t.Fatal()
	}
}
