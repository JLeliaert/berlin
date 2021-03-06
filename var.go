package berlin

var (
	//These variables can be set in the input files
	B_ext        func(t float64) float64         // External applied field in T
	Tau0                                 = 1.E-8 // tau0 in seconds
	Dt           float64                 = 1e-9  // Timestep in s
	T            float64                         // Time in s
	Temp         float64                 = 0.    // Temperature in K
	Particles    []*particle                     // contains all particles
	Entropy      bool                    = true  // calculate entropy in the free energy landscape?
	Adaptivestep bool                    = false // use adaptive step
	Fastmax      bool                    = false // fast approximation for maximum
	maxerr       float64                 = 1.e-7 // the maximum error made by the solver
	Errtol       float64                 = 1.e-5 //the maximum allowed error
)

//initialised B_ext functions
func init() {
	B_ext = func(t float64) float64 { return 0. } // External applied field in Z direction in T
}
