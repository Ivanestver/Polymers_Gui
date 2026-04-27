package views

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"math/rand"
	"polymers/base"
	dt "polymers/datatypes"
	"polymers/outputformat"
	"sort"
	"strconv"
	"strings"
)

type GlobulaProperty int

const (
	GlobulaAged GlobulaProperty = iota
	GlobulaWaterized
	GlobulaGlobulaType
	GlobulaThreadType
	GlobulaSurfaceType
)
const crosslinksCount = 0.5

type GlobulaView struct {
	polymers          []*PolymerView
	xClusters         *ClusterView
	yClusters         *ClusterView
	zClusters         *ClusterView
	commonClusterDone bool
	literalsTable     map[dt.MonomerType]string
	globulaProperties map[GlobulaProperty]bool
}

func NewGlobulaView(polymers []dt.IPolymer, globulaType GlobulaProperty, literalsTable map[dt.MonomerType]string) *GlobulaView {
	newGlobulaView := new(GlobulaView)
	newGlobulaView.polymers = make([]*PolymerView, len(polymers))
	for i := 0; i < len(polymers); i++ {
		newGlobulaView.polymers[i] = NewPolymerView(polymers[i])
	}
	monomerNumber := int64(1)
	for _, polymer := range newGlobulaView.polymers {
		ForEachMonomer(polymer, func(mon *dt.Monomer) bool {
			mon.Number = monomerNumber
			monomerNumber++
			return true
		})
	}
	newGlobulaView.commonClusterDone = false
	newGlobulaView.literalsTable = make(map[dt.MonomerType]string)
	maps.Copy(newGlobulaView.literalsTable, literalsTable)
	newGlobulaView.globulaProperties = make(map[GlobulaProperty]bool)
	newGlobulaView.globulaProperties[globulaType] = true
	return newGlobulaView
}

func (globula *GlobulaView) Len() int {
	return len(globula.polymers)
}

func (globula *GlobulaView) GetPolymerByIdx(idx int) *PolymerView {
	if idx < 0 || idx >= globula.Len() {
		return nil
	}
	return globula.polymers[idx]
}

func (globula *GlobulaView) Is(prop GlobulaProperty) bool {
	return globula.globulaProperties[prop]
}

func (globula *GlobulaView) Reset() {
	for _, pol := range globula.polymers {
		ForEachMonomer(pol, func(mon *dt.Monomer) bool { mon.MonomerType = dt.MonomerTypeUsual; return true })
	}
	for gp := range globula.globulaProperties {
		delete(globula.globulaProperties, gp)
	}
}

func (globula *GlobulaView) SetLiterals(literals map[dt.MonomerType]string) {
	for monType, str := range literals {
		globula.literalsTable[monType] = str
	}
}

func (globula *GlobulaView) GetLiteral(monType dt.MonomerType) string {
	return globula.literalsTable[monType]
}

func (globula *GlobulaView) GetLiterals() *map[dt.MonomerType]string {
	return &globula.literalsTable
}

func (globula *GlobulaView) GetMonomerTypeByLiteral(letter string) dt.MonomerType {
	for monType, l := range globula.literalsTable {
		if letter == l {
			return monType
		}
	}
	return dt.MonomerTypeUndefined
}

/*
FullReset The idea is to break all the connections and recover the original globula
using the polumer's connections information
*/
func (globula *GlobulaView) FullReset() {
	for _, pol := range globula.polymers {
		ForEachMonomer(pol, func(mon *dt.Monomer) bool {
			// Make monomer as usual
			mon.MonomerType = dt.MonomerTypeUsual
			// Break all the connections
			for _, side := range dt.GetAllSides() {
				sibling, err := mon.GetSibling(side)
				if err != nil || sibling == nil {
					continue
				}
				dt.BreakConnection(mon, sibling, side)
			}
			// Connection with the next monomer in the chain
			dt.MakeConnection(mon, mon.NextMonomer, dt.ConnectionTypeOne)
			return true
		})
	}
	for gp := range globula.globulaProperties {
		delete(globula.globulaProperties, gp)
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

func ForEachPolymerIf(globula *GlobulaView, pred func(*PolymerView) bool) {
	for _, pol := range globula.polymers {
		if !pred(pol) {
			break
		}
	}
}

func (globula *GlobulaView) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name       string
		Polymers   []*PolymerView
		Properties map[GlobulaProperty]bool
	}{
		Polymers:   globula.polymers,
		Properties: globula.globulaProperties,
	})
}

func (globula *GlobulaView) XClusters(avg float64) *ClusterView {
	if globula.xClusters == nil {
		globula.xClusters = NewClusterView(globula, avg, base.AxisX)
	}
	return globula.xClusters
}

func (globula *GlobulaView) YClusters(avg float64) *ClusterView {
	if globula.yClusters == nil {
		globula.yClusters = NewClusterView(globula, avg, base.AxisY)
	}
	return globula.yClusters
}

func (globula *GlobulaView) ZClusters(avg float64) *ClusterView {
	if globula.zClusters == nil {
		globula.zClusters = NewClusterView(globula, avg, base.AxisZ)
	}
	return globula.zClusters
}

func (globula *GlobulaView) CommonClusters() (*ClusterView, *ClusterView, *ClusterView) {
	if globula.Is(GlobulaGlobulaType) {
		return globula.commonClustersGlobula()
	} else {
		return globula.commonClustersThread()
	}
}

func (globula *GlobulaView) commonClustersGlobula() (*ClusterView, *ClusterView, *ClusterView) {
	if !globula.commonClusterDone {
		clustersX := globula.XClusters(0.0)
		clustersY := globula.YClusters(0.0)
		clustersZ := globula.ZClusters(0.0)

		ForEachCluster(clustersX, func(cluster_X *Cluster) {
			ForEachCluster(clustersY, func(cluster_Y *Cluster) {
				commonMonomers := IntersectClustersSoft(cluster_X, cluster_Y)
				for _, commonMonomer := range commonMonomers {
					ForEachCluster(clustersX, func(cluster *Cluster) {
						cluster.RemoveMonomer(commonMonomer)
					})
					ForEachCluster(clustersY, func(cluster *Cluster) {
						cluster.RemoveMonomer(commonMonomer)
					})
				}
			})
		})

		ForEachCluster(clustersY, func(cluster_Y *Cluster) {
			ForEachCluster(clustersZ, func(cluster_Z *Cluster) {
				commonMonomers := IntersectClustersSoft(cluster_Y, cluster_Z)
				for _, commonMonomer := range commonMonomers {
					ForEachCluster(clustersY, func(cluster *Cluster) {
						cluster.RemoveMonomer(commonMonomer)
					})
					ForEachCluster(clustersZ, func(cluster *Cluster) {
						cluster.RemoveMonomer(commonMonomer)
					})
				}
			})
		})

		ForEachCluster(clustersX, func(cluster_X *Cluster) {
			ForEachCluster(clustersZ, func(cluster_Z *Cluster) {
				commonMonomers := IntersectClustersSoft(cluster_X, cluster_Z)
				for _, commonMonomer := range commonMonomers {
					ForEachCluster(clustersX, func(cluster *Cluster) {
						cluster.RemoveMonomer(commonMonomer)
					})
					ForEachCluster(clustersZ, func(cluster *Cluster) {
						cluster.RemoveMonomer(commonMonomer)
					})
				}
			})
		})
		// In case we have empty cluster, we need to remove them
		clustersX.Trunkate()
		clustersY.Trunkate()
		clustersZ.Trunkate()

		// Everything is done. We're ready to fully-connect them
		for _, cluster := range []*ClusterView{clustersX, clustersY, clustersZ} {
			ForEachCluster(cluster, func(c *Cluster) {
				c.MakeFullyConnected()
			})
		}
		// Now let's colorize them
		clustersX.Colorize(false)
		clustersY.Colorize(false)
		clustersZ.Colorize(false)

		globula.commonClusterDone = true
	}
	return globula.xClusters, globula.yClusters, globula.zClusters
}

func (globula *GlobulaView) commonClustersThread() (*ClusterView, *ClusterView, *ClusterView) {
	if !globula.commonClusterDone {
		size := globula.polymers[0].Len()
		clusterLength := int(float64(size) * 0.2)
		firstMonomer := globula.getFirstMonomer()
		globula.zClusters = globula.getClustersInThread(firstMonomer, clusterLength)
	}
	return globula.xClusters, globula.yClusters, globula.zClusters
}

func (globula *GlobulaView) getFirstMonomer() *dt.Monomer {
	polymer := globula.polymers[0]
	current := polymer.polymer.GetMonomerByIdx(polymer.polymer.Len() / 2)
	return getBorderMonomer(current, dt.SideForward)
}

func getBorderMonomer(startingMonomer *dt.Monomer, side dt.Side) *dt.Monomer {
	current := startingMonomer
	sibling, err := current.GetSibling(side)
	for err == nil && sibling.MonomerType != dt.MonomerTypeUndefined {
		current = sibling
		sibling, err = current.GetSibling(side)
	}
	return current
}

func (globula *GlobulaView) getClustersInThread(firstMonomer *dt.Monomer, clusterLength int) *ClusterView {
	clusters := make([]*Cluster, 0)
	startingMonomers := make([]*dt.Monomer, 0)
	for i := 0; i < clusterLength; i++ {
		if i != 0 {
			firstMonomer, _ = firstMonomer.GetSibling(dt.SideUp)
			if firstMonomer == nil || firstMonomer.IsTypeOf(dt.MonomerTypeUndefined) {
				break
			}
		}
		startingMonomer := getBorderMonomer(firstMonomer, dt.SideLeft)
		for {
			cluster := new(Cluster)
			if base.Contains(startingMonomers, startingMonomer) {
				startingMonomer, _ = startingMonomer.GetSibling(dt.SideRight)
				if startingMonomer == nil || startingMonomer.IsTypeOf(dt.MonomerTypeUndefined) {
					startingMonomer, _ = getBorderMonomer(startingMonomer, dt.SideLeft).GetSibling(dt.SideBackward)
					if startingMonomer == nil || startingMonomer.IsTypeOf(dt.MonomerTypeUndefined) {
						break
					}
				}
			}
			startingMonomers = append(startingMonomers, startingMonomer)
			mon2, _ := startingMonomer.GetSibling(dt.SideRight)
			if mon2.IsTypeOf(dt.MonomerTypeUndefined) {
				continue
			}
			mon3, _ := mon2.GetSibling(dt.SideBackward)
			if mon3.IsTypeOf(dt.MonomerTypeUndefined) {
				continue
			}
			mon4, _ := mon3.GetSibling(dt.SideLeft)
			if mon4.IsTypeOf(dt.MonomerTypeUndefined) {
				continue
			}
			mon11, _ := startingMonomer.GetSibling(dt.SideUp)
			mon21, _ := mon2.GetSibling(dt.SideUp)
			mon31, _ := mon3.GetSibling(dt.SideUp)
			mon41, _ := mon4.GetSibling(dt.SideUp)
			clusterUnit := NewClusterUnit([]*dt.Monomer{
				startingMonomer, mon2, mon3, mon4,
				mon11, mon21, mon31, mon41,
			}, dt.SideUp, base.AxisZ)
			clusterUnit.MakeFullyConnected()
			cluster.units = append(cluster.units, clusterUnit)
			clusters = append(clusters, cluster)
		}
	}
	return NewClusterViewRaw(clusters, base.AxisZ)
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
	Cs_ := globula.breakConnections(int(float64(groupsCount)*0.26), dt.MonomerTypeOContaining)

	// 2. Turn some Bs into C
	turnIntoAnotherGroup(&Cs_, int(float64(groupsCount)*0.03), dt.MonomerTypeVynil)

	// 3. Turn random bins into Cs (excluding Bs)
	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// Do it twice
	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))

	// 4. Create connections
	globula.createCrosslinks1(int(float64(groupsCount)*0.59), &Cs_)
	globula.globulaProperties[GlobulaAged] = true
}

func (globula *GlobulaView) DoAging2(groupsCount int, doCrosslinks bool) {
	// 1. Break connections
	Bs_ := globula.breakConnections(int(float64(groupsCount)*0.44), dt.MonomerTypeVynil)

	// 2. Turn some Bs into C
	turnIntoAnotherGroup(&Bs_, int(float64(groupsCount)*0.15), dt.MonomerTypeVynil)

	// 3. Turn random bins into Cs (excluding Bs)
	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))
	// Do it twice
	globula.turnRandomBinsIntoC(int(float64(groupsCount) * 0.06))

	// 4. Create connections
	if doCrosslinks {
		globula.createCrosslinks2(int(float64(groupsCount) * crosslinksCount))
	}
	globula.globulaProperties[GlobulaAged] = true
}

func (globula *GlobulaView) breakConnections(groupsCount int, monomerTypeToGather dt.MonomerType) []*dt.Monomer {
	Bs := make([]*dt.Monomer, 0)
	triesNumber := 0
	for len(Bs) != groupsCount && triesNumber < groupsCount {
		chosenPoly := rand.Intn(globula.Len()) // Take a random poly
		poly := globula.polymers[chosenPoly]
		if poly.Len() < 2 {
			continue
		}
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		nextMonomerNumber := chosenMonomerNumber + 1
		chosenMonomer := poly.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != dt.MonomerTypeUsual {
			triesNumber++
			continue
		}
		nextMonomer := poly.polymer.GetMonomerByIdx(nextMonomerNumber)
		if nextMonomer.MonomerType != dt.MonomerTypeUsual {
			triesNumber++
			continue
		}
		triesNumber = 0
		side := chosenMonomer.GetSideOfSibling(nextMonomer)
		dt.TierConnection(chosenMonomer, nextMonomer, side)
		if rand.Intn(2) == 0 {
			chosenMonomer.MonomerType = dt.MonomerTypeOContaining
			nextMonomer.MonomerType = dt.MonomerTypeVynil
			if monomerTypeToGather == dt.MonomerTypeOContaining {
				Bs = append(Bs, chosenMonomer)
			} else {
				Bs = append(Bs, nextMonomer)
			}
		} else {
			chosenMonomer.MonomerType = dt.MonomerTypeVynil
			nextMonomer.MonomerType = dt.MonomerTypeOContaining
			if monomerTypeToGather == dt.MonomerTypeVynil {
				Bs = append(Bs, chosenMonomer)
			} else {
				Bs = append(Bs, nextMonomer)
			}
		}
	}

	return Bs
}

func turnIntoAnotherGroup(Bs *[]*dt.Monomer, groupsCount int, monomerTypeToTurn dt.MonomerType) {
	mapUsedBs := make(map[int]bool)
	triesCount := 0
	for len(mapUsedBs) != groupsCount && triesCount < 3*groupsCount {
		mapUsedBs[rand.Intn(len(*Bs))] = true
		triesCount++
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
	triesNumber := 0
	for i < groupsCount && triesNumber < groupsCount {
		chosenPoly := rand.Intn(globula.Len()) // Take a random poly
		poly := globula.polymers[chosenPoly]
		if poly.Len() < 2 {
			triesNumber++
			continue
		}
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType == dt.MonomerTypeUsual {
			chosenMonomer.MonomerType = dt.MonomerTypeVynil
			i++
			triesNumber = 0
		} else {
			triesNumber++
		}
	}
}

func (globula *GlobulaView) createCrosslinks1(groupsCount int, Bs *[]*dt.Monomer) {
	for len(*Bs) != groupsCount {
		chosenPoly := rand.Intn(globula.Len()) // Take a random poly
		poly := globula.polymers[chosenPoly]
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != dt.MonomerTypeUsual {
			continue
		}

		movementSides := dt.GetMovementSides()
		for i := 0; i < len(movementSides); i++ {
			chosenSide := movementSides[rand.Intn(len(movementSides))]
			nextMonomer, err := chosenMonomer.GetSibling(chosenSide)
			if err == nil && nextMonomer != nil && nextMonomer.MonomerType == dt.MonomerTypeUsual {
				dt.MakeConnection(chosenMonomer, nextMonomer, dt.ConnectionTypeCrosslinks)
				chosenMonomer.MonomerType = dt.MonomerTypeOContaining
				nextMonomer.MonomerType = dt.MonomerTypeOContaining
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
		if poly.Len() < 2 {
			continue
		}
		chosenMonomerNumber := rand.Intn(poly.Len() - 1) // Take a random monomer in it
		chosenMonomer := poly.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != dt.MonomerTypeUsual {
			timesRepeated++
			continue
		}

		movementSides := dt.GetMovementSides()
		for i := 0; i < len(movementSides); i++ {
			chosenSide := movementSides[rand.Intn(len(movementSides))]
			nextMonomer, err := chosenMonomer.GetSibling(chosenSide)
			if err == nil &&
				nextMonomer != nil &&
				nextMonomer.MonomerType == dt.MonomerTypeUsual &&
				(!dt.MonomersAreEqual(chosenMonomer.NextMonomer, nextMonomer) &&
					!dt.MonomersAreEqual(chosenMonomer.PrevMonomer, nextMonomer)) {
				dt.MakeConnection(chosenMonomer, nextMonomer, dt.ConnectionTypeCrosslinks)
				chosenMonomer.MonomerType = dt.MonomerTypeCrosslinked
				nextMonomer.MonomerType = dt.MonomerTypeCrosslinked
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

func (globula *GlobulaView) DoAging3(ncut, OcontainingCount, ncross int) error {
	if ncut < OcontainingCount {
		return errors.New("ncut is less that the O-containing monomers count")
	}
	printer := outputformat.GetPrint()
	// First, distribute O containing monomers
	// if warning, err := globula.aging3DistributeCutMonomers(OcontainingCount,
	// 	dt.MONOMER_TYPE_O_CONTAINING,
	// 	dt.MONOMER_TYPE_VYNIL); warning != nil {
	// 	printer.PrintlnWarning(warning.Error())
	// } else if err != nil {
	// 	return err
	// }
	globula.breakConnections(OcontainingCount, dt.MonomerTypeVynil)

	// Then, distribute what's left
	// if warning, err := globula.aging3DistributeCutMonomers(ncut-OcontainingCount,
	// 	dt.MONOMER_TYPE_VYNIL,
	// 	dt.MONOMER_TYPE_VYNIL); warning != nil {
	// 	printer.PrintlnWarning(warning.Error())
	// } else if err != nil {
	// 	return err
	// }
	globula.turnRandomBinsIntoC(ncut - OcontainingCount)

	// At the end, distribute crosslinks
	if warning, err := globula.aging3DistributeCrosslinks(ncross); warning != nil {
		printer.PrintlnWarning(warning.Error())
	} else if err != nil {
		return err
	}
	return nil
}

func (globula *GlobulaView) aging3DistributeCutMonomers(OcontainingCount int, typePrev, typeNext dt.MonomerType) (warning, err error) {
	warning = nil
	err = nil
	trialsCount := 0
	for OcontainingCount > 0 && trialsCount < 10 {
		chosenPolymerNumber := rand.Intn(globula.Len())
		chosenPolymer := globula.polymers[chosenPolymerNumber]
		chosenMonomerNumber := rand.Intn(chosenPolymer.Len()-2) + 1
		chosenMonomer := chosenPolymer.polymer.GetMonomerByIdx(chosenMonomerNumber)
		nextChosenMonomer := chosenPolymer.polymer.GetMonomerByIdx(chosenMonomerNumber + 1)
		if chosenMonomer.MonomerType != dt.MonomerTypeUsual ||
			nextChosenMonomer.MonomerType != dt.MonomerTypeUsual {
			trialsCount++
			continue
		}
		dt.BreakConnection1(chosenMonomer, nextChosenMonomer)
		chosenMonomer.MonomerType = typePrev
		nextChosenMonomer.MonomerType = typeNext
		OcontainingCount--
		trialsCount = 0
	}
	if trialsCount != 0 {
		return fmt.Errorf("%d O containing monomers weren't distributed", OcontainingCount), nil
	}
	return nil, nil
}

func (globula *GlobulaView) aging3DistributeCrosslinks(ncross int) (warning, err error) {
	warning = nil
	err = nil
	trialsCount := 0
	for ncross > 0 && trialsCount < 10 {
		chosenPolymerNumber := rand.Intn(globula.Len())
		chosenPolymer := globula.polymers[chosenPolymerNumber]
		chosenMonomerNumber := rand.Intn(chosenPolymer.Len()-2) + 1
		chosenMonomer := chosenPolymer.polymer.GetMonomerByIdx(chosenMonomerNumber)
		nextChosenMonomer := chosenPolymer.polymer.GetMonomerByIdx(chosenMonomerNumber + 1)
		if chosenMonomer.MonomerType != dt.MonomerTypeUsual ||
			nextChosenMonomer.MonomerType != dt.MonomerTypeUsual {
			trialsCount++
			continue
		}
		dt.MakeConnection(chosenMonomer, nextChosenMonomer, dt.ConnectionTypeCrosslinks)
		chosenMonomer.MonomerType = dt.MonomerTypeCrosslinked
		nextChosenMonomer.MonomerType = dt.MonomerTypeCrosslinked
		ncross--
		trialsCount = 0
	}
	if trialsCount != 0 {
		return fmt.Errorf("%d crosslinks weren't distributed", ncross), nil
	}
	return nil, nil
}

func (globula *GlobulaView) aging4DistributeCrosslinks(ncross int) (warning, err error) {
	warning = nil
	err = nil
	trialsCount := 0
	for ncross > 0 && trialsCount < 10 {
		chosenPolymerNumber := rand.Intn(globula.Len())
		chosenPolymer := globula.polymers[chosenPolymerNumber]
		chosenMonomerNumber := rand.Intn(chosenPolymer.Len())
		chosenMonomer := chosenPolymer.polymer.GetMonomerByIdx(chosenMonomerNumber)
		if chosenMonomer.MonomerType != dt.MonomerTypeUsual {
			trialsCount++
			continue
		}
		movementSides := dt.GetMovementSides()
		for i := 0; i < len(movementSides); i++ {
			chosenSide := movementSides[rand.Intn(len(movementSides))]
			nextMonomer, _ := chosenMonomer.GetSibling(chosenSide)
			getBorderMonomer := func(side dt.Side, startMonomer *dt.Monomer) *dt.Monomer {
				curr := startMonomer
				prev, _ := curr.GetSibling(side)
				for prev != nil && prev.MonomerType != dt.MonomerTypeUndefined {
					t := prev
					prev, _ = curr.GetSibling(side)
					curr = t
				}
				return curr
			}
			if (chosenMonomer.NextMonomer != nil && dt.MonomersAreEqual(chosenMonomer.NextMonomer, nextMonomer)) ||
				(chosenMonomer.PrevMonomer != nil && dt.MonomersAreEqual(chosenMonomer.PrevMonomer, nextMonomer)) {
				continue
			}
			if nextMonomer == nil ||
				nextMonomer.MonomerType == dt.MonomerTypeUndefined {
				switch chosenSide {
				case dt.SideForward:
					nextMonomer = getBorderMonomer(dt.SideBackward, chosenMonomer)
				case dt.SideBackward:
					nextMonomer = getBorderMonomer(dt.SideForward, chosenMonomer)
				case dt.SideUp:
					nextMonomer = getBorderMonomer(dt.SideDown, chosenMonomer)
				case dt.SideDown:
					nextMonomer = getBorderMonomer(dt.SideUp, chosenMonomer)
				default:
					continue
				}
				dt.MakeConnectionUnsafe(chosenMonomer, nextMonomer, chosenSide, dt.ConnectionTypeCrosslinks)
			} else {
				dt.MakeConnection(chosenMonomer, nextMonomer, dt.ConnectionTypeCrosslinks)
			}
			chosenMonomer.MonomerType = dt.MonomerTypeCrosslinked
			nextMonomer.MonomerType = dt.MonomerTypeCrosslinked
			trialsCount = 0
			break
		}
	}
	if trialsCount != 0 {
		return fmt.Errorf("%d crosslinks weren't distributed", ncross), nil
	}
	return nil, nil
}

func (globula *GlobulaView) HighlightBorders() bool {
	for _, polymerView := range globula.polymers {
		go doDST(polymerView)
	}
	return true
}

func doDST(polymerView *PolymerView) {

}

func (globula *GlobulaView) Waterize() {
	field := globula.polymers[0].GetUnderlinedField().(*dt.Field)
	field.Waterize()
	globula.globulaProperties[GlobulaWaterized] = true
}

func (globula *GlobulaView) MakeHomogenousAsShortest() error {
	if globula.Len() == 0 {
		return errors.New("the globula is empty")
	}
	// find the length of the shortest chain
	shortestChainLength := globula.polymers[0].Len()
	for _, polymer := range globula.polymers {
		if polymer.Len() < shortestChainLength {
			shortestChainLength = polymer.Len()
		}
	}

	// Trunkate others to the minimum found
	return globula.trunkAllPolymers(shortestChainLength)
}

func (globula *GlobulaView) MakeHomogenousAsCustom(newSize int) error {
	if newSize > 128 {
		return errors.New("new size must not be more than 128")
	}

	return globula.trunkAllPolymers(newSize)
}

func (globula *GlobulaView) trunkAllPolymers(newSize int) error {
	for _, polymer := range globula.polymers {
		polymer.TrunkTo(newSize)
	}

	var monomerNumber int64 = 1
	for _, polymer := range globula.polymers {
		ForEachMonomer(polymer, func(m *dt.Monomer) bool {
			m.Number = int64(monomerNumber)
			monomerNumber++
			return true
		})
	}

	return nil
}

func (globula *GlobulaView) GetStatistics() string {
	builder := strings.Builder{}
	builder.WriteString(globula.showNumberOfParticles())
	builder.WriteString(globula.showNumberOfChains())
	builder.WriteString(globula.showTheoreticalAgeStatistics())
	builder.WriteString(globula.showActualAgeStatistics())
	builder.WriteString(globula.showMeanLengthOfChains())
	return builder.String()
}

func (globula *GlobulaView) showNumberOfParticles() string {
	builder := strings.Builder{}
	builder.WriteString("1. Число частиц: ")
	builder.WriteString(strconv.Itoa(globula.GetAtomsCount()))
	builder.WriteString("\n")
	return builder.String()
}

func (globula *GlobulaView) GetAtomsCount() int {
	atomsCount := 0
	for _, pol := range globula.polymers {
		atomsCount += pol.Len()
	}
	return atomsCount
}

func (globula *GlobulaView) showNumberOfChains() string {
	builder := strings.Builder{}
	builder.WriteString("2. Исходное число цепей ")
	builder.WriteString(strconv.Itoa(globula.Len()))
	builder.WriteString("\n")
	return builder.String()
}

func (globula *GlobulaView) showTheoreticalAgeStatistics() string {
	outputformat.GetPrint().Println("Please, specify the theoretical age ratio in percents")
	ageRatioPercent := 25
	ageRatio := float64(ageRatioPercent) * 0.01
	atomsCount := globula.GetAtomsCount()
	outputformat.GetPrint().Readln(&ageRatioPercent)
	builder := strings.Builder{}

	builder.WriteString("3. Ожидаемая степень старения: ")
	builder.WriteString(strconv.Itoa(int(float64(atomsCount) * ageRatio)))
	builder.WriteString(" (")
	builder.WriteString(strconv.Itoa(ageRatioPercent))
	builder.WriteString("%)\n")

	agedParticlesCount := math.Ceil(float64(atomsCount) * ageRatio)
	cutsCount := int(math.Ceil(agedParticlesCount * 0.44))
	crossCount := int(agedParticlesCount * crosslinksCount)

	builder.WriteString("4. Ожидаемое количество разрывов: ")
	builder.WriteString(strconv.Itoa(cutsCount))
	builder.WriteString("\n")

	builder.WriteString("   Ожидаемое количество сшивок: ")
	builder.WriteString(strconv.Itoa(crossCount))
	builder.WriteString("\n")

	CCount := int(agedParticlesCount * 0.59)
	NCount := int(agedParticlesCount) - CCount
	HCount := crossCount * 2

	builder.WriteString("5. Ожидаемое распределение по состаренным группам:\n")
	builder.WriteString("   C: ")
	builder.WriteString(strconv.Itoa(CCount))
	builder.WriteString("\n")
	builder.WriteString("   N: ")
	builder.WriteString(strconv.Itoa(NCount))
	builder.WriteString("\n")
	builder.WriteString("   H: ")
	builder.WriteString(strconv.Itoa(HCount))
	builder.WriteString("\n")
	return builder.String()
}

func (globula *GlobulaView) showActualAgeStatistics() string {
	NCount := 0
	CCount := 0
	cutsCount := 0
	crossCount := make(map[int64]int)
	for _, pol := range globula.polymers {
		ForEachMonomer(pol, func(mon *dt.Monomer) bool {
			if mon.NextMonomer != nil &&
				((mon.MonomerType == dt.MonomerTypeVynil || mon.MonomerType == dt.MonomerTypeOContaining) &&
					(mon.NextMonomer.MonomerType == dt.MonomerTypeVynil ||
						mon.NextMonomer.MonomerType == dt.MonomerTypeOContaining)) {
				cutsCount++
			}
			if mon.MonomerType == dt.MonomerTypeVynil {
				CCount++
			} else if mon.MonomerType == dt.MonomerTypeOContaining {
				NCount++
			} else if _, ok := crossCount[mon.Number]; !ok && mon.MonomerType == dt.MonomerTypeCrosslinked {
				crossCount[mon.Number] = 1
				for _, side := range dt.GetMovementSides() {
					connectionType := mon.GetTypeOfConnectionWithSide(side)
					if connectionType == dt.ConnectionTypeCrosslinks {
						crossMon, _ := mon.GetSibling(side)
						crossCount[crossMon.Number] = 1
					}
				}
			} else {
				return true
			}
			return true
		})
	}

	builder := strings.Builder{}
	ageGroupsCount := float64(NCount + CCount)
	builder.WriteString("6. Фактическое количество разрывов: ")
	builder.WriteString(strconv.Itoa(cutsCount))
	builder.WriteString("\n")
	builder.WriteString("   Фактическое количество сшивок: ")
	builder.WriteString(strconv.Itoa(len(crossCount) / 2))
	builder.WriteString("\n")
	builder.WriteString("7. Фактическое распределение по состаренным группам:\n")
	builder.WriteString("   C: ")
	builder.WriteString(strconv.Itoa(CCount) + " (" + strconv.FormatFloat(float64(CCount)/ageGroupsCount, 'f', 2, 64) + "%)")
	builder.WriteString("\n")
	builder.WriteString("   N: ")
	builder.WriteString(strconv.Itoa(NCount) + " (" + strconv.FormatFloat(float64(NCount)/ageGroupsCount, 'f', 2, 64) + "%)")
	builder.WriteString("\n")
	builder.WriteString("   H: ")
	builder.WriteString(strconv.Itoa(len(crossCount)))
	builder.WriteString("\n")
	builder.WriteString("   Фактическая степень старения: ")
	builder.WriteString(strconv.FormatFloat(ageGroupsCount/float64(globula.GetAtomsCount())*100.0, 'f', 2, 64))
	builder.WriteString("%\n")
	return builder.String()
}

func (globula *GlobulaView) showMeanLengthOfChains() string {
	builder := strings.Builder{}
	builder.WriteString("8. Средняя длина цепи: ")
	atomsCount := globula.GetAtomsCount()
	builder.WriteString(strconv.FormatFloat(float64(atomsCount)/float64(globula.Len()), 'f', 2, 64))
	return builder.String()
}

func (globula *GlobulaView) DoAgingSurface(ncut, nOContaining, ncross int) error {
	if ncut < nOContaining {
		return errors.New("ncut is less that the O-containing monomers count")
	}
	printer := outputformat.GetPrint()
	// First, distribute O containing monomers
	globula.breakConnections(nOContaining, dt.MonomerTypeVynil)

	// Then, distribute what's left
	globula.turnRandomBinsIntoC(ncut - nOContaining)

	// At the end, distribute crosslinks
	if warning, err := globula.aging4DistributeCrosslinks(ncross); warning != nil {
		printer.PrintlnWarning(warning.Error())
	} else if err != nil {
		return err
	}
	return nil
}
