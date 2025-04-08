from alg.field_lib import Monomer, MonomerType
from alg.config import Side, Axis, get_reversed_side, get_side, MoveDirection

class Cluster:
    def __init__(self, monomers: list[Monomer], main_direction: Side, axis: Axis):
        self._monomers = monomers
        self._main_direction = main_direction
        self._axis = axis

    @property
    def monomers(self):
        return self._monomers
    
    @property
    def main_direction(self):
        return self._main_direction

    @property
    def axis(self):
        return self._axis

    def set_type_of_monomers(self, type: MonomerType):
        for mon in self.monomers:
            mon.type = type

    def remove_monomer(self, monomer: Monomer):
        idx = self._monomers.index(monomer)
        if idx >= 0:
            del self._monomers[idx]
            return True
        return False

    def get_avg_length_by_axis(self, required_axis: Axis):
        side_backward = get_reversed_side(self.main_direction)

        lengths = []
        used_monomers = list[Monomer]()
        for mon in self.monomers:
            if mon in used_monomers:
                continue
            
            curr_monomer = mon
            while curr_monomer.sides[side_backward] is not None and curr_monomer.sides[side_backward] in self.monomers:
                curr_monomer = curr_monomer.sides[side_backward]

            used_monomers.append(curr_monomer)
            length = 1
            while curr_monomer.sides[self.main_direction] is not None and curr_monomer.sides[self.main_direction] in self.monomers:
                curr_monomer = curr_monomer.sides[self.main_direction]
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