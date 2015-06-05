package berlin

import ()

//A particle has a K,V,u,Ms
type particle struct {
	mz     float64 // magnetisation along the z-direction
	m1     float64 // percentage of partices in first minimum
	m2     float64 // percentage of particles in second minimum
	u_anis Vector  // Uniaxial anisotropy axis
	Ku1    float64 // Uniaxial anisotropy constant in J/m**3
	r      float64 // radius
	msat   float64 // Saturation magnetisation in A/m
}

func (p particle) E_anis(theta float64) float64 {
	return 0.
}

func (p particle) E_ext(theta float64) float64 {
	return 0.
}

func (p particle) TS(theta float64) float64 {
	return 0.
}

func (p particle) F(theta float64) float64 {
	return p.E_anis(theta) + p.E_ext(theta) - p.TS(theta)
}
