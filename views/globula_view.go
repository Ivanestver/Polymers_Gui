package views

import (
	"encoding/json"
	"polymers/datatypes"
	dt "polymers/datatypes"
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

func (this *GlobulaView) Name() string {
	return this.name
}

func (this *GlobulaView) Len() int {
	return len(this.polymers)
}

func (this *GlobulaView) Reset() {
	for _, pol := range this.polymers {
		ForEachMonomer(pol, func(mon *datatypes.Monomer) { mon.MonomerType = datatypes.MONOMER_TYPE_USUAL })
	}
}

// func (this *GlobulaView) ToJson() ([]byte, error) {
// 	pols_countMap := make(map[string]interface{})
// 	pols_countMap["value"] = len(this.polymers)
// 	pols_countMap["name"] = "Количество полимеров"

// 	pols_map := make(map[string]interface{})
// 	pols_map["value"] = make([]string, len(this.polymers))

// 	var jsonDict map[string]interface{}
// 	jsonDict["polymers_count"] = pols_countMap
// 	return nil, nil
// }

func ForEachPolymer(globula *GlobulaView, pred func(*PolymerView)) {
	for _, pol := range globula.polymers {
		pred(pol)
	}
}

func (this *GlobulaView) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name     string
		Polymers []*PolymerView
	}{
		Name:     this.name,
		Polymers: this.polymers,
	})
}

func (this *GlobulaView) XClusters(avg float64) *ClusterView {
	if this.xClusters == nil {
		this.xClusters = NewClusterView(this, avg, dt.X_AXIS)
	}
	return this.xClusters
}

func (this *GlobulaView) YClusters(avg float64) *ClusterView {
	if this.yClusters == nil {
		this.yClusters = NewClusterView(this, avg, dt.Y_AXIS)
	}
	return this.yClusters
}

func (this *GlobulaView) ZClusters(avg float64) *ClusterView {
	if this.zClusters == nil {
		this.zClusters = NewClusterView(this, avg, dt.Z_AXIS)
	}
	return this.zClusters
}

func (this *GlobulaView) CommonClusters() (*ClusterView, *ClusterView, *ClusterView) {
	if !this.commonClusterDone {
		clusters_X := this.XClusters(0.0)
		clusters_Y := this.YClusters(0.0)
		clusters_Z := this.ZClusters(0.0)

		ForEachCluster(clusters_X, func(cluster_X *Cluster) {
			ForEachCluster(clusters_Y, func(cluster_Y *Cluster) {
				common_monomers := IntersectClusters(cluster_X, cluster_Y)
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
				common_monomers := IntersectClusters(cluster_Y, cluster_Z)
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
				common_monomers := IntersectClusters(cluster_X, cluster_Z)
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
		this.commonClusterDone = true
	}
	return this.xClusters, this.yClusters, this.zClusters
}
