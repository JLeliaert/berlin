package berlin

var (
	//These variables can be set in the input files
	B_ext     func(t float64) float64      // External applied field in T
	Dt        float64                 = -1 // Timestep in s
	T         float64                      // Time in s
	Temp      float64                 = -1 // Temperature in K
	Particles []*particle                  // contains all particles
)

//initialised B_ext functions
func init() {
	B_ext = func(t float64) float64 { return 0. } // External applied field in Z direction in T
}
