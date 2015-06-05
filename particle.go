package berlin

import (
	"math"
)

//A particle has a K,V,u,Ms
type particle struct {
	mz     float64 // magnetisation along the z-direction
	u_anis float64 // Uniaxial anisotropy axis (angle)
	ku1    float64 // Uniaxial anisotropy constant in J/m**3
	r      float64 // radius
	msat   float64 // Saturation magnetisation in A/m

	//variables related to the energy landscape
	min1   float64 // position of first minimum
	m1     float64 // percentage of partices in first minimum
	E1     float64 // energy of first minimum
	min2   float64 // position of second minimum
	m2     float64 // percentage of particles in second minimum
	E2     float64 // energy of second minimum
	Ebar1  float64 // energy barrier to jump from E1 to E2
	Ebar2  float64 // energy barrier to jump from E2 to E1
	onemin bool    // boolean that states if there is one or 2 minima in the energy landscape

}

func (p *particle) V() float64 {
	return 4 / 3 * math.Pi * p.r * p.r * p.r
}

// returns the energy due to the anisotropy of the particle as function of theta
func (p *particle) E_anis(theta float64) float64 {
	return -p.ku1 * p.V() * (math.Sin(theta)*math.Sin(p.u_anis) + math.Cos(theta)*math.Cos(p.u_anis)) * (math.Sin(theta)*math.Sin(p.u_anis) + math.Cos(theta)*math.Cos(p.u_anis))
}

// returns the energy due to the external field as function of theta
func (p *particle) E_ext(theta float64) float64 {
	return -p.msat * p.V() * B_ext(T) * math.Cos(theta)
}

// returns the entropy times temperature as function of theta
func (p *particle) TS(theta float64) float64 {
	if (T==0.){return 0.}
	if (theta==0.){return 10000.}
	return Temp * kb * math.Log(1/2.*math.Sin(theta)*math.Sin(theta))
}

// returns the total free energy as function of theta
func (p *particle) F(theta float64) float64 {
	return p.E_anis(theta) + p.E_ext(theta) - p.TS(theta)
}

// looks for the position and energies of the minima in the free energy, angle accuracy is 0.001 rad
func (p *particle) Update_minima() {
	//find first minimum
	theta := 0.
	dt := 0.001
	ref := p.F(theta)
	theta += dt

	for p.F(theta) < ref {
		ref = p.F(theta)

		theta += dt
	}
	p.min1 = theta
	p.E1 = ref

	//find second minimum
	theta = math.Pi
	dt = 0.001
	ref = p.F(theta)
	theta -= dt

	for p.F(theta) < ref {
		ref = p.F(theta)
		theta += dt
	}
	p.min2 = theta
	p.E2 = ref
}

// looks for the position with maximum energy between the two minima (returns 0 if min1=min2)
func (p *particle) Update_maximum() {
	if math.Abs(p.min1-p.min2) < 0.001 {
		p.Ebar1 = 0.
		p.Ebar2 = 0.
		p.onemin = true
		return
	}
	theta := p.min1
	ref := p.E1
	dt := 0.001
	for p.F(theta) > ref {
		ref = p.F(theta)
		theta += dt
	}
	p.Ebar1 = ref - p.E1
	p.Ebar2 = ref - p.E2
	p.onemin = false
	return
}

// returns the z-component of the particle magnetisation
func (p *particle) M() float64 {
	return p.mz
}

//performs one timestep with stepsize Dt, using euler forward method
func (p *particle) step() {
	p.Update_minima()
	p.Update_maximum()

	if p.onemin {
		if p.m1 < math.Pi/2. {
			p.m1 = 1.
			p.m2 = 0.
			p.mz = math.Cos(p.min1)
		} else {
			p.m1 = 0.
			p.m2 = 1.
			p.mz = math.Cos(p.min2)
		}
		return
	}

	//update m1 and m2
	if p.onemin == false {
		onetotwo := tau0 * math.Exp(p.Ebar1/kb/Temp)
		twotoone := tau0 * math.Exp(p.Ebar2/kb/Temp)
		p.m1 += Dt * (p.m2*twotoone - p.m1*onetotwo)
		p.m2 += Dt * (p.m1*onetotwo - p.m2*twotoone)

	}
	//norm in between
	if p.m1+p.m2 != 1. {
		p.m1 = p.m1 / (p.m1 + p.m2)
		p.m2 = p.m2 / (p.m1 + p.m2)
	}

	//update M
	p.mz = p.m1*math.Cos(p.min1) + p.m2*math.Cos(p.min2)
}

// makes a new particle, given its anisotropy angle/constant,radius and msat
func NewParticle(radius float64, Msat float64, U_anis float64, Ku1 float64) *particle {
	return &particle{r: radius, msat: Msat, ku1: Ku1, u_anis: U_anis}
}

// Add a particle to the Particles list
func AddParticle(p *particle) {
	Particles = append(Particles, p)
}
