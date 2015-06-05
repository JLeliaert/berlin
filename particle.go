package berlin

import (
	"math"
)

//A particle has a K,V,u,Ms
type particle struct {
	mz     float64 // magnetisation along the z-direction
	u_anis float64 // Uniaxial anisotropy axis (angle)
	Ku1    float64 // Uniaxial anisotropy constant in J/m**3
	r      float64 // radius
	msat   float64 // Saturation magnetisation in A/m

	//variables related to the energy landscape
	min1 float64 // position of first minimum
	m1   float64 // percentage of partices in first minimum
	E1   float64 // energy of first minimum
	min2 float64 // position of second minimum
	m2   float64 // percentage of particles in second minimum
	E2   float64 // energy of second minimum
	Emax float64 // energy of maximum

}

func (p particle) V() float64 {
	return 4 / 3 * math.Pi * p.r * p.r * p.r
}

// returns the energy due to the anisotropy of the particle as function of theta
func (p particle) E_anis(theta float64) float64 {
	return -p.Ku1 * p.V() * math.Cos(theta) * math.Cos(p.u_anis) * math.Cos(theta) * math.Cos(p.u_anis)
}

// returns the energy due to the external field as function of theta
func (p particle) E_ext(theta float64) float64 {
	return -p.msat * p.V() * B_ext(T) * math.Cos(theta)
}

// returns the entropy times temperature as function of theta
func (p particle) TS(theta float64) float64 {
	return Temp * kb * math.Log(1/2.*math.Sin(theta)*math.Sin(theta))
}

// returns the total free energy as function of theta
func (p particle) F(theta float64) float64 {
	return p.E_anis(theta) + p.E_ext(theta) - p.TS(theta)
}

//TODO
func (p particle) Update_minima() {
}

//TODO
func (p particle) Update_maximum() {
}

//TODO
func (p particle) M() float64 {
	return p.mz
}
