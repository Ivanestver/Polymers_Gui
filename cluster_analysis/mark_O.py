from alg.polymer_lib import Monomer, MonomerType
from polymer_view import GlobulaView
from alg.config import Side, Axis
from cluster_analysis.cluster import Cluster, join_clusters, intersect_clusters

def get_direction(monomer: Monomer):
    next_monomer = monomer.next_monomer
    if next_monomer is None:
        return Side.Undefined
    
    for side, side_monomer in monomer.sides.items():
        if side_monomer is None:
            continue
        if side_monomer == next_monomer:
            return side

    return Side.Undefined

def do_traverse(main_direction: Side, directions: list[Side], curr_monomer: Monomer, pot_cluster: list[Monomer]):
    if len(directions) == 0:
        return
    
    next_monomer: Monomer = curr_monomer.sides[directions[0]]
    if next_monomer is None or next_monomer.is_of_type(MonomerType.Undefined):
        pot_cluster.clear()
        return
    
    next_in_main_direction_monomer: Monomer = next_monomer.sides[main_direction]
    if next_in_main_direction_monomer is None:
        pot_cluster.clear()
        return

    next_next_monomer = next_monomer.next_monomer
    prev_next_monomer = next_monomer.prev_monomer
    if next_next_monomer is None and prev_next_monomer is None:
        pot_cluster.clear()
        return

    if next_next_monomer == next_in_main_direction_monomer or prev_next_monomer == next_in_main_direction_monomer:
        pot_cluster.append(next_monomer)
        pot_cluster.append(next_in_main_direction_monomer)
        do_traverse(main_direction, directions[1:], next_monomer, pot_cluster)
    else:
        pot_cluster.clear()

def find_clusters(current_globula: GlobulaView, axis: Axis):
    clusters = list[Cluster]()
    def fill_clusters(directions, main_direction):
        pot_cluster: list[Monomer] = [monomer, monomer.sides[main_direction]]
        do_traverse(main_direction, directions, monomer, pot_cluster)
        if len(pot_cluster) == 8:
            clusters.append(Cluster(pot_cluster, main_direction, axis))
    for pol in current_globula:
        for monomer in pol:
            main_direction = get_direction(monomer)
            if main_direction == Side.Undefined:
                continue

            if axis == Axis.X_AXIS and main_direction in [Side.Forward, Side.Backward]:
                fill_clusters([Side.Left, Side.Down, Side.Right], main_direction)
                fill_clusters([Side.Left, Side.Up, Side.Right], main_direction)
                fill_clusters([Side.Right, Side.Down, Side.Left], main_direction)
                fill_clusters([Side.Right, Side.Up, Side.Left], main_direction)
            elif axis == Axis.Y_AXIS and main_direction in [Side.Left, Side.Right]:
                fill_clusters([Side.Backward, Side.Down, Side.Forward], main_direction)
                fill_clusters([Side.Backward, Side.Up, Side.Forward], main_direction)
                fill_clusters([Side.Forward, Side.Down, Side.Backward], main_direction)
                fill_clusters([Side.Forward, Side.Up, Side.Backward], main_direction)
            elif axis == Axis.Z_AXIS and main_direction in [Side.Up, Side.Down]:
                fill_clusters([Side.Left, Side.Forward, Side.Right], main_direction)
                fill_clusters([Side.Left, Side.Backward, Side.Right], main_direction)
                fill_clusters([Side.Right, Side.Forward, Side.Left], main_direction)
                fill_clusters([Side.Right, Side.Backward, Side.Left], main_direction)
            else:
                continue

    return clusters

def get_monomers_count_in_clusters(clusters: list[Cluster]):
    s = set[Monomer]()
    for cluster in clusters:
        s.update(set(cluster.monomers))

    return len(s)

def gather_clusters(clusters: list[Cluster], avg_tuple : tuple[int, Axis] = None):
    if avg_tuple is not None:
        avg, avg_axis = avg_tuple
        while True:
            clusters_new = list[Cluster]()
            d = dict[int, list[int]]()
            for i in range(len(clusters) - 1):
                for j in range(i + 1, len(clusters)):
                    joined_cluster = join_clusters(clusters[i], clusters[j])
                    if joined_cluster is not None:
                        if i not in d.keys():
                            d[i] = list[int]()
                            
                        d[i].append(j)

            if len(d) == 0:
                break

            def join_joined_clusters(current_node: int):
                joined_clusters_for_current_node = list[int]()
                if current_node not in d.keys() or len(d[current_node]) == 0:
                    return joined_clusters_for_current_node
                
                for value in d[current_node]:
                    joined_clusters_for_current_node.append(value)
                    received = join_joined_clusters(value)
                    for rec_value in received:
                        if rec_value not in joined_clusters_for_current_node:
                            joined_clusters_for_current_node.append(rec_value)
                    
                d[current_node].clear()
                return joined_clusters_for_current_node
            
            for key in d.keys():
                d[key] = join_joined_clusters(key) 
            
            joined_clusters_indexes = set()
            for key in d.keys():
                joined_clusters_indexes.add(key)
                joined_clusters_indexes.update(d[key])

                if len(d[key]) == 0:
                    continue

                curr_cluster = clusters[key]
                for value in d[key]:
                    c = join_clusters(curr_cluster, clusters[value])
                    if c is None:
                        continue
                    curr_cluster = c
                clusters_new.append(curr_cluster)

            indexes_not_joined = set([i for i in range(len(clusters)) if i not in joined_clusters_indexes])
            for i in indexes_not_joined:
                clusters_new.append(clusters[i])
            
            clusters = clusters_new

        clusters = set([c for c in clusters if c.get_avg_length_by_axis(avg_axis) >= avg])
    return clusters

def mark_clusters(current_globula: GlobulaView, axis: Axis, avg_tuple: tuple[int, Axis] = None):
    clusters = find_clusters(current_globula, axis)
    clusters = join_clusters(clusters, avg_tuple)
            
    axis_type_dict = {
        Axis.X_AXIS: MonomerType.Owise,
        Axis.Y_AXIS: MonomerType.Nwise,
        Axis.Z_AXIS: MonomerType.Fwise
    }
    for cluster in clusters:
        cluster.set_type_of_monomers(axis_type_dict[axis])

def mark_common_clusters(current_globula: GlobulaView):
    clusters_X = gather_clusters(find_clusters(current_globula, Axis.X_AXIS), (0, Axis.X_AXIS))
    clusters_Y = gather_clusters(find_clusters(current_globula, Axis.Y_AXIS), (0, Axis.Y_AXIS))
    clusters_Z = gather_clusters(find_clusters(current_globula, Axis.Z_AXIS), (0, Axis.Z_AXIS))

    clusters_X_Y = list([intersect_clusters(cluster_X, cluster_Y) for cluster_X in clusters_X for cluster_Y in clusters_Y])
    clusters_X_Y = list([mon for cluster in clusters_X_Y for mon in cluster])

    clusters_X_Z = list([intersect_clusters(cluster_X, cluster_Z) for cluster_X in clusters_X for cluster_Z in clusters_Z])
    clusters_X_Z = list([mon for cluster in clusters_X_Z for mon in cluster])

    clusters_Y_Z = list([intersect_clusters(cluster_Y, cluster_Z) for cluster_Y in clusters_Y for cluster_Z in clusters_Z])
    clusters_Y_Z = list([mon for cluster in clusters_Y_Z for mon in cluster])

    for mon in clusters_X:
        mon.set_type_of_monomers(MonomerType.Owise)
        
    for mon in clusters_Y:
        mon.set_type_of_monomers(MonomerType.Nwise)
        
    for mon in clusters_Z:
        mon.set_type_of_monomers(MonomerType.Fwise)

    for mon in clusters_X_Y:
        mon.type = MonomerType.Usual
        for cluster in clusters_X:
            if cluster.remove_monomer(mon):
                break
        for cluster in clusters_Y:
            if cluster.remove_monomer(mon):
                break

    for mon in clusters_X_Z:
        mon.type = MonomerType.Usual
        for cluster in clusters_X:
            if cluster.remove_monomer(mon):
                break
        for cluster in clusters_Z:
            if cluster.remove_monomer(mon):
                break

    for mon in clusters_Y_Z:
        mon.type = MonomerType.Usual
        for cluster in clusters_Y:
            if cluster.remove_monomer(mon):
                break
        for cluster in clusters_Z:
            if cluster.remove_monomer(mon):
                break