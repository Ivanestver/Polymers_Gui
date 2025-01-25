import math
from typing import Iterable
from alg.config import axis_count, Direction, Axis, move_cell, move_cell_no_borders
import numpy as np
from space import Space

class Field:
    def __init__(self, sphere_radius):
        shape = tuple([Space.space_dimention for i in range(axis_count())])
        self._field = np.zeros(shape, dtype=int)
        self._sphere_radius = sphere_radius
        
    def __set_cell_state(self, position_to_set: Iterable, state: int):
        item = self._field
        for i in range(axis_count()-1):
            item = item[position_to_set[i]]
        item[position_to_set[axis_count()-1]] = state

    def make_filled(self, position_to_make: Iterable):
        self.__set_cell_state(position_to_make, 1)
    
    def make_free(self, position_to_make: Iterable):
        self.__set_cell_state(position_to_make, 0)

    def is_free(self, coords: Iterable) -> bool:
        item = self._field
        for i in range(len(coords)-1):
            if coords[i] < 0 or coords[i] >= len(item):
                return False
            item = item[coords[i]]

        if coords[-1] < 0 or coords[-1] >= len(item):
            return False
        return item[coords[-1]] == 0

    def get_cell_within_borders(self, point):
        return tuple((v % Space.space_dimention for v in point))

    def __get_length_metrics(self, start_point):
        end_point = Space.global_zero
        return math.sqrt(sum([(e - s) ** 2 for e, s in zip(end_point, start_point)]))

    def is_within_sphere(self, point):
        return self.__get_length_metrics(point) <= self._sphere_radius
    
    def is_within_borders(self, point: tuple):
        return (0 <= point[0] <= Space.space_dimention) \
        and (0 <= point[1] <= Space.space_dimention) \
        and (0 <= point[2] <= Space.space_dimention) \
        and self.is_within_sphere(point)

    def get_available_cells(self, curr_pos):
        available_cells = []
        for axis in range(axis_count()):
            for dir in [Direction.DIRECTION_BACKWARD, Direction.DIRECTION_FORWARD]:
                curr_moved = move_cell_no_borders(curr_pos, Axis(axis), dir)
                if self.is_within_borders(curr_moved) and self.is_free(curr_moved):
                    available_cells.append(curr_moved)

        return available_cells
    
    def define_start_position(self):
        start_position = tuple([int(Space.space_dimention / 2) for i in range(axis_count())])
        while not self.is_free(start_position):
            available_cells = self.get_available_cells(start_position)
            if len(available_cells) != 0:
                return available_cells[np.random.randint(0, len(available_cells))]
            random_axis = np.random.randint(0, Axis.AXIS_COUNT.value)
            random_direction = np.random.randint(0, Direction.DIRECTION_COUNT.value)
            start_position = move_cell(start_position, Axis(random_axis), Direction.DIRECTION_FORWARD if random_direction == 0 else Direction.DIRECTION_BACKWARD)
        return start_position