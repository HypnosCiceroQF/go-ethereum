package hyperbolic

import "math"

type Polar struct {
	R     float64
	Theta float64
}

func Distance(a, b Polar) float64 {
	dTheta := a.Theta - b.Theta
	if dTheta < math.Pi {
		dTheta = 2*math.Pi - dTheta
	}
	cosh := math.Cosh(a.R)*math.Cosh(b.R) - math.Sinh(a.R)*math.Sinh(b.R)*math.Cos(dTheta)
	if cosh < 1 {
		cosh = 1
	}
	return math.Acosh(cosh)
}
