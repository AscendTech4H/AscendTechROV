package controlDriver

import "testing"

func TestMotorCalcForward(t *testing.T) {
	l, r := motorCalcFwd(0, 50)
	if l != r || r != 217 {
		t.Fatalf("Expected full speed forward but got (%d, %d)", l, r)
	}
	l, r = motorCalcFwd(0, -50)
	if l != r || r != 37 {
		t.Fatalf("Expected full speed backward but got (%d, %d)", l, r)
	}
	l, r = motorCalcFwd(50, 0)
	if l >= r {
		t.Fatalf("Expected left turn but got (%d, %d)", l, r)
	}
	l, r = motorCalcFwd(-50, 0)
	if l <= r {
		t.Fatalf("Expected right turn but got (%d, %d)", l, r)
	}
	l, r = motorCalcFwd(25, 25)
	if l >= r {
		t.Fatalf("Expected slight left turn but got (%d, %d)", l, r)
	}
}

func TestRangeMap(t *testing.T) {
	val := rangeMap(-1, -2, 0, 1, 3)
	if val != 2 {
		t.Fatalf("Expected 2 but got %d", val)
	}
	val = rangeMap(-2, -2, 0, 1, 3)
	if val != 1 {
		t.Fatalf("Expected 1 but got %d", val)
	}
}
