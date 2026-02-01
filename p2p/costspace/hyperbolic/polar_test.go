package hyperbolic

import (
	"math"
	"testing"
)

func almost(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestDistanceZero(t *testing.T) {
	a := Polar{R: 1.2, Theta: 0.7}
	if !almost(Distance(a, a), 0) {
		t.Fatalf("expected 0")
	}
}

func TestDistanceSymmetry(t *testing.T) {
	a := Polar{R: 2, Theta: 0.1}
	b := Polar{R: 3, Theta: 2.2}
	dab := Distance(a, b)
	dba := Distance(b, a)
	if math.Abs(dab-dba) > 1e-9 {
		t.Fatalf("not symmetric")
	}
}

func TestAngleWrap(t *testing.T) {
	a := Polar{R: 2, Theta: 0.01}
	b := Polar{R: 2, Theta: 2*math.Pi - 0.01}
	d := Distance(a, b)
	if d > 0.2 { // 很小的角差应该距离不大（这里只做粗断言）
		t.Fatalf("wrap seems wrong: %v", d)
	}
}
