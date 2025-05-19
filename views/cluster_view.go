package views

import (
	"polymers/base"
	dt "polymers/datatypes"
)

type Cluster struct {
	monomers      []*dt.Monomer
	mainDirection dt.Side
	axis          dt.Axis
}

func NewCluster(monomers []*dt.Monomer, mainDirection dt.Side, axis dt.Axis) *Cluster {
	newCluster := new(Cluster)
	newCluster.monomers = monomers
	newCluster.mainDirection = mainDirection
	newCluster.axis = axis
	return newCluster
}

func (this *Cluster) Monomers() []*dt.Monomer {
	return this.monomers
}

func (this *Cluster) MainDirection() dt.Side {
	return this.mainDirection
}

func (this *Cluster) Axis() dt.Axis {
	return this.axis
}

func (this *Cluster) Size() int {
	return len(this.monomers)
}

func (this *Cluster) SetTypeOfMonomers(monomerType dt.MonomerType) {
	for _, mon := range this.monomers {
		mon.MonomerType = monomerType
		additionalSides := dt.GetAdditionalSides()
		if monomerType == dt.MONOMER_TYPE_UNDEFINED || monomerType == dt.MONOMER_TYPE_USUAL {
			for _, side := range additionalSides {
				sideMon, err := mon.GetSibling(side)
				if err == nil {
					dt.BreakConnection(mon, sideMon, side)
				}
			}
		}
	}
}

func (this *Cluster) RemoveMonomer(monomer *dt.Monomer) bool {
	idx := base.Index[*dt.Monomer](this.monomers, monomer)
	if idx < 0 {
		return false
	}

	additionalSides := dt.GetAdditionalSides()
	if monomer.MonomerType == dt.MONOMER_TYPE_UNDEFINED || monomer.MonomerType == dt.MONOMER_TYPE_USUAL {
		// Break all the connections between the current monomer and the others.
		// Meanwhile, gather those other monomers in order to break connections between them
		candidates := make([]*dt.Monomer, 0)
		currMon := this.monomers[idx]
		for _, side := range additionalSides {
			sideMon, err := currMon.GetSibling(side)
			// It's essential to remove only those monomers which belong to the cluster
			if err == nil && base.Contains(this.monomers, sideMon) {
				candidates = append(candidates, sideMon)
				dt.BreakConnection(currMon, sideMon, side)
			}
		}
		for _, side := range dt.GetMovementSides() {
			sideMon, err := currMon.GetSibling(side)
			if err == nil && base.Contains(this.monomers, sideMon) {
				candidates = append(candidates, sideMon)
				if !dt.MonomersAreEqual(sideMon, currMon.NextMonomer) &&
					!dt.MonomersAreEqual(sideMon, currMon.PrevMonomer) {
					dt.BreakConnection(currMon, sideMon, side)
				}
			}
		}
		this.monomers = append(this.monomers[:idx], this.monomers[idx+1:]...)

		// Unattached the current monomer. It's time to break all the connections
		// between candidates. We gotta consider all the possible pairs
		for i := 0; i < len(candidates); i++ {
			for j := i + 1; j < len(candidates); j++ {
				dt.BreakConnection(candidates[i], candidates[j], dt.GetSideByMonomers(candidates[i], candidates[j]))
			}
		}
		// Now we may have the situation when a candidate has no connection to the cluster's monomers
		// Therefore, such candidates must be removed as well

		// Firstly, don't consider those which don't belong to this cluster
		candidates_temp := make([]*dt.Monomer, 0)
		for _, candidate := range candidates {
			idx := base.Index_if(this.monomers, func(m *dt.Monomer) bool { return dt.MonomersAreEqual(m, candidate) })
			if idx >= 0 {
				candidates_temp = append(candidates_temp, candidate)
			}
		}
		candidates = candidates_temp

		// Then, the idea is that if a candidate doesn't have any additional connections
		// excepting potential connections as the polymer's part, it has to be removed
		candidates_temp = make([]*dt.Monomer, 0)
		for _, candidate := range candidates {
			// If there is any additional side connection, skip this monomer
			noConnection := true
			for _, side := range additionalSides {
				if candidate.GetTypeOfConnectionWithSide(side) != dt.CONNECTION_TYPE_UNDEFINED {
					noConnection = false
					break
				}
			}
			if !noConnection {
				continue
			}
			nextMonSide := candidate.GetSideOfSibling(candidate.NextMonomer)
			prevMonSide := candidate.GetSideOfSibling(candidate.PrevMonomer)
			noConnection = true
			for _, side := range dt.GetMovementSides() {
				// If this is the side of the polymer's chain, skip
				if nextMonSide != dt.SIDE_Undefined && side == nextMonSide {
					continue
				}
				if prevMonSide != dt.SIDE_Undefined && side == prevMonSide {
					continue
				}
				connType := candidate.GetTypeOfConnectionWithSide(side)
				if connType != dt.CONNECTION_TYPE_UNDEFINED {
					noConnection = false
					break
				}
			}
			if noConnection {
				candidates_temp = append(candidates_temp, candidate)
			}
		}
		for _, candidate := range candidates_temp {
			idx := base.Index_if(this.monomers, func(m *dt.Monomer) bool { return dt.MonomersAreEqual(m, candidate) })
			if idx >= 0 {
				this.monomers = append(this.monomers[:idx], this.monomers[idx+1:]...)
			}
		}
	}
	return true
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
	sideBackward := dt.GetReversedSide(this.mainDirection)

	lengths := make([]float64, 0)
	usedMonomers := make([]*dt.Monomer, 0)
	for _, mon := range this.monomers {
		if containsMonomer(&usedMonomers, mon) {
			continue
		}

		currMonomer := mon
		sibling, err := currMonomer.GetSibling(sideBackward)
		for err == nil && containsMonomer(&this.monomers, sibling) {
			currMonomer = sibling
			sibling, err = currMonomer.GetSibling(sideBackward)
		}

		usedMonomers = append(usedMonomers, currMonomer)
		length := 1.0
		sibling, err = currMonomer.GetSibling(this.mainDirection)
		for err == nil && containsMonomer(&this.monomers, sibling) {
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
	// for i, mon_i := range this.monomers {
	// 	for j := i + 1; j < len(this.monomers); j++ {
	// 		dt.MakeConnection(mon_i, this.monomers[j], dt.GetConnectionType(mon_i, this.monomers[j]))
	// 	}
	// }

	for i := 0; i < len(this.monomers)-1; i++ {
		currMon := this.monomers[i]
		for j := i + 1; j < len(this.monomers); j++ {
			sideMon := this.monomers[j]
			side := dt.GetSideByMonomers(currMon, sideMon)
			if side != dt.SIDE_Undefined {
				dt.MakeConnection(currMon, sideMon, dt.GetConnectionType(currMon, sideMon))
			}
		}
	}
}

func JoinClusters(cluster1, cluster2 *Cluster) *Cluster {
	if cluster1.mainDirection != cluster2.mainDirection &&
		cluster1.mainDirection != dt.GetReversedSide(cluster2.mainDirection) &&
		cluster1.axis != cluster2.axis {
		return nil
	}

	intersection := IntersectClusters(cluster1, cluster2)

	// If there's no or only one point to join, then they're possibly of different cluster
	if len(intersection) < 2 {
		return nil
	}
	monomers := make([]*dt.Monomer, 0)
	for _, mon := range cluster1.monomers {
		monomers = append(monomers, mon)
	}
	for _, mon := range cluster2.monomers {
		if !base.Any(intersection, func(m *dt.Monomer) bool { return dt.MonomersAreEqual(m, mon) }) {
			monomers = append(monomers, mon)
		}
	}

	return NewCluster(monomers, cluster1.mainDirection, cluster1.axis)
}

func IntersectClusters(cluster1, cluster2 *Cluster) []*dt.Monomer {
	m := make(map[*dt.Monomer]bool)
	for _, mon := range cluster1.monomers {
		m[mon] = true
	}
	intersection := make([]*dt.Monomer, 0)
	for _, mon := range cluster2.monomers {
		if _, ok := m[mon]; ok {
			intersection = append(intersection, mon)
		}
	}
	return intersection
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
		new_cluster := NewCluster(potCluster, mainDirection, axis)
		new_cluster.MakeFullyConnected()
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
	s := make(map[*dt.Monomer]bool)
	for _, cluster := range clusters {
		for _, mon := range cluster.monomers {
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
