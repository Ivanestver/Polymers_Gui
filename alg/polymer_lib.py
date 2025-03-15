from alg.config import Axis, Side
from alg.field_lib import Field
from space import Space
from alg.common_funcs import distance
from alg.monomer_lib import Monomer, MonomerType

class Polymer:
    __number__ = 0
    def __init__(self, field: Field, number = None):
        self.__field = field
        self.__polymer: list[Monomer] = []
        if number is None:
            self.__number = Polymer.__number__
            Polymer.__number__ += 1
        else:
            self.__number = number
        
    def __iter__(self):
        return iter(self.__polymer)
    
    def add_monomer(self, m: Monomer):
        if len(self.__polymer) != 0:
            self.__polymer[-1].next_monomer = m
            m.prev_monomer = self.__polymer[-1]

        self.__polymer.append(m)
        
    def len(self):
        return len(self.__polymer)

    def front(self):
        return self.__polymer[0]
    
    def back(self):
        return self.__polymer[-1]
    
    def name(self):
        return f"Polymer {self.__number}"
    
    def number(self):
        return self.__number

    def field(self):
        return self.__field
    
    def copy(self):
        new_polymer = Polymer(self.__field, self.number())
        new_polymer.__polymer = self.__polymer.copy()
        return new_polymer
        
    def get_monomer_idx(self, curr_monomer):
        for i, monomer in enumerate(self.__polymer):
            if all([item.left == item.right for item in list(zip(curr_monomer, monomer))]):
                return i
        return -1

    def get_monomer_by_idx(self, idx: int):
        if (idx < 0 or idx >= len(self.__polymer)):
            return None
        return self.__polymer[idx]

    def calc_energy(self):
        u = 0.0
        last_point = self.__polymer[-1]
        if len(self.__polymer) > 1:
            prelast_point = self.__polymer[-2]
        else:
            prelast_point = last_point
        for side in Side:
            if last_point.sides[side] is not None and not last_point.sides[side].is_of_type(MonomerType.Undefined) and last_point.sides[side] != prelast_point:
                u += -1.0

        return u

    def calc_whole_energy(self):
        u = 0.0
        for i in range(len(self.__polymer)):
            for j in range(len(self.__polymer)):
                if i == j:
                    continue
                    
                r_ij = 1 / distance(self.__polymer[i], self.__polymer[j])
                u += 4 * 0.01 * ((r_ij ** 12) - (r_ij ** 6))

        return u

    def make_step_back(self):
        if (len(self.__polymer) <= 1):
            return False
        self.__polymer.remove(self.__polymer[-1])
        return True

    def get_min_max_width_height(self):
        min_width = Space.space_dimention
        max_width = 0
        min_height = Space.space_dimention
        max_height = 0

        for monomer in self.__polymer:
            min_width = min(min_width, monomer[Axis.X_AXIS.value])
            max_width = max(max_width, monomer[Axis.Y_AXIS.value])
            min_height = min(min_height, monomer[Axis.X_AXIS.value])
            max_height = max(max_height, monomer[Axis.Y_AXIS.value])
        return (float(min_width), float(max_width), float(min_height), float(max_height))