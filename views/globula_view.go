package views

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"polymers/base"
	dt "polymers/datatypes"
	"polymers/outputformat"
	"slices"
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
	GlobulaPatternType
)

type GlobulaView struct {
	polymers          []*PolymerView
	xClusters         *ClusterView
	yClusters         *ClusterView
	zClusters         *ClusterView
	commonClusterDone bool
	globulaProperties map[GlobulaProperty]bool
}

func NewGlobulaView(polymers []dt.IPolymer, globulaType GlobulaProperty) *GlobulaView {
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
		ForEachMonomer(pol, func(mon *dt.Monomer) bool { mon.MonomerType = base.C; return true })
	}
	for gp := range globula.globulaProperties {
		if gp == GlobulaGlobulaType || gp == GlobulaThreadType || gp == GlobulaSurfaceType {
			continue
		}
		delete(globula.globulaProperties, gp)
	}
}

/*
FullReset The idea is to break all the connections and recover the original globula
using the polumer's connections information
*/
func (globula *GlobulaView) FullReset() {
	for _, pol := range globula.polymers {
		ForEachMonomer(pol, func(mon *dt.Monomer) bool {
			// Make monomer as usual
			mon.MonomerType = base.C
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
// 	pols_countMap := make(map[string]any)
// 	pols_countMap["value"] = len(globula.polymers)
// 	pols_countMap["name"] = "Количество полимеров"

// 	pols_map := make(map[string]any)
// 	pols_map["value"] = make([]string, len(globula.polymers))

// 	var jsonDict map[string]any
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
	for err == nil && sibling.IsNotTypeOf(base.MendeleevTableElementUndefined) {
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
			if firstMonomer == nil || firstMonomer.IsTypeOf(base.MendeleevTableElementUndefined) {
				break
			}
		}
		startingMonomer := getBorderMonomer(firstMonomer, dt.SideLeft)
		for {
			cluster := new(Cluster)
			if base.Contains(startingMonomers, startingMonomer) {
				startingMonomer, _ = startingMonomer.GetSibling(dt.SideRight)
				if startingMonomer == nil || startingMonomer.IsTypeOf(base.MendeleevTableElementUndefined) {
					startingMonomer, _ = getBorderMonomer(startingMonomer, dt.SideLeft).GetSibling(dt.SideBackward)
					if startingMonomer == nil || startingMonomer.IsTypeOf(base.MendeleevTableElementUndefined) {
						break
					}
				}
			}
			startingMonomers = append(startingMonomers, startingMonomer)
			mon2, _ := startingMonomer.GetSibling(dt.SideRight)
			if mon2.IsTypeOf(base.MendeleevTableElementUndefined) {
				continue
			}
			mon3, _ := mon2.GetSibling(dt.SideBackward)
			if mon3.IsTypeOf(base.MendeleevTableElementUndefined) {
				continue
			}
			mon4, _ := mon3.GetSibling(dt.SideLeft)
			if mon4.IsTypeOf(base.MendeleevTableElementUndefined) {
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
	builder.WriteString(globula.showLengthDistribution())
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
	//outputformat.GetPrint().Readln(&ageRatioPercent)
	builder := strings.Builder{}

	builder.WriteString("3. Ожидаемая степень старения: ")
	builder.WriteString(strconv.Itoa(int(float64(atomsCount) * ageRatio)))
	builder.WriteString(" (")
	builder.WriteString(strconv.Itoa(ageRatioPercent))
	builder.WriteString("%)\n")

	agedParticlesCount := math.Ceil(float64(atomsCount) * ageRatio)
	cutsCount := int(math.Ceil(agedParticlesCount * 0.44))
	crossCount := int(agedParticlesCount * 0.25)

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
		for monNumber := range pol.Len() {
			mon := pol.GetMonomerByIdx(monNumber)
			if mon.NextMonomer != nil &&
				((mon.MonomerType == base.N || mon.MonomerType == base.O) &&
					(mon.NextMonomer.MonomerType == base.N ||
						mon.NextMonomer.MonomerType == base.O)) {
				cutsCount++
			}
			if mon.MonomerType == base.N {
				CCount++
			} else if mon.MonomerType == base.O {
				NCount++
			} else if _, ok := crossCount[mon.Number]; !ok && mon.MonomerType == base.H {
				crossCount[mon.Number] = 1
				for _, side := range dt.GetMovementSides() {
					connectionType := mon.GetTypeOfConnectionWithSide(side)
					if connectionType == dt.ConnectionTypeCrosslinks {
						crossMon, _ := mon.GetSibling(side)
						crossCount[crossMon.Number] = 1
					}
				}
			}
		}
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
	builder.WriteString(strconv.Itoa(CCount))
	builder.WriteString(" (")
	if ageGroupsCount != 0 {
		builder.WriteString(strconv.FormatFloat(float64(CCount)/ageGroupsCount*100, 'f', 2, 64))
	} else {
		builder.WriteString("0")
	}
	builder.WriteString("%)")
	builder.WriteString("\n")
	builder.WriteString("   N: ")
	builder.WriteString(strconv.Itoa(NCount))
	builder.WriteString(" (")
	if ageGroupsCount != 0 {
		builder.WriteString(strconv.FormatFloat(float64(NCount)/ageGroupsCount*100, 'f', 2, 64))
	} else {
		builder.WriteString("0")
	}
	builder.WriteString("%)")
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
	builder.WriteString("\n")
	return builder.String()
}

func (globula *GlobulaView) showLengthDistribution() string {
	builder := strings.Builder{}
	distribution := make(map[int]int) // [length]count_of_such_chains
	for _, polymer := range globula.polymers {
		currChainLength := 0
		for monNumber := range polymer.Len() {
			currMon := polymer.GetMonomerByIdx(monNumber)
			currChainLength++
			if currMon.NextMonomer == nil {
				distribution[currChainLength]++
				currChainLength = 0
			}
		}
	}
	lengths := make([]int, len(distribution))
	i := 0
	for length := range distribution {
		lengths[i] = length
		i++
	}
	slices.Sort(lengths)
	builder.WriteString("9. Распределение по длинам цепей\n")
	builder.WriteString("Длина: Количество\n")
	for _, length := range lengths {
		count := distribution[length]
		fmt.Fprintf(&builder, "%d: %d\n", length, count)
	}
	builder.WriteString("\n")
	return builder.String()
}

func (globula *GlobulaView) DeepCopy() *GlobulaView {
	newGlobula := &GlobulaView{}
	newGlobula.polymers = make([]*PolymerView, len(globula.polymers))
	if len(globula.polymers) > 0 {
		firstPolymer := globula.polymers[0]
		field := firstPolymer.GetUnderlinedField().DeepCopy()
		for i, pol := range globula.polymers {
			newGlobula.polymers[i] = pol.DeepCopy(field)
		}
	}
	newGlobula.globulaProperties = make(map[GlobulaProperty]bool)
	maps.Copy(newGlobula.globulaProperties, globula.globulaProperties)
	return newGlobula
}

func (globula *GlobulaView) SetProperty(property GlobulaProperty) {
	globula.globulaProperties[property] = true
}
