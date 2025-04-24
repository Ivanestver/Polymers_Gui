from alg.monomer_lib import Monomer, MonomerType, make_connection, break_connection, get_connection_type, get_side_by_monomers
from alg.config import Side, Axis, get_reversed_side, get_additional_sides

class Cluster:
    def __init__(self, monomers: list[Monomer], main_direction: Side, axis: Axis):
        self._monomers = monomers
        self._main_direction = main_direction
        self._axis = axis

    def make_fully_connected(self):
        for i in range(len(self._monomers)):
            for j in range(i + 1, len(self._monomers)):
                if i == j:
                    continue
                side = get_side_by_monomers(self._monomers[i], self._monomers[j])
                make_connection( \
                    self._monomers[i], \
                    self._monomers[j],\
                    side,\
                    get_connection_type(self._monomers[i], self._monomers[j], side)
                )

    @property
    def monomers(self):
        return self._monomers
    
    @property
    def main_direction(self):
        return self._main_direction

    @property
    def axis(self):
        return self._axis

    @property
    def size(self):
        return len(self._monomers)

    def set_type_of_monomers(self, type: MonomerType):
        for mon in self.monomers:
            mon.type = type
            additional_sides = get_additional_sides()
            if type in [MonomerType.Undefined, MonomerType.Usual]:
                for side in additional_sides:
                    side_mon: Monomer = mon.get_sibling(side)
                    if side_mon is not None:
                        break_connection(mon, side_mon, side)

    def remove_monomer(self, monomer: Monomer):
        try:
            idx = self._monomers.index(monomer)
            if idx >= 0:
                additional_sides = get_additional_sides()
                if monomer.type in [MonomerType.Undefined, MonomerType.Usual]:
                    for side in additional_sides:
                        curr_mon: Monomer = self._monomers[idx]
                        side_mon: Monomer = curr_mon.get_sibling(side)
                        if side_mon is not None:
                            break_connection(curr_mon, side_mon, side)
                del self._monomers[idx]
                return True
        except ValueError as er:
            return False

    def get_avg_length_by_axis(self, required_axis: Axis):
        side_backward = get_reversed_side(self.main_direction)

        lengths = []
        used_monomers = list[Monomer]()
        for mon in self.monomers:
            if mon in used_monomers:
                continue
            
            curr_monomer = mon
            while curr_monomer.get_sibling(side_backward) is not None and curr_monomer.get_sibling(side_backward) in self.monomers:
                curr_monomer = curr_monomer.get_sibling(side_backward)

            used_monomers.append(curr_monomer)
            length = 1
            while curr_monomer.get_sibling(self.main_direction) is not None and curr_monomer.get_sibling(self.main_direction) in self.monomers:
                curr_monomer = curr_monomer.get_sibling(self.main_direction)
                used_monomers.append(curr_monomer)
                length += 1

            lengths.append(length)
    
        return sum(lengths) / len(lengths)

def intersect_clusters(cluster1: Cluster, cluster2: Cluster):
    return set(cluster1.monomers).intersection(set(cluster2.monomers))

def join_clusters(cluster1: Cluster, cluster2: Cluster) -> Cluster:
    if (cluster1.main_direction != cluster2.main_direction and \
        cluster1.main_direction != get_reversed_side(cluster2.main_direction)) or \
        cluster1.axis != cluster2.axis:
        return None

    intersection = set(cluster1.monomers).intersection(set(cluster2.monomers))
    len_inter = len(intersection)
    if len_inter < 2:
        return None
    
    s = set[Monomer]()
    s.update(set(cluster1.monomers))
    s.update(set(cluster2.monomers))

    return Cluster(list(s), cluster1.main_direction, cluster1.axis)