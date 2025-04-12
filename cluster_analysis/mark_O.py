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

def find_clusters(current_globula: GlobulaView, axis: Axis, avg_tuple: tuple[int, Axis] = None):
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

            old_size = len(clusters)
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

            if old_size != len(clusters):
                clusters = gather_clusters(clusters, avg_tuple, old_size)

    # Although we seem to have all the clusters already joined
    # it is possible that some clusters can be joined.
    # Since here we don't have a lot of them, this step shouldn't be
    # too time-consuming
    return gather_clusters(clusters, avg_tuple, None)

def get_monomers_count_in_clusters(clusters: list[Cluster]):
    s = set[Monomer]()
    for cluster in clusters:
        s.update(set(cluster.monomers))

    return len(s)

def gather_clusters(clusters: list[Cluster], avg_tuple : tuple[int, Axis] = None, start_from: int = None):
    if avg_tuple is None:
        return clusters

    avg, avg_axis = avg_tuple
    while True:
        clusters_new = list[Cluster]()
        # Build a dict where for each i we map a list of js that can be joined with i
        d = dict[int, list[int]]()
        # If start_from is None, then we need to compare all the cluster between each other
        # to do the deep check of clusters possible to join
        if start_from is None:
            for i in range(len(clusters) - 1):
                for j in range(i + 1, len(clusters)):
                    joined_cluster = join_clusters(clusters[i], clusters[j])
                    if joined_cluster is not None:
                        if i not in d.keys():
                            d[i] = list[int]()
                            
                        d[i].append(j)
        # The algorithm of highlighing clusters determine that all the new clusters are added
        # into the end of the clusters list, therefore, it is assumed it isn't able to join clusters
        # before start_from. So we need to check only the new ones
        else:
            for i in range(start_from, len(clusters)):
                for j in range(start_from):
                    joined_cluster = join_clusters(clusters[i], clusters[j])
                    if joined_cluster is not None:
                        if i not in d.keys():
                            d[i] = list[int]()
                            
                        d[i].append(j)

        # If there's nothing to join, stop the algorithm
        if len(d) == 0:
            break

        # Now join the extracted clusters recursively
        # that means that if we, e.g, joined clusters 1 and 2 into a 1-2 cluster
        # and cluster 2 and 3 into a 2-3 cluster, therefore, the joined ones have a common set of monomers
        # which are in the cluster 2 and can be joined as well into a 1-2-3 cluster
        def join_joined_clusters(current_node: int):
            joined_clusters_for_current_node = list[int]()
            # If nothing to join, come back
            if current_node not in d.keys() or len(d[current_node]) == 0:
                return joined_clusters_for_current_node
            
            for value in d[current_node]:
                # Add the current cluster as the connection
                joined_clusters_for_current_node.append(value)
                # Recursively find others ready to be joined
                received = join_joined_clusters(value)
                for rec_value in received:
                    if rec_value not in joined_clusters_for_current_node:
                        joined_clusters_for_current_node.append(rec_value)
                
            # All clusters for this cluster are to become a part of another list,
            # therefore, clear the current to avoid issues
            d[current_node].clear()
            return joined_clusters_for_current_node
        
        # Join the clusters. In fact, find the numbers belonging to the same cluster
        for key in d.keys():
            d[key] = join_joined_clusters(key) 
        
        # Do actual join
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

        # Those clusters not been touched must be moved as well
        indexes_not_joined = set([i for i in range(len(clusters)) if i not in joined_clusters_indexes])
        for i in indexes_not_joined:
            clusters_new.append(clusters[i])
        
        clusters = clusters_new

    # return those which satisfy the condition
    return list(set([c for c in clusters if c.get_avg_length_by_axis(avg_axis) >= avg]))

def mark_clusters(current_globula: GlobulaView, axis: Axis, avg_tuple: tuple[int, Axis] = None):
    clusters = find_clusters(current_globula, axis, avg_tuple)
            
    axis_type_dict = {
        Axis.X_AXIS: MonomerType.Owise,
        Axis.Y_AXIS: MonomerType.Nwise,
        Axis.Z_AXIS: MonomerType.Fwise
    }
    for cluster in clusters:
        cluster.set_type_of_monomers(axis_type_dict[axis])

def mark_common_clusters(current_globula: GlobulaView):
    clusters_X = set(find_clusters(current_globula, Axis.X_AXIS, (0, Axis.X_AXIS)))
    clusters_Y = set(find_clusters(current_globula, Axis.Y_AXIS, (0, Axis.Y_AXIS)))
    clusters_Z = set(find_clusters(current_globula, Axis.Z_AXIS, (0, Axis.Z_AXIS)))

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

def mark_biggest_clusters(current_globula: GlobulaView):
    clusters_X = find_clusters(current_globula, Axis.X_AXIS, (0, Axis.X_AXIS))
    clusters_Y = find_clusters(current_globula, Axis.Y_AXIS, (0, Axis.Y_AXIS))
    clusters_Z = find_clusters(current_globula, Axis.Z_AXIS, (0, Axis.Z_AXIS))
    total_monomers_count = sum([pol.len() for pol in current_globula])

    clusters_X = sorted(clusters_X, key=lambda c1: c1.size, reverse=True)
    clusters_Y = sorted(clusters_Y, key=lambda c1: c1.size, reverse=True)
    clusters_Z = sorted(clusters_Z, key=lambda c1: c1.size, reverse=True)

    chosen_clusters: list[Cluster] = []
    
    def clusters_monomers_count():
        return sum([c.size for c in chosen_clusters]) if len(chosen_clusters) > 0 else 0

    def test(testing, first, second):
        if len(testing) == 0:
            return False
        
        if len(first) != 0 and len(second) != 0:
            return testing[0].size >= first[0].size and testing[0].size >= second[0].size
        elif len(first) != 0:
            return testing[0].size >= first[0].size
        elif len(second) != 0:
            return testing[0].size >= second[0].size
        else:
            return True
    
    prev_crystallization_ratio = clusters_monomers_count() / total_monomers_count
    while True:
        percentage = prev_crystallization_ratio * 100
        if percentage > 55 or sum([len(clusters_X), len(clusters_Y), len(clusters_Z)]) == 0:
            break

        if 20 <= percentage:
            for cluster in chosen_clusters:
                if cluster.axis == Axis.X_AXIS:
                    cluster.set_type_of_monomers(MonomerType.Owise)
                elif cluster.axis == Axis.Y_AXIS:
                    cluster.set_type_of_monomers(MonomerType.Nwise)
                else:
                    cluster.set_type_of_monomers(MonomerType.Fwise)
            yield (True, percentage)
            for cluster in chosen_clusters:
                cluster.set_type_of_monomers(MonomerType.Usual)
            
        if test(clusters_X, clusters_Y, clusters_Z):
            chosen_clusters.append(clusters_X[0])
            clusters_X = clusters_X[1:]
        elif test(clusters_Y, clusters_X, clusters_Z):
            chosen_clusters.append(clusters_Y[0])
            clusters_Y = clusters_Y[1:]
        elif test(clusters_Z, clusters_X, clusters_Y):
            chosen_clusters.append(clusters_Z[0])
            clusters_Z = clusters_Z[1:]
        else:
            assert False, "It couldn't happen"

        prev_crystallization_ratio = clusters_monomers_count() / total_monomers_count


    for clusters in [clusters_X, clusters_Y, clusters_Z]:
        for cluster in clusters:
            cluster.set_type_of_monomers(MonomerType.Usual)
    return (False, 0)