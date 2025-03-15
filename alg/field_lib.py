import math
from typing import Iterable
from alg.config import axis_count, MoveDirection, Axis, Side, get_side
import numpy as np
from space import Space
from alg.monomer_lib import Monomer, MonomerType

class Field:
    def __init__(self, sphere_radius):
        self._sphere_radius = sphere_radius
        self.__init_field()

    def __init_field(self):
        shape = tuple([Space.space_dimention for i in range(axis_count())])
        self._field = [
            [
                [
                    Monomer((i, j, k)) for k in range(shape[2])
                    ] for j in range(shape[1])
            ] for i in range(shape[0])
        ] 

        for i in range(shape[0]):
            for j in range(shape[1]):
                for k in range(shape[2]):
                    monomer = self._field[i][j][k]
                    monomer.sides[Side.Forward] = self._field[i + 1][j][k] if i < shape[0] - 1 else None
                    monomer.sides[Side.Backward] = self._field[i - 1][j][k] if i > 0 else None
                    monomer.sides[Side.Up] = self._field[i][j + 1][k] if j < shape[1] - 1 else None
                    monomer.sides[Side.Down] = self._field[i][j - 1][k] if j > 0 else None
                    monomer.sides[Side.Left] = self._field[i][j][k + 1] if k < shape[2] - 1 else None
                    monomer.sides[Side.Right] = self._field[i][j][k - 1] if k > 0 else None
        
    def __set_cell_state(self, position_to_set: Iterable, state: MonomerType):
        item = self._field
        for i in range(axis_count()-1):
            item = item[position_to_set[i]]
        item[position_to_set[axis_count()-1]] = state

    def make_filled(self, monomer: Monomer):
        monomer.type = MonomerType.Usual
    
    def make_free(self, monomer: Monomer):
        monomer.type = MonomerType.Undefined

    def is_free(self, coords: Iterable) -> bool:
        item = self._field
        for i in range(len(coords)-1):
            if coords[i] < 0 or coords[i] >= len(item):
                return False
            item = item[coords[i]]

        if coords[-1] < 0 or coords[-1] >= len(item):
            return False
        return item[coords[-1]].is_of_type(MonomerType.Undefined)

    def get_cell_within_borders(self, point):
        return tuple((v % Space.space_dimention for v in point))

    def is_busy(self):
        for i in self._field:
            for j in i:
                for k in j:
                    if k.is_of_type(MonomerType.Undefined):
                        return False
        True

    def __get_length_metrics(self, start_point):
        end_point = Space.global_zero
        return math.sqrt(sum([(e - s) ** 2 for e, s in zip(end_point, start_point)]))

    def is_within_sphere(self, point):
        return self.__get_length_metrics(point) <= self._sphere_radius
    
    def is_within_borders(self, point: tuple):
        return Space.point_within_borders(point) and self.is_within_sphere(point)

    def get_monomer_by_coords(self, coords: tuple) -> Monomer:
        item = self._field
        for i in range(axis_count()-1):
            item = item[coords[i]]
        return item[coords[-1]]

    def get_available_cells(self, curr_pos: tuple):
        available_cells: list[Monomer] = []
        monomer = self.get_monomer_by_coords(curr_pos)
        for side in Side:
            if monomer.sides[side] is not None and monomer.sides[side].is_of_type(MonomerType.Undefined):
                available_cells.append(monomer.sides[side])

        return available_cells
    
    def define_start_position(self) -> Monomer:
        start_position = tuple([int(Space.space_dimention / 2) for i in range(axis_count())])
        while not self.is_free(start_position):
            available_cells = self.get_available_cells(start_position)
            if len(available_cells) != 0:
                return available_cells[np.random.randint(0, len(available_cells))]
            random_axis = np.random.randint(0, Axis.AXIS_COUNT.value)
            random_direction = np.random.randint(0, MoveDirection.DIRECTION_COUNT.value)
            start_position = self.get_monomer_by_coords(start_position).get_sibling(get_side(random_axis, random_direction)).coords
        return self.get_monomer_by_coords(start_position)