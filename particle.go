package berlin

import (
	//		"fmt"
	"math"
)

//A particle has a K,V,u,Ms
type particle struct {
	mz     float64 // magnetisation along the z-direction
	u_anis float64 // Uniaxial anisotropy axis (angle)
	ku1    float64 // Uniaxial anisotropy constant in J/m**3
	r      float64 // radius
	msat   float64 // Saturation magnetisation in A/m
	weight float64 // number of particles with these properties

	//variables related to the energy landscape
	min1   Coord   // position of first minimum
	m1     float64 // percentage of partices in first minimum
	E1     float64 // energy of first minimum
	min2   Coord   // position of second minimum
	m2     float64 // percentage of particles in second minimum
	E2     float64 // energy of second minimum
	Ebar1  float64 // energy barrier to jump from E1 to E2
	Ebar2  float64 // energy barrier to jump from E2 to E1
	onemin bool    // boolean that states if there is one or 2 minima in the energy landscape

	previousm1 float64
	previousm2 float64

}

func (p *particle) V() float64 {
	//fmt.Println(p.r)
	return 4. / 3. * math.Pi * p.r * p.r * p.r
}

// returns the energy due to the anisotropy of the particle as function of theta
func (p *particle) E_anis(theta, psi float64) float64 {
	angles := (math.Sin(theta)*math.Sin(p.u_anis)*math.Cos(psi) + math.Cos(theta)*math.Cos(p.u_anis))
	return -p.ku1 * p.V() * angles * angles
}

// returns the energy due to the external field as function of theta
func (p *particle) E_ext(theta float64) float64 {
	return -p.msat * p.V() * B_ext(T) * math.Cos(theta)
}

// returns the entropy times temperature as function of theta
func (p *particle) TS(theta float64) float64 {
	if Entropy == false {
		return 0.
	}
	if Temp == 0. {
		return 0.
	}
	if theta < 0.0001 || theta > math.Pi-0.0001 {
		return -1.e-10
	}
	return Temp * kb * math.Log(1./2.*math.Sin(theta)*math.Sin(theta))
}

// returns the total free energy as function of theta
func (p *particle) F(theta, psi float64) float64 {
	return p.E_anis(theta, psi) + p.E_ext(theta) - p.TS(theta)
}

// returns the derivative of the energy due to the anisotropy of the particle as function to theta
func (p *particle) dE_anis(theta, psi float64) float64 {
	angles1 := math.Cos(theta)*math.Cos(psi)*math.Sin(p.u_anis) - math.Sin(theta)*math.Cos(p.u_anis)
	angles2 := math.Cos(theta)*math.Cos(p.u_anis) + math.Cos(psi)*math.Sin(theta)*math.Sin(p.u_anis)
	return -2. * p.ku1 * p.V() * angles1 * angles2
}

// returns the derivative of energy due to the external field as function to theta
func (p *particle) dE_ext(theta float64) float64 {
	return p.msat * p.V() * B_ext(T) * math.Sin(theta)
}

// returns the derivative of the entropy times temperature as function to theta
func (p *particle) dTS(theta float64) float64 {
	if Entropy == false {
		return 0.
	}
	if Temp == 0. {
		return 0.
	}
	if theta < 0.0001 {
		return 1.e-10
	}
	if theta > math.Pi-0.0001 {
		return -1.e-10
	}

	return Temp * kb / math.Tan(theta)
}

// returns the derivative of the total free energy to theta
func (p *particle) dFdtheta(theta, psi float64) float64 {
	return p.dE_anis(theta, psi) + p.dE_ext(theta) - p.dTS(theta)
}

// looks for the position and energies of the minima in the free energy, angle accuracy is 0.0001 rad
func (p *particle) Update_minima() {

	//for i :=0.;i<math.Pi;i+=0.01{
	//	fmt.Println(i, p.F(i,0.))
	//}

	//find first minimum
	theta := 0.
	psi := 0.
	dt := 0.1
	ref := p.F(theta, psi)
	theta += dt

	for dt > 0.00001 {
		for p.F(theta, psi) <= ref {
			//	fmt.Println(dt,theta,ref-p.F(theta,psi))
			ref = p.F(theta, psi)
			theta += dt
		}
		theta -= 2 * dt
		if theta < 0. {
			theta = 0.
		}
		ref = p.F(theta, psi)
		dt /= 2.
	}

	p.min1 = Coord{theta, psi}
	p.E1 = ref

	//find second minimum
	theta = math.Pi
	psi = math.Pi
	dt = 0.1
	ref = p.F(theta, psi)

	for dt > 0.00001 {
		for p.F(theta, psi) <= ref {
			ref = p.F(theta, psi)
			theta -= dt
		}
		theta += 2 * dt
		if theta > math.Pi {
			theta = math.Pi
		}
		ref = p.F(theta, psi)
		dt /= 2.
	}

	p.min2 = Coord{theta, psi}
	p.E2 = ref

	//fmt.Println("min1", p.min1[0])
	//fmt.Println(T, p.min1[0], p.min2[0])
	//fmt.Println("min1psi", p.min1[1])
	//fmt.Println("m1  ", p.m1)
	//fmt.Println("E1  ", p.E1)
	//fmt.Println("min2", p.min2[0])
	//fmt.Println("min2psi", p.min2[1])
	//fmt.Println("m2  ", p.m2)
	//fmt.Println("E2  ", p.E2)
	//fmt.Println("m1+m2", p.m1+p.m2)
	//fmt.Println("onemine", p.onemin)
	//fmt.Println("ebar1", p.Ebar1)
	//fmt.Println("ebar2", p.Ebar2)
	//fmt.Println()
}

// looks for the position with maximum energy between the two minima (returns 0 if min1=min2)
func (p *particle) Update_maximum() {
	if p.min1.Dist(p.min2) < 0.00001 {
		p.Ebar1 = 0.
		p.Ebar2 = 0.
		p.onemin = true
		return
	}

	//fmt.Println(p.min1[0],p.min2[0])
	if math.Abs(p.min1[0]-p.min2[0]) < 0.0021 {
		ref := p.F(p.min1[0], math.Pi/2.)
		p.Ebar1 = ref - p.E1
		p.Ebar2 = ref - p.E2
		p.onemin = false
		return
	}

	theta := p.min1[0]
	psi := math.Pi / 2.
	ref := p.dFdtheta(theta, psi)

	//find maximum
	//fmt.Println("ref", ref)
	dt := 0.
	//if p.min1[0] < p.min2[0] {
	//	dt = 0.1
	//} else {
	//	dt = -0.1
	//}
	dt = (p.min2[0] - p.min1[0]) / 3.

	//until df/dtheta changes sign or not found
	for math.Abs(dt) > 0.00001 {
		for p.dFdtheta(theta, psi)*ref >= 0. {
			//	fmt.Println(p.dFdtheta(theta, psi))
			ref = p.dFdtheta(theta, psi)
			theta += dt
			if theta > math.Pi {
				theta = p.min1[0] + p.min2[0]/2.
				break
			}
		}
		theta -= 2 * dt
		if theta < 0. {
			theta = 0.
		}
		if theta > math.Pi {
			theta = math.Pi
		}
		ref = p.dFdtheta(theta, psi)
		dt /= 2.
	}

	ref = p.F(theta, psi)

	p.Ebar1 = ref - p.E1
	p.Ebar2 = ref - p.E2
	//fmt.Println(p.Ebar1)
	//fmt.Println(p.Ebar2)
	//fmt.Println()
	p.onemin = false
	return
}

// returns the z-component of the particle magnetisation
func (p *particle) M() float64 {
	return p.mz
}

//puts all the particles in their ground state
func Relax() {
	for i := range Particles {
		Particles[i].relax()
	}

}

//puts the magnetisation of the particle in its ground state
func (p *particle) relax() {

	p.Update_minima()
	p.Update_maximum()

	// if there is only one minimum, this is all there is to it
	if p.onemin {
		if p.min1[0] < math.Pi/2. {
			p.m1 = 1.
			p.m2 = 0.
			p.mz = math.Cos(p.min1[0])
		} else {
			p.m1 = 0.
			p.m2 = 1.
			p.mz = math.Cos(p.min2[0])
		}
		return
	}

	p.m1 = math.Exp(p.E1-p.E2) / (1 + math.Exp(p.E1-p.E2))

	p.m2 = 1 - p.m1

	//update mz
	p.mz = p.m1*math.Cos(p.min1[0]) + p.m2*math.Cos(p.min2[0])
}

//performs one timestep with stepsize Dt, using euler forward method
func (p *particle) step() {

	p.previousm1=p.m1
	p.previousm2=p.m2
	//step 1 of heun scheme
	p.Update_minima()
	p.Update_maximum()

	// if there is only one minimum, 1 step is all there is to it
	if p.onemin {
		if p.min1[0] < math.Pi/2. {
			p.m1 = 1.
			p.m2 = 0.
			p.mz = math.Cos(p.min1[0])
		} else {
			p.m1 = 0.
			p.m2 = 1.
			p.mz = math.Cos(p.min2[0])
		}
		return
	}
	k1 := 0.
	k2 := 0.

	predicted := 0.
	corrected := 0.
	//update m1 and m2
	if p.onemin == false {
		if Temp != 0. {
			onetotwo := p.m1 / Tau0 / math.Exp(p.Ebar1/kb/Temp)
			twotoone := p.m2 / Tau0 / math.Exp(p.Ebar2/kb/Temp)
			k1 = onetotwo
			k2 = twotoone

			p.m1 += Dt * (twotoone - onetotwo)
			predicted = Dt*(twotoone - onetotwo)
			p.m2 -= Dt * (twotoone - onetotwo)

		}
	}

	//step 2 of heun scheme
	p.Update_minima()
	p.Update_maximum()
	if Temp != 0. {
		onetotwo := p.m1 / Tau0 / math.Exp(p.Ebar1/kb/Temp)
		twotoone := p.m2 / Tau0 / math.Exp(p.Ebar2/kb/Temp)

		p.m1 += 0.5 * Dt * (twotoone - onetotwo - k2 + k1)
		p.m2 -= 0.5 * Dt * (twotoone - onetotwo - k2 + k1)
		corrected=Dt * (twotoone - onetotwo)

	}
	if Adaptivestep && math.Abs(corrected-predicted) > maxerr {
		maxerr = math.Abs(corrected - predicted)
	}
	p.m1 = p.m1 / (p.m1 + p.m2)
	p.m2 = p.m2 / (p.m1 + p.m2)

	//update mz
	p.mz = p.m1*math.Cos(p.min1[0]) + p.m2*math.Cos(p.min2[0])
}

// makes a new particle, given its anisotropy angle/constant,radius and msat
func NewParticle(radius float64, Msat float64, U_anis float64, Ku1 float64) *particle {
	return &particle{r: radius, msat: Msat, ku1: Ku1, u_anis: U_anis, m1: 1., weight: 1.}
}

// Add a particle to the Particles list
func AddParticle(p *particle) {
	Particles = append(Particles, p)
}

// Copies a given particle, can be used to make size or anisotropy axis distributions
func CopyParticle(p *particle) *particle {
	return &particle{r: p.r, msat: p.msat, ku1: p.ku1, u_anis: p.u_anis, m1: 1., weight: 1.}
}

//helper function that returns the weight for the axis distribution
func axisweight(N int) float64 {

	totalweight := 0.
	for i := 0; i < N; i += 1 {
		totalweight += math.Sin(math.Pi / 2. * float64(i) / float64(N-1))
	}
	return totalweight
}

// gives the ensemble random anisotropy axes (pi/2N is discritization in radians, N is number of u_anis directions)
func Random_anis_axis(N int) {
	var Newparticles []*particle
	totalweight := axisweight(N)
	for j := range Particles {
		for i := 0; i < N; i += 1 {
			weight := math.Sin(math.Pi/2.*float64(i)/float64(N-1)) / totalweight
			if weight*Particles[j].weight > 0.0001 {
				newparticle := CopyParticle(Particles[j])
				newparticle.u_anis = math.Pi / 2. * float64(i) / float64(N-1)
				newparticle.weight = Particles[j].weight * weight
				Newparticles = append(Newparticles, newparticle)
			}
		}
	}
	Particles = Newparticles
}

//helper function that returns the weight for the lognormal distribution
func lognormal(D, avg, stdev float64) float64 {
	prefactor := 1. / (math.Sqrt(2.*math.Pi) * stdev * D)
	exponent := -math.Log(D/avg) * math.Log(D/avg) / 2. / stdev / stdev
	return prefactor * math.Exp(exponent)
}

// gives the ensemble lognormal DIAMETER(!) size distribution, based on top cutoff (in nm), discretization (in nm), average (in nm) and stdev (in nm)
func Lognormal_sizes(top, discr, avg, stdev float64) {
	var Newparticles []*particle
	for i := discr; i <= top; i += discr {
		for j := range Particles {
			dist := lognormal(i, avg, stdev)
			if dist*Particles[j].weight > 0.0001 {
				newparticle := CopyParticle(Particles[j])
				newparticle.r = i / 2. * 1e-9
				newparticle.weight = Particles[j].weight * dist
				Newparticles = append(Newparticles, newparticle)
			}
		}
	}
	Particles = Newparticles
}
