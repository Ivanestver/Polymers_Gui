package views

import (
	"polymers/base"
	dt "polymers/datatypes"
)

const c_CLUSTER_UNIT_SIZE int = 8

type ClusterUnit struct {
	monomers      [c_CLUSTER_UNIT_SIZE]*dt.Monomer
	mainDirection dt.Side
	axis          dt.Axis
}

func NewClusterUnit(monomers []*dt.Monomer, mainDirection dt.Side, axis dt.Axis) *ClusterUnit {
	if len(monomers) != c_CLUSTER_UNIT_SIZE {
		return nil
	}
	newClusterUnit := new(ClusterUnit)
	for i := 0; i < c_CLUSTER_UNIT_SIZE; i++ {
		newClusterUnit.monomers[i] = monomers[i]
	}
	newClusterUnit.mainDirection = mainDirection
	newClusterUnit.axis = axis
	return newClusterUnit
}

func (this *ClusterUnit) Size() int {
	return c_CLUSTER_UNIT_SIZE
}

func (this *ClusterUnit) MakeFullyConnected() {
	for i := 0; i < this.Size()-1; i++ {
		currMon := this.monomers[i]
		for j := i + 1; j < this.Size(); j++ {
			sideMon := this.monomers[j]
			side := dt.GetSideByMonomers(currMon, sideMon)
			if side != dt.SIDE_Undefined {
				dt.MakeConnection(currMon, sideMon, dt.GetConnectionType(currMon, sideMon))
			}
		}
	}
}

func (this *ClusterUnit) contains(monomer *dt.Monomer) bool {
	for _, mon := range this.monomers {
		if dt.MonomersAreEqual(mon, monomer) {
			return true
		}
	}
	return false
}

func (this *ClusterUnit) setTypeOfMonomers(monomerType dt.MonomerType) {
	for i := 0; i < this.Size(); i++ {
		currMon := this.monomers[i]
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
	axis          dt.Axis
}

func NewCluster(unit []*ClusterUnit, mainDirection dt.Side, axis dt.Axis) *Cluster {
	newCluster := new(Cluster)
	newCluster.units = unit
	newCluster.mainDirection = mainDirection
	newCluster.axis = axis
	return newCluster
}

func (this *Cluster) Size() int {
	return len(this.units)
}

func (this *Cluster) MainDirection() dt.Side {
	return this.mainDirection
}

func (this *Cluster) Axis() dt.Axis {
	return this.axis
}

func (this *Cluster) SetTypeOfMonomers(monomerType dt.MonomerType) {
	for _, unit := range this.units {
		unit.setTypeOfMonomers(monomerType)
	}
}

func (this *Cluster) RemoveMonomer(monomer *dt.Monomer) bool {
	for i, unit := range this.units {
		if unit.contains(monomer) {
			this.units = append(this.units[:i], this.units[i+1:]...)
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

func (this *Cluster) GetAvgLengthByAxis() float64 {
	allMonomers := make([]*dt.Monomer, 0)
	for _, unit := range this.units {
		for _, mon := range unit.monomers {
			allMonomers = append(allMonomers, mon)
		}
	}
	sideBackward := dt.GetReversedSide(this.mainDirection)

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
		sibling, err = currMonomer.GetSibling(this.mainDirection)
		for err == nil && containsMonomer(&allMonomers, sibling) {
			currMonomer = sibling
			usedMonomers = append(usedMonomers, currMonomer)
			length += 1.0
			sibling, err = currMonomer.GetSibling(this.mainDirection)
		}
		lengths = append(lengths, length)
	}

	return base.Sum(lengths) / float64(len(lengths))
}

func (this *Cluster) MakeFullyConnected() {
	for _, unit := range this.units {
		unit.MakeFullyConnected()
	}
}

func JoinClusters(cluster1, cluster2 *Cluster) *Cluster {
	if cluster1.mainDirection != cluster2.mainDirection &&
		cluster1.mainDirection != dt.GetReversedSide(cluster2.mainDirection) &&
		cluster1.axis != cluster2.axis {
		return nil
	}

	intersection := IntersectClusters_StickToDirection(cluster1, cluster2)

	// If there's nothing to connect, these are probably different clusters
	if len(intersection) < 2 {
		return nil
	}

	units := make([]*ClusterUnit, 0)
	for _, unit1 := range cluster1.units {
		units = append(units, unit1)
	}

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

func IntersectClusters_Soft(cluster1, cluster2 *Cluster) []*dt.Monomer {
	intersectedClusterUnits_map := make(map[*dt.Monomer]bool, 0)
	for i := 0; i < cluster1.Size(); i++ {
		for j := 0; j < cluster2.Size(); j++ {
			inters := intersectClusterUnits(cluster1.units[i], cluster2.units[j])
			if len(inters) > 0 {
				for _, mon := range inters {
					intersectedClusterUnits_map[mon] = true
				}
			}
		}
	}
	intersection := make([]*dt.Monomer, 0)
	for key := range intersectedClusterUnits_map {
		intersection = append(intersection, key)
	}
	return intersection
}

func IntersectClusters_StickToDirection(cluster1, cluster2 *Cluster) []*dt.Monomer {
	if cluster1.mainDirection != cluster2.mainDirection &&
		cluster1.mainDirection != dt.GetReversedSide(cluster2.mainDirection) &&
		cluster1.axis != cluster2.axis {
		return nil
	}

	return IntersectClusters_Soft(cluster1, cluster2)
}

type ClusterView struct {
	avg      float64
	axis     dt.Axis
	clusters []*Cluster
}

func NewClusterView(globula *GlobulaView, avg float64, axis dt.Axis) *ClusterView {
	newClusterView := new(ClusterView)
	newClusterView.avg = avg
	newClusterView.axis = axis
	newClusterView.clusters = findClusters(globula, axis, avg)
	return newClusterView
}

func (this *ClusterView) Colorize(reset bool) {
	for _, cluster := range this.clusters {
		if reset {
			cluster.SetTypeOfMonomers(dt.MONOMER_TYPE_USUAL)
		} else {
			cluster.SetTypeOfMonomers(dt.GetAxisColor(this.axis))
		}
	}
}

func (this *ClusterView) Trunkate() {
	idxsToRemove := make([]int, 0)
	for i, cluster := range this.clusters {
		if cluster.Size() == 0 {
			idxsToRemove = append(idxsToRemove, i)
		}
	}

	for i := len(idxsToRemove) - 1; i >= 0; i-- {
		this.clusters = append(this.clusters[:idxsToRemove[i]], this.clusters[idxsToRemove[i]+1:]...)
	}
}

func ForEachCluster(view *ClusterView, pred func(*Cluster)) {
	for _, cluster := range view.clusters {
		pred(cluster)
	}
}

func get_direction(monomer *dt.Monomer) dt.Side {
	nextMonomer := monomer.NextMonomer
	if nextMonomer == nil {
		return dt.SIDE_Undefined
	}

	return monomer.GetSideOfSibling(nextMonomer)
}

func doTraverse(mainDirection dt.Side, directions []dt.Side, currMon *dt.Monomer, potCluster *[]*dt.Monomer) {
	if len(directions) == 0 {
		return
	}

	nextMonomer, err := currMon.GetSibling(directions[0])
	if err != nil || nextMonomer.IsTypeOf(dt.MONOMER_TYPE_UNDEFINED) {
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
func fillClusters(directions []dt.Side, mainDirection dt.Side, axis dt.Axis, monomer *dt.Monomer, clusters *[]*Cluster) {
	potCluster := make([]*dt.Monomer, 2)
	potCluster[0] = monomer
	sibling, _ := monomer.GetSibling(mainDirection)
	potCluster[1] = sibling
	doTraverse(mainDirection, directions, monomer, &potCluster)
	if len(potCluster) == 8 {
		new_cluster := NewCluster([]*ClusterUnit{NewClusterUnit(potCluster, mainDirection, axis)}, mainDirection, axis)
		//new_cluster.MakeFullyConnected()
		*clusters = append(*clusters, new_cluster)
	}
}

func findClusters(currentGlobula *GlobulaView, axis dt.Axis, avg float64) []*Cluster {
	clusters := make([]*Cluster, 0)
	for _, pol := range currentGlobula.polymers {
		ForEachMonomer(pol, func(monomer *dt.Monomer) {
			for i := 0; i < pol.Len(); i++ {
				monomer := pol.polymer.GetMonomerByIdx(i)
				mainDirection := get_direction(monomer)
				if mainDirection == dt.SIDE_Undefined || monomer.IsTypeOf(dt.MONOMER_TYPE_UNDEFINED) {
					continue
				}

				old_size := len(clusters)
				if axis == dt.X_AXIS && (mainDirection == dt.SIDE_Forward || mainDirection == dt.SIDE_Backward) {
					fillClusters([]dt.Side{dt.SIDE_Left, dt.SIDE_Down, dt.SIDE_Right}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Left, dt.SIDE_Up, dt.SIDE_Right}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Right, dt.SIDE_Down, dt.SIDE_Left}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Right, dt.SIDE_Up, dt.SIDE_Left}, mainDirection, axis, monomer, &clusters)
				} else if axis == dt.Y_AXIS && (mainDirection == dt.SIDE_Left || mainDirection == dt.SIDE_Right) {
					fillClusters([]dt.Side{dt.SIDE_Backward, dt.SIDE_Down, dt.SIDE_Forward}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Backward, dt.SIDE_Up, dt.SIDE_Forward}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Forward, dt.SIDE_Down, dt.SIDE_Backward}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Forward, dt.SIDE_Up, dt.SIDE_Backward}, mainDirection, axis, monomer, &clusters)
				} else if axis == dt.Z_AXIS && (mainDirection == dt.SIDE_Up || mainDirection == dt.SIDE_Down) {
					fillClusters([]dt.Side{dt.SIDE_Left, dt.SIDE_Forward, dt.SIDE_Right}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Left, dt.SIDE_Backward, dt.SIDE_Right}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Right, dt.SIDE_Forward, dt.SIDE_Left}, mainDirection, axis, monomer, &clusters)
					fillClusters([]dt.Side{dt.SIDE_Right, dt.SIDE_Backward, dt.SIDE_Left}, mainDirection, axis, monomer, &clusters)
				} else {
					continue
				}

				if old_size != len(clusters) {
					clusters = gather_clusters(clusters, avg, axis, -1)
				}
			}
		})
	}
	// Although we seem to have all the clusters already joined
	// it is possible that some clusters can be joined.
	// Since here we don't have a lot of them, this step shouldn't be
	// too time-consuming
	return gather_clusters(clusters, avg, axis, -1)
}

func get_monomers_count_in_clusters(clusters []*Cluster) int {
	s := make(map[*ClusterUnit]bool)
	for _, cluster := range clusters {
		for _, mon := range cluster.units {
			s[mon] = true
		}
	}

	return len(s)
}

// Now join the extracted clusters recursively
// that means that if we, e.g, joined clusters 1 and 2 into a 1-2 cluster
// and cluster 2 and 3 into a 2-3 cluster, therefore, the joined ones have a common set of monomers
// which are in the cluster 2 and can be joined as well into a 1-2-3 cluster
func join_joined_clusters(current_node int, d *map[int][]int) []int {
	joined_clusters_for_current_node := make([]int, 0)
	// If nothing to join, come back
	if v, ok := (*d)[current_node]; !ok || len(v) == 0 {
		return joined_clusters_for_current_node
	}

	v, _ := (*d)[current_node]
	for _, value := range v {
		// Add the current cluster as the connection
		joined_clusters_for_current_node = append(joined_clusters_for_current_node, value)
		// Recursively find others ready to be joined
		received := join_joined_clusters(value, d)
		for _, rec_value := range received {
			if !base.Contains(joined_clusters_for_current_node, rec_value) {
				joined_clusters_for_current_node = append(joined_clusters_for_current_node, rec_value)
			}
		}
	}

	// All clusters for this cluster are to become a part of another list,
	// therefore, clear the current to avoid issues
	(*d)[current_node] = make([]int, 0)
	return joined_clusters_for_current_node
}

func gather_clusters(clusters []*Cluster, avg float64, avg_axis dt.Axis, start_from int) []*Cluster {
	if avg_axis == dt.AXIS_COUNT {
		return clusters
	}

	for true {
		clusters_new := make([]*Cluster, 0)
		// Build a dict where for each i we map a list of js that can be joined with i
		d := make(map[int][]int)
		// If start_from is None, then we need to compare all the cluster between each other
		// to do the deep check of clusters possible to join
		if start_from < 0 {
			for i := 0; i < len(clusters)-1; i++ {
				for j := i + 1; j < len(clusters); j++ {
					joined_cluster := JoinClusters(clusters[i], clusters[j])
					if joined_cluster != nil {
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
			for i := start_from; i < len(clusters)-1; i++ {
				for j := start_from; j < len(clusters); j++ {
					joined_cluster := JoinClusters(clusters[i], clusters[j])
					if joined_cluster != nil {
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
			ret := join_joined_clusters(key, &d)
			d[key] = ret
		}

		// Do actual join
		joined_clusters_indexes := make(map[int]bool)
		for key := range d {
			joined_clusters_indexes[key] = true
			for _, val := range d[key] {
				joined_clusters_indexes[val] = true
			}

			if len(d[key]) == 0 {
				continue
			}

			curr_cluster := clusters[key]
			for _, value := range d[key] {
				c := JoinClusters(curr_cluster, clusters[value])
				if c == nil {
					continue
				}
				curr_cluster = c
			}
			clusters_new = append(clusters_new, curr_cluster)
		}

		// Those clusters not been touched must be moved as well
		indexes_not_joined := make([]int, 0)
		for i := 0; i < len(clusters); i++ {
			if _, ok := joined_clusters_indexes[i]; !ok {
				indexes_not_joined = append(indexes_not_joined, i)
			}
		}

		temp := make([]*Cluster, len(indexes_not_joined))
		for i := 0; i < len(indexes_not_joined); i++ {
			temp[i] = clusters[indexes_not_joined[i]]
		}
		for _, cluster := range clusters_new {
			temp = append(temp, cluster)
		}

		clusters = temp
	}

	filtered_clusters := make([]*Cluster, 0)
	for _, cluster := range clusters {
		if cluster.GetAvgLengthByAxis() >= avg {
			filtered_clusters = append(filtered_clusters, cluster)
		}
	}
	// return those which satisfy the condition
	return filtered_clusters
}
