package berlin

import (
	"math"
)

//Runs the simulation for a certain time
func Run(time float64) {

	write()
	for j := T; T < j+time; {
		maxerr = 1.e-6
		step(Particles)
		if Adaptivestep {
			if maxerr > Errtol {
				undobadstep(Particles)
			}
			Dt = 0.9 * Dt * math.Pow(Errtol/maxerr, (1./2.))
		}
		T += Dt
		write()
	}

}

//Perform a timestep using euler forward method
func step(Particles []*particle) {
	for _, p := range Particles {
		p.step()
	}
}

//Undoes a bad step
func undobadstep(Particles []*particle) {
	for _, p := range Particles {
		p.m1 = p.previousm1
		p.m2 = p.previousm2
	}
	T -= Dt
}
