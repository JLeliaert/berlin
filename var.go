package berlin

var (
	//These variables can be set in the input files
	B_ext       func(t float64) (float64, float64, float64)               // External applied field in T
	B_ext_space func(t, x, y, z float64) (float64, float64, float64)      // External applied field in T
	Dt          float64                                              = -1 // Timestep in s
	T           float64                                                   // Time in s
	Temp        float64                                              = -1 // Temperature in K
)

//initialised B_ext functions
func init() {
	B_ext = func(t float64) (float64, float64, float64) { return 0, 0, 0 }                // External applied field in T
	B_ext_space = func(t, x, y, z float64) (float64, float64, float64) { return 0, 0, 0 } // External applied field in T
}
