package berlin

import ()

//Runs the simulation for a certain time
func Run(time float64) {
	write()
	for j := T; T < j+time; {
		step(Particles)
		write()
		T += Dt
	}

}

//Perform a timestep using euler forward method
func step(Particles []*particle) {
	for _, p := range Particles {
		p.step()
	}
}
