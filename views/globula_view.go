package views

import (
	"encoding/json"
	"math/rand"
	"polymers/datatypes"
	dt "polymers/datatypes"
	"sort"
	"strconv"
)

type GlobulaView struct {
	name              string
	polymers          []*PolymerView
	xClusters         *ClusterView
	yClusters         *ClusterView
	zClusters         *ClusterView
	commonClusterDone bool
}

func NewGlobulaView(name string, polymers []*datatypes.Polymer) *GlobulaView {
	newGlobulaView := new(GlobulaView)
	newGlobulaView.name = name
	newGlobulaView.polymers = make([]*PolymerView, len(polymers))
	for i := 0; i < len(polymers); i++ {
		newGlobulaView.polymers[i] = NewPolymerView(polymers[i])
	}
	monomerNumber := int64(1)
	for _, polymer := range newGlobulaView.polymers {
		ForEachMonomer(polymer, func(mon *dt.Monomer) {
			mon.Number = monomerNumber
			monomerNumber++
		})
	}
	newGlobulaView.commonClusterDone = false
	return newGlobulaView
}

func (globula *GlobulaView) Name() string {
	return globula.name
}

func (globula *GlobulaView) Len() int {
	return len(globula.polymers)
}

func (globula *GlobulaView) Reset() {
	for _, pol := range globula.polymers {
		ForEachMonomer(pol, func(mon *datatypes.Monomer) { mon.MonomerType = datatypes.MONOMER_TYPE_USUAL })
	}
}

// func (globula *GlobulaView) ToJson() ([]byte, error) {
// 	pols_countMap := make(map[string]interface{})
// 	pols_countMap["value"] = len(globula.polymers)
// 	pols_countMap["name"] = "Количество полимеров"

// 	pols_map := make(map[string]interface{})
// 	pols_map["value"] = make([]string, len(globula.polymers))

// 	var jsonDict map[string]interface{}
// 	jsonDict["polymers_count"] = pols_countMap
// 	return nil, nil
// }

func ForEachPolymer(globula *GlobulaView, pred func(*PolymerView)) {
	for _, pol := range globula.polymers {
		pred(pol)
	}
}

func (globula *GlobulaView) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name     string
		Polymers []*PolymerView
	}{
		Name:     globula.name,
		Polymers: globula.polymers,
	})
}

func (globula *GlobulaView) XClusters(avg float64) *ClusterView {
	if globula.xClusters == nil {
		globula.xClusters = NewClusterView(globula, avg, dt.X_AXIS)
	}
	return globula.xClusters
}

func (globula *GlobulaView) YClusters(avg float64) *ClusterView {
	if globula.yClusters == nil {
		globula.yClusters = NewClusterView(globula, avg, dt.Y_AXIS)
	}
	return globula.yClusters
}

func (globula *GlobulaView) ZClusters(avg float64) *ClusterView {
	if globula.zClusters == nil {
		globula.zClusters = NewClusterView(globula, avg, dt.Z_AXIS)
	}
	return globula.zClusters
}

func (globula *GlobulaView) CommonClusters() (*ClusterView, *ClusterView, *ClusterView) {
	if !globula.commonClusterDone {
		clusters_X := globula.XClusters(0.0)
		clusters_Y := globula.YClusters(0.0)
		clusters_Z := globula.ZClusters(0.0)

		ForEachCluster(clusters_X, func(cluster_X *Cluster) {
			ForEachCluster(clusters_Y, func(cluster_Y *Cluster) {
				common_monomers := IntersectClusters_Soft(cluster_X, cluster_Y)
				for _, common_monomer := range common_monomers {
					ForEachCluster(clusters_X, func(cluster *Cluster) {
						cluster.RemoveMonomer(common_monomer)
					})
					ForEachCluster(clusters_Y, func(cluster *Cluster) {
						cluster.RemoveMonomer(common_monomer)
					})
				}
			})
		})

		ForEachCluster(clusters_Y, func(cluster_Y *Cluster) {
			ForEachCluster(clusters_Z, func(cluster_Z *Cluster) {
				common_monomers := IntersectClusters_Soft(cluster_Y, cluster_Z)
				for _, common_monomer := range common_monomers {
					ForEachCluster(clusters_Y, func(cluster *Cluster) {
						cluster.RemoveMonomer(common_monomer)
					})
					ForEachCluster(clusters_Z, func(cluster *Cluster) {
						cluster.RemoveMonomer(common_monomer)
					})
				}
			})
		})

		ForEachCluster(clusters_X, func(cluster_X *Cluster) {
			ForEachCluster(clusters_Z, func(cluster_Z *Cluster) {
				common_monomers := IntersectClusters_Soft(cluster_X, cluster_Z)
				for _, common_monomer := range common_monomers {
					ForEachCluster(clusters_X, func(cluster *Cluster) {
						cluster.RemoveMonomer(common_monomer)
					})
					ForEachCluster(clusters_Z, func(cluster *Cluster) {
						cluster.RemoveMonomer(common_monomer)
					})
				}
			})
		})
		// In case we have empty cluster, we need to remove them
		clusters_X.Trunkate()
		clusters_Y.Trunkate()
		clusters_Z.Trunkate()

		// Everything is done. We're ready to fully-connect them
		for _, cluster := range []*ClusterView{clusters_X, clusters_Y, clusters_Z} {
			ForEachCluster(cluster, func(c *Cluster) {
				c.MakeFullyConnected()
			})
		}
		// Now let's colorize them
		clusters_X.Colorize(false)
		clusters_Y.Colorize(false)
		clusters_Z.Colorize(false)

		globula.commonClusterDone = true
	}
	return globula.xClusters, globula.yClusters, globula.zClusters
}

// var turn int = 0
// var Bs_ = make([]*dt.Monomer, 0)

func (globula *GlobulaView) DoAging1(groupsCount int) {
	/*
		The aging process consists of 4 stages:
		1. Break connections and randomly assign the new ends as B and C so that we have 26% of B and 26% of C
		2. A part of Bs (3% of all groups to take) turn into C so that 29% of all groups are C
		3. Randomly turn other bins into C (6% in the first time and 6% in the second time) so that 41% of all groups are C
		4. Randomly create connections between bins untouched so that 59% of groups are B

		Example:
		Imagine we have 1000 bins at all and 100 aged groups required. The process is going to be the following:
		1. Break randomly (26/2)% of 100 = 26 connections and create 26 Bs and 26 Cs. 48 bins are left untouched
		2. Turn 3% of 100 = 3 Bs into C. 23 Bs left, 29 Cs left. 48 bins are left untouched
		3. Turn 6% of 100 = 6 bins randomly taken in globula into C twice. For now, we don't take Bs into account. 23 Bs left, 41 Cs left. 36 bins are left untouched
		4. Turn 59%-23%=36% of 100 = 36 bins untouched into Hs and create connections. 2 turns happen for 1 connection so that the number of connections is 36/2 = 18
	*/

	// =================DEBUG=================
	// // 4. Create connections
	// if turn == 3 {
	// 	globula.createBConnections(int(float64(groupsCount)*0.59), &Bs_)
	// 	turn++
	// }

	// // 3. Turn random bins into Cs (excluding Bs)
	// if turn == 2 {
	// 	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// 	// Do it twice
	// 	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// 	turn++
	// }

	// // 2. Turn some Bs into C
	// if turn == 1 {
	// 	turnBIntoC(&Bs_, int(float64(groupsCount)*0.03))
	// 	turn++
	// }

	// // 1. Break connections
	// if turn == 0 {
	// 	Bs_ = globula.breakConnections(int(float64(groupsCount) * 0.26))
	// 	turn++
	// }
	// =================DEBUG=================

	// 1. Break connections
	Cs_ := globula.breakConnections(int(float64(groupsCount)*0.26), dt.MONOMER_TYPE_NWISE)

	// 2. Turn some Bs into C
	turnIntoAnotherGroup(&Cs_, int(float64(groupsCount)*0.03), dt.MONOMER_TYPE_OWISE)

	// 3. Turn random bins into Cs (excluding Bs)
	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// Do it twice
	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))

	// 4. Create connections
	globula.createCrosslinks1(int(float64(groupsCount)*0.59), &Cs_)
}

func (globula *GlobulaView) DoAging2(groupsCount int) {
	// 1. Break connections
	Bs_ := globula.breakConnections(int(float64(groupsCount)*0.44), dt.MONOMER_TYPE_OWISE)

	// 2. Turn some Bs into C
	turnIntoAnotherGroup(&Bs_, int(float64(groupsCount)*0.15), dt.MONOMER_TYPE_NWISE)

	// 3. Turn random bins into Cs (excluding Bs)
	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// Do it twice
	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))

	// 4. Create connections
	globula.createCrosslinks2(int(float64(groupsCount) * 0.50))
}

func (globula *GlobulaView) breakConnections(groupsCount int, monomerTypeToGather dt.MonomerType) []*dt.Monomer {
	Bs := make([]*dt.Monomer, 0)
	for len(Bs) != groupsCount {
		chosenPoly := rand.Intn(globula.Len()) // Take a random poly
		poly := globula.polymers[chosenPoly]
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		nextMonomerNumber := chosenMonomerNumber + 1
		chosenMonomer := poly.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != dt.MONOMER_TYPE_USUAL {
			continue
		}
		nextMonomer := poly.polymer.GetMonomerByIdx(nextMonomerNumber)
		if nextMonomer.MonomerType != dt.MONOMER_TYPE_USUAL {
			continue
		}
		side := chosenMonomer.GetSideOfSibling(nextMonomer)
		datatypes.TierConnection(chosenMonomer, nextMonomer, side)
		if rand.Intn(2) == 0 {
			chosenMonomer.MonomerType = datatypes.MONOMER_TYPE_NWISE
			nextMonomer.MonomerType = datatypes.MONOMER_TYPE_OWISE
			if monomerTypeToGather == dt.MONOMER_TYPE_NWISE {
				Bs = append(Bs, chosenMonomer)
			} else {
				Bs = append(Bs, nextMonomer)
			}
		} else {
			chosenMonomer.MonomerType = datatypes.MONOMER_TYPE_OWISE
			nextMonomer.MonomerType = datatypes.MONOMER_TYPE_NWISE
			if monomerTypeToGather == dt.MONOMER_TYPE_OWISE {
				Bs = append(Bs, chosenMonomer)
			} else {
				Bs = append(Bs, nextMonomer)
			}
		}
	}

	return Bs
}

func turnIntoAnotherGroup(Bs *[]*dt.Monomer, groupsCount int, monomerTypeToTurn datatypes.MonomerType) {
	mapUsedBs := make(map[int]bool)
	for len(mapUsedBs) != groupsCount {
		mapUsedBs[rand.Intn(len(*Bs))] = true
	}

	arr := make([]int, 0)
	for key := range mapUsedBs {
		arr = append(arr, key)
	}
	sort.Ints(arr)

	for i := len(arr) - 1; i >= 0; i-- {
		(*Bs)[arr[i]].MonomerType = monomerTypeToTurn
		*Bs = append((*Bs)[:arr[i]], (*Bs)[arr[i]+1:]...)
	}
}

func (globula *GlobulaView) turnRandomBinsIntoC(groupsCount int) {
	i := 0
	for i < groupsCount {
		chosenPoly := rand.Intn(globula.Len()) // Take a random poly
		poly := globula.polymers[chosenPoly]
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType == dt.MONOMER_TYPE_USUAL {
			chosenMonomer.MonomerType = datatypes.MONOMER_TYPE_OWISE
			i++
		}
	}
}

func (globula *GlobulaView) createCrosslinks1(groupsCount int, Bs *[]*dt.Monomer) {
	for len(*Bs) != groupsCount {
		chosenPoly := rand.Intn(globula.Len()) // Take a random poly
		poly := globula.polymers[chosenPoly]
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != dt.MONOMER_TYPE_USUAL {
			continue
		}

		movementSides := dt.GetMovementSides()
		for i := 0; i < len(movementSides); i++ {
			chosenSide := movementSides[rand.Intn(len(movementSides))]
			nextMonomer, err := chosenMonomer.GetSibling(chosenSide)
			if err == nil && nextMonomer != nil && nextMonomer.MonomerType == dt.MONOMER_TYPE_USUAL {
				dt.MakeConnection(chosenMonomer, nextMonomer, dt.CONNECTION_TYPE_ONE)
				chosenMonomer.MonomerType = dt.MONOMER_TYPE_NWISE
				nextMonomer.MonomerType = dt.MONOMER_TYPE_NWISE
				*Bs = append(*Bs, chosenMonomer, nextMonomer)
				break
			}
		}
	}
}

func (globula *GlobulaView) createCrosslinks2(crosslinksCount int) {
	currentCount := 0
	timesRepeated := 0
	const maxTimesRepeated = 1000
	for currentCount != crosslinksCount && timesRepeated != maxTimesRepeated {
		chosenPoly := rand.Intn(globula.Len()) // Take a random poly
		poly := globula.polymers[chosenPoly]
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != dt.MONOMER_TYPE_USUAL {
			timesRepeated++
			continue
		}

		movementSides := dt.GetMovementSides()
		for i := 0; i < len(movementSides); i++ {
			chosenSide := movementSides[rand.Intn(len(movementSides))]
			nextMonomer, err := chosenMonomer.GetSibling(chosenSide)
			if err == nil &&
				nextMonomer != nil &&
				nextMonomer.MonomerType == dt.MONOMER_TYPE_USUAL &&
				(!dt.MonomersAreEqual(chosenMonomer.NextMonomer, nextMonomer) &&
					!dt.MonomersAreEqual(chosenMonomer.PrevMonomer, nextMonomer)) {
				dt.MakeConnection(chosenMonomer, nextMonomer, dt.CONNECTION_TYPE_ONE)
				chosenMonomer.MonomerType = dt.MONOMER_TYPE_HWISE
				nextMonomer.MonomerType = dt.MONOMER_TYPE_HWISE
				currentCount += 1
				timesRepeated = 0
				break
			}
		}
	}
	if timesRepeated == maxTimesRepeated {
		print("Didn't make all crosslinks! The number of crosslinks done: " + strconv.FormatInt(int64(currentCount), 10))
	}
}
