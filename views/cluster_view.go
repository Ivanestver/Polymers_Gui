package views

import (
	"polymers/base"
	dt "polymers/datatypes"
)

const _ClusterUnitSize int = 8

type ClusterUnit struct {
	monomers      [_ClusterUnitSize]*dt.Monomer
	mainDirection dt.Side
	axis          base.Axis
}

func NewClusterUnit(monomers []*dt.Monomer, mainDirection dt.Side, axis base.Axis) *ClusterUnit {
	if len(monomers) != _ClusterUnitSize {
		return nil
	}
	newClusterUnit := new(ClusterUnit)
	for i := 0; i < _ClusterUnitSize; i++ {
		newClusterUnit.monomers[i] = monomers[i]
	}
	newClusterUnit.mainDirection = mainDirection
	newClusterUnit.axis = axis
	return newClusterUnit
}

func (clusterUnit *ClusterUnit) Size() int {
	return _ClusterUnitSize
}

func (clusterUnit *ClusterUnit) MakeFullyConnected() {
	for i := 0; i < clusterUnit.Size(); i++ {
		currMon := clusterUnit.monomers[i]
		for j := 0; j < clusterUnit.Size(); j++ {
			if i == j {
				continue
			}
			sideMon := clusterUnit.monomers[j]
			side := dt.GetSideByMonomers(currMon, sideMon)
			if side != dt.SideUndefined {
				dt.MakeConnection(currMon, sideMon, dt.GetConnectionType(currMon, sideMon))
			}
		}
	}
}

func (clusterUnit *ClusterUnit) contains(monomer *dt.Monomer) bool {
	for _, mon := range clusterUnit.monomers {
		if dt.MonomersAreEqual(mon, monomer) {
			return true
		}
	}
	return false
}

func (clusterUnit *ClusterUnit) setTypeOfMonomers(monomerType base.MendeleevTableElement) {
	for i := 0; i < clusterUnit.Size(); i++ {
		currMon := clusterUnit.monomers[i]
		currMon.MonomerType = monomerType
	}
}

func clusterUnitsAreEqual(unit1, unit2 *ClusterUnit) bool {
	if unit1 == nil && unit2 == nil {
		return true
	}
	if unit1 == nil || unit2 == nil {
		return false
	}

	if unit1.axis != unit2.axis &&
		unit1.mainDirection != unit2.mainDirection &&
		unit1.mainDirection != dt.GetReversedSide(unit2.mainDirection) {
		return false
	}

	for _, mon := range unit1.monomers {
		contains := false
		for _, m := range unit2.monomers {
			if dt.MonomersAreEqual(mon, m) {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true
}

type Cluster struct {
	units         []*ClusterUnit
	mainDirection dt.Side
	axis          base.Axis
}

func NewCluster(unit []*ClusterUnit, mainDirection dt.Side, axis base.Axis) *Cluster {
	newCluster := new(Cluster)
	newCluster.units = unit
	newCluster.mainDirection = mainDirection
	newCluster.axis = axis
	return newCluster
}

func (cluster *Cluster) Size() int {
	return len(cluster.units)
}

func (cluster *Cluster) MainDirection() dt.Side {
	return cluster.mainDirection
}

func (cluster *Cluster) Axis() base.Axis {
	return cluster.axis
}

func (cluster *Cluster) SetTypeOfMonomers(monomerType base.MendeleevTableElement) {
	for _, unit := range cluster.units {
		unit.setTypeOfMonomers(monomerType)
	}
}

func (cluster *Cluster) RemoveMonomer(monomer *dt.Monomer) bool {
	for i, unit := range cluster.units {
		if unit.contains(monomer) {
			cluster.units = append(cluster.units[:i], cluster.units[i+1:]...)
			return true
		}
	}
	return false
}

func containsMonomer(usedMonomers *[]*dt.Monomer, mon *dt.Monomer) bool {
	for _, usedMon := range *usedMonomers {
		if dt.MonomersAreEqual(usedMon, mon) {
			return true
		}
	}
	return false
}

func (cluster *Cluster) GetAvgLengthByAxis() float64 {
	allMonomers := make([]*dt.Monomer, 0)
	for _, unit := range cluster.units {
		for _, mon := range unit.monomers {
			allMonomers = append(allMonomers, mon)
		}
	}
	sideBackward := dt.GetReversedSide(cluster.mainDirection)

	lengths := make([]float64, 0)
	usedMonomers := make([]*dt.Monomer, 0)
	for _, mon := range allMonomers {
		if containsMonomer(&usedMonomers, mon) {
			continue
		}

		currMonomer := mon
		sibling, err := currMonomer.GetSibling(sideBackward)
		for err == nil && containsMonomer(&allMonomers, sibling) {
			currMonomer = sibling
			sibling, err = currMonomer.GetSibling(sideBackward)
		}

		usedMonomers = append(usedMonomers, currMonomer)
		length := 1.0
		sibling, err = currMonomer.GetSibling(cluster.mainDirection)
		for err == nil && containsMonomer(&allMonomers, sibling) {
			currMonomer = sibling
			usedMonomers = append(usedMonomers, currMonomer)
			length += 1.0
			sibling, err = currMonomer.GetSibling(cluster.mainDirection)
		}
		lengths = append(lengths, length)
	}

	return base.Sum(lengths) / float64(len(lengths))
}

func (cluster *Cluster) MakeFullyConnected() {
	for _, unit := range cluster.units {
		unit.MakeFullyConnected()
	}
}

func JoinClusters(cluster1, cluster2 *Cluster) *Cluster {
	if cluster1.mainDirection != cluster2.mainDirection &&
		cluster1.mainDirection != dt.GetReversedSide(cluster2.mainDirection) &&
		cluster1.axis != cluster2.axis {
		return nil
	}

	intersection := IntersectClustersStickToDirection(cluster1, cluster2)

	// If there's nothing to connect, these are probably different clusters
	if len(intersection) < 2 {
		return nil
	}

	units := make([]*ClusterUnit, 0)
	units = append(units, cluster1.units...)

	for _, unit2 := range cluster2.units {
		contains := false
		for _, unit := range units {
			if clusterUnitsAreEqual(unit, unit2) {
				contains = true
				break
			}
		}
		if !contains {
			units = append(units, unit2)
		}
	}

	return NewCluster(units, cluster1.mainDirection, cluster1.axis)
}

func intersectClusterUnits(unit1, unit2 *ClusterUnit) []*dt.Monomer {
	m := make(map[*dt.Monomer]bool)
	for _, mon := range unit1.monomers {
		m[mon] = true
	}
	intersection := make([]*dt.Monomer, 0)
	for _, mon := range unit2.monomers {
		if _, ok := m[mon]; ok {
			intersection = append(intersection, mon)
		}
	}
	return intersection
}

func IntersectClustersSoft(cluster1, cluster2 *Cluster) []*dt.Monomer {
	intersectedClusterUnitsMap := make(map[*dt.Monomer]bool, 0)
	for i := 0; i < cluster1.Size(); i++ {
		for j := 0; j < cluster2.Size(); j++ {
			inters := intersectClusterUnits(cluster1.units[i], cluster2.units[j])
			if len(inters) > 0 {
				for _, mon := range inters {
					intersectedClusterUnitsMap[mon] = true
				}
			}
		}
	}
	intersection := make([]*dt.Monomer, 0)
	for key := range intersectedClusterUnitsMap {
		intersection = append(intersection, key)
	}
	return intersection
}

func IntersectClustersStickToDirection(cluster1, cluster2 *Cluster) []*dt.Monomer {
	if cluster1.mainDirection != cluster2.mainDirection &&
		cluster1.mainDirection != dt.GetReversedSide(cluster2.mainDirection) &&
		cluster1.axis != cluster2.axis {
		return nil
	}

	return IntersectClustersSoft(cluster1, cluster2)
}

type ClusterView struct {
	avg      float64
	axis     base.Axis
	clusters []*Cluster
}

func NewClusterView(globula *GlobulaView, avg float64, axis base.Axis) *ClusterView {
	newClusterView := new(ClusterView)
	newClusterView.avg = avg
	newClusterView.axis = axis
	newClusterView.clusters = findClusters(globula, axis, avg)
	return newClusterView
}

func NewClusterViewRaw(clusters []*Cluster, axis base.Axis) *ClusterView {
	newClusterView := new(ClusterView)
	newClusterView.avg = 0.0
	newClusterView.axis = axis
	newClusterView.clusters = clusters
	return newClusterView
}

func (clusterView *ClusterView) Colorize(reset bool) {
	for _, cluster := range clusterView.clusters {
		if reset {
			cluster.SetTypeOfMonomers(base.C)
		} else {
			cluster.SetTypeOfMonomers(dt.GetAxisColor(clusterView.axis))
		}
	}
}

func (clusterView *ClusterView) Trunkate() {
	idxsToRemove := make([]int, 0)
	for i, cluster := range clusterView.clusters {
		if cluster.Size() == 0 {
			idxsToRemove = append(idxsToRemove, i)
		}
	}

	for i := len(idxsToRemove) - 1; i >= 0; i-- {
		clusterView.clusters = append(clusterView.clusters[:idxsToRemove[i]], clusterView.clusters[idxsToRemove[i]+1:]...)
	}
}

func ForEachCluster(view *ClusterView, pred func(*Cluster)) {
	for _, cluster := range view.clusters {
		pred(cluster)
	}
}

func getDirection(monomer *dt.Monomer) dt.Side {
	nextMonomer := monomer.NextMonomer
	if nextMonomer == nil {
		return dt.SideUndefined
	}

	return monomer.GetSideOfSibling(nextMonomer)
}

func doTraverse(mainDirection dt.Side, directions []dt.Side, currMon *dt.Monomer, potCluster *[]*dt.Monomer) {
	if len(directions) == 0 {
		return
	}

	nextMonomer, err := currMon.GetSibling(directions[0])
	if err != nil || nextMonomer.IsTypeOf(base.MendeleevTableElementUndefined) {
		*potCluster = nil
		return
	}

	nextInMainDirectionMonomer, err := nextMonomer.GetSibling(mainDirection)
	if err != nil {
		*potCluster = nil
		return
	}

	nextNextMon := nextMonomer.NextMonomer
	prevNextMon := nextMonomer.PrevMonomer
	if nextNextMon == nil && prevNextMon == nil {
		*potCluster = nil
		return
	}

	if nextNextMon == nextInMainDirectionMonomer || prevNextMon == nextInMainDirectionMonomer {
		*potCluster = append(*potCluster, nextMonomer)
		*potCluster = append(*potCluster, nextInMainDirectionMonomer)
		doTraverse(mainDirection, directions[1:], nextMonomer, potCluster)
	} else {
		*potCluster = nil
	}
}
func fillClusters(directions []dt.Side, mainDirection dt.Side, axis base.Axis, monomer *dt.Monomer, clusters *[]*Cluster) {
	potCluster := make([]*dt.Monomer, 2)
	potCluster[0] = monomer
	sibling, _ := monomer.GetSibling(mainDirection)
	potCluster[1] = sibling
	doTraverse(mainDirection, directions, monomer, &potCluster)
	if len(potCluster) == 8 {
		newCluster := NewCluster([]*ClusterUnit{NewClusterUnit(potCluster, mainDirection, axis)}, mainDirection, axis)
		//new_cluster.MakeFullyConnected()
		*clusters = append(*clusters, newCluster)
	}
}

func findClusters(currentGlobula *GlobulaView, axis base.Axis, avg float64) []*Cluster {
	clusters := make([]*Cluster, 0)
	for _, pol := range currentGlobula.polymers {
		ForEachMonomer(pol, func(monomer *dt.Monomer) bool {
			for i := 0; i < pol.Len(); i++ {
				monomer := pol.polymer.GetMonomerByIdx(i)
				mainDirection := getDirection(monomer)
				if mainDirection == dt.SideUndefined || monomer.IsTypeOf(base.MendeleevTableElementUndefined) {
					continue
				}

				oldSize := len(clusters)
				if axis == base.AxisX && (mainDirection == dt.SideForward || mainDirection == dt.SideBackward) {
					fillClusters([]dt.Side{dt.SideLeft, dt.SideDown, dt.SideRight}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideLeft, dt.SideUp, dt.SideRight}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideRight, dt.SideDown, dt.SideLeft}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideRight, dt.SideUp, dt.SideLeft}, mainDirection, axis, monomer, &clusters)
				} else if axis == base.AxisY && (mainDirection == dt.SideLeft || mainDirection == dt.SideRight) {
					fillClusters([]dt.Side{dt.SideBackward, dt.SideDown, dt.SideForward}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideBackward, dt.SideUp, dt.SideForward}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideForward, dt.SideDown, dt.SideBackward}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideForward, dt.SideUp, dt.SideBackward}, mainDirection, axis, monomer, &clusters)
				} else if axis == base.AxisZ && (mainDirection == dt.SideUp || mainDirection == dt.SideDown) {
					fillClusters([]dt.Side{dt.SideLeft, dt.SideForward, dt.SideRight}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideLeft, dt.SideBackward, dt.SideRight}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideRight, dt.SideForward, dt.SideLeft}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SideRight, dt.SideBackward, dt.SideLeft}, mainDirection, axis, monomer, &clusters)
				} else {
					continue
				}

				if oldSize != len(clusters) {
					clusters = gatherClusters(clusters, avg, axis, -1)
				}
			}
			return true
		})
	}
	// Although we seem to have all the clusters already joined
	// it is possible that some clusters can be joined.
	// Since here we don't have a lot of them, this step shouldn't be
	// too time-consuming
	return gatherClusters(clusters, avg, axis, -1)
}

// Now join the extracted clusters recursively
// that means that if we, e.g, joined clusters 1 and 2 into a 1-2 cluster
// and cluster 2 and 3 into a 2-3 cluster, therefore, the joined ones have a common set of monomers
// which are in the cluster 2 and can be joined as well into a 1-2-3 cluster
func joinJoinedClusters(currentNode int, d *map[int][]int) []int {
	joinedClustersForCurrentNode := make([]int, 0)
	// If nothing to join, come back
	if v, ok := (*d)[currentNode]; !ok || len(v) == 0 {
		return joinedClustersForCurrentNode
	}

	v := (*d)[currentNode]
	for _, value := range v {
		// Add the current cluster as the connection
		joinedClustersForCurrentNode = append(joinedClustersForCurrentNode, value)
		// Recursively find others ready to be joined
		received := joinJoinedClusters(value, d)
		for _, recValue := range received {
			if !base.Contains(joinedClustersForCurrentNode, recValue) {
				joinedClustersForCurrentNode = append(joinedClustersForCurrentNode, recValue)
			}
		}
	}

	// All clusters for this cluster are to become a part of another list,
	// therefore, clear the current to avoid issues
	(*d)[currentNode] = make([]int, 0)
	return joinedClustersForCurrentNode
}

func gatherClusters(clusters []*Cluster, avg float64, avgAxis base.Axis, startFrom int) []*Cluster {
	if avgAxis == base.AxisCount {
		return clusters
	}

	for {
		clustersNew := make([]*Cluster, 0)
		// Build a dict where for each i we map a list of js that can be joined with i
		d := make(map[int][]int)
		// If start_from is None, then we need to compare all the cluster between each other
		// to do the deep check of clusters possible to join
		if startFrom < 0 {
			for i := 0; i < len(clusters)-1; i++ {
				for j := i + 1; j < len(clusters); j++ {
					joinedCluster := JoinClusters(clusters[i], clusters[j])
					if joinedCluster != nil {
						if _, ok := d[i]; !ok {
							d[i] = make([]int, 0)
						}

						d[i] = append(d[i], j)
					}
				}
			}
		} else {
			// The algorithm of highlighing clusters determine that all the new clusters are added
			// into the end of the clusters list, therefore, it is assumed it isn't able to join clusters
			// before start_from. So we need to check only the new ones
			for i := startFrom; i < len(clusters)-1; i++ {
				for j := startFrom; j < len(clusters); j++ {
					joinedCluster := JoinClusters(clusters[i], clusters[j])
					if joinedCluster != nil {
						if _, ok := d[i]; !ok {
							d[i] = make([]int, 0)
						}

						d[i] = append(d[i], j)
					}
				}
			}
		}

		// If there's nothing to join, stop the algorithm
		if len(d) == 0 {
			break
		}

		// Join the clusters. In fact, find the numbers belonging to the same cluster
		for key := range d {
			ret := joinJoinedClusters(key, &d)
			d[key] = ret
		}

		// Do actual join
		joinedClustersIndexes := make(map[int]bool)
		for key := range d {
			joinedClustersIndexes[key] = true
			for _, val := range d[key] {
				joinedClustersIndexes[val] = true
			}

			if len(d[key]) == 0 {
				continue
			}

			currCluster := clusters[key]
			for _, value := range d[key] {
				c := JoinClusters(currCluster, clusters[value])
				if c == nil {
					continue
				}
				currCluster = c
			}
			clustersNew = append(clustersNew, currCluster)
		}

		// Those clusters not been touched must be moved as well
		indexesNotJoined := make([]int, 0)
		for i := 0; i < len(clusters); i++ {
			if _, ok := joinedClustersIndexes[i]; !ok {
				indexesNotJoined = append(indexesNotJoined, i)
			}
		}

		temp := make([]*Cluster, len(indexesNotJoined))
		for i := 0; i < len(indexesNotJoined); i++ {
			temp[i] = clusters[indexesNotJoined[i]]
		}
		temp = append(temp, clustersNew...)

		clusters = temp
	}

	filteredClusters := make([]*Cluster, 0)
	for _, cluster := range clusters {
		if cluster.GetAvgLengthByAxis() >= avg {
			filteredClusters = append(filteredClusters, cluster)
		}
	}
	// return those which satisfy the condition
	return filteredClusters
}
