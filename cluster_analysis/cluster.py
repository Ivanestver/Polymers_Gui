from alg.field_lib import Monomer, MonomerType
from cluster_analysis.function import get_direction, Direction, get_side
from alg.config import Side, get_reversed_side

class Cluster:
    def __init__(self, monomer: Monomer):
        self.start_monomer = monomer
        self.side = get_side(self.start_monomer, self.start_monomer.next_monomer)
        self.monomers: list[Monomer] = []
        self.first_monomers: list[Monomer] = []

        self.sides_to_look = {
            Side.Forward: [Side.Up, Side.Down, Side.Left, Side.Right],
            Side.Backward: [Side.Up, Side.Down, Side.Left, Side.Right],
            Side.Up: [Side.Forward, Side.Backward, Side.Left, Side.Right],
            Side.Down: [Side.Forward, Side.Backward, Side.Left, Side.Right],
            Side.Left: [Side.Up, Side.Down, Side.Forward, Side.Backward],
            Side.Right: [Side.Up, Side.Down, Side.Forward, Side.Backward],
        }

    def make(self):
        sides_to_look = self.sides_to_look[self.side]
        reversed_side = get_reversed_side(self.side)
        monomers_queue = [self.start_monomer]


    def fits_constraint_avg(self, avg: float):
        def calc_chain_length(first_monomer: Monomer):
            length = 1
            curr = first_monomer
            while curr.sides[self.side] is not None and curr.sides[self.side].is_of_type(MonomerType.Nwise):
                curr = curr.sides[self.side]
                length += 1
            return length
        
        lengths = [calc_chain_length(first_monomer) for first_monomer in self.first_monomers]
        calculated_avg = sum(lengths) / len(lengths)
        return abs(avg - calculated_avg) <= 1e-3 or calculated_avg > avg

    def destroy_cluster(self):
        self.set_cluster_type(MonomerType.Usual)

    def set_cluster_type(self, type: MonomerType):
        for monomer in self.monomers:
            monomer.type = type