package server

import (
	"math"

	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

func distance(p, q *ups.DeliveryLocation) float64 {
	return math.Hypot(float64(p.GetX())-float64(q.GetX()), float64(p.GetY())-float64(q.GetY()))
}

func pathPlanning(x, y int32, dLocs []*ups.DeliveryLocation) (r []*ups.DeliveryLocation) {
	curr := &ups.DeliveryLocation{X: &x, Y: &y}
	for len(dLocs) > 0 {
		nearest := 0
		dist := distance(curr, dLocs[0])
		for i := 1; i < len(dLocs); i++ {
			dist1 := distance(curr, dLocs[i])
			if dist1 < dist {
				nearest, dist = i, dist1
			}
		}
		// move to nearest
		curr = dLocs[nearest]
		r = append(r, curr)
		// remove nearest from dLocs
		dLocs[nearest] = dLocs[len(dLocs)-1]
		dLocs = dLocs[:len(dLocs)-1]
	}
	return
}
