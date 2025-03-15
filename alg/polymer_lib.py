from alg.config import axis_count, MoveDirection, Axis, move_cell, move_cell_no_borders
from alg.field_lib import Field
from space import Space
from alg.common_funcs import distance
from enum import IntEnum

class MonomerType(IntEnum):
    Usual = 0
    Owise = 1

class Monomer:
    def __init__(self, coords: list[int], type: MonomerType):
        self._coords = coords
        self._type = type
        self._prev_monomer = None
        self._next_monomer = None

    @property
    def coords(self):
        return self._coords

    @property
    def type(self):
        return self._type
    
    @type.setter
    def type(self, type):
        self._type = type

    @property
    def prev_monomer(self):
        return self._prev_monomer
    
    @prev_monomer.setter
    def prev_monomer(self, prev):
        self._prev_monomer = prev

    @property
    def next_monomer(self):
        return self._next_monomer
    
    @next_monomer.setter
    def next_monomer(self, next):
        self._next_monomer = next
    
    def __iter__(self):
        return iter(self._coords)

    def __len__(self):
        return len(self._coords)

    def __getitem__(self, i):
        return self._coords[i]

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
    
    def add_monomer(self, monomer):
        m = Monomer(monomer, MonomerType.Usual)
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
        next_points = []
        for axis in range(axis_count()):
            for direction in [MoveDirection.DIRECTION_BACKWARD, MoveDirection.DIRECTION_FORWARD]:
                next_points.append(move_cell_no_borders(last_point, Axis(axis), direction))
        for next_point in next_points:
            if Space.point_within_borders(next_point) and not self.__field.is_free(next_point) and next_point != prelast_point:
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