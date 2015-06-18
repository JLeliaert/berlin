package berlin

import (
	"math"
)

//Runs the simulation for a certain time
func Run(time float64) {
	write()
	for j := T; T < j+time; {
		maxerr = 1.e-9
		step(Particles)
		write()
		T += Dt
		if Adaptivestep {
			Dt = 0.95 * Dt * math.Pow(1.e-7/maxerr, (1./2.))
		}
	}

}

//Perform a timestep using euler forward method
func step(Particles []*particle) {
	for _, p := range Particles {
		p.step()
	}
}
