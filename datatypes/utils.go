package datatypes

import "polymers/base"

func distance_of_monomers(mon1, mon2 *Monomer) float64 {
	return base.EcludianDistanceF(mon1.coords, mon2.coords)
}
