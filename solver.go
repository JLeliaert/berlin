package berlin

import ()

//Runs the simulation for a certain time
func Run(time float64) {
	//TODO, remove DEBUG only!!!!
	for _, p := range Particles {
		p.m1=1
	}




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
