package berlin

import (
)

//A particle has a K,V,u,Ms
type particle struct {
	mz                  float64 // magnetisation along the z-direction
	m1		    float64 // percentage of partices in first minimum
	m2		    float64 // percentage of particles in second minimum
	u_anis              vector  // Uniaxial anisotropy axis
	r                   float64 // radius
	msat                float64 // Saturation magnetisation in A/m
}
