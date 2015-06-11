package berlin

import "math"

// 2-component coordinate
type Coord [2]float64

func (v Coord) Dist(c Coord) float64 {
	return math.Sqrt((v[0]-c[0])*(v[0]-c[0]) + (v[1]-c[1])*(v[1]-c[1]))
}
