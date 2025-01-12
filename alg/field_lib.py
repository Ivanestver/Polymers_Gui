from typing import Iterable
from config import DIMENTION, axis_count, Direction, Axis, move_cell
import numpy as np

class Field:
    def __init__(self):
        shape = tuple([DIMENTION for i in range(axis_count())])
        self.__field = np.zeros(shape, dtype=int)
        
    def __set_cell_state(self, position_to_set: Iterable, state: int):
        item = self.__field
        for i in range(axis_count()-1):
            item = item[position_to_set[i]]
        item[position_to_set[axis_count()-1]] = state

    def make_filled(self, position_to_make: Iterable):
        self.__set_cell_state(position_to_make, 1)
    
    def make_free(self, position_to_make: Iterable):
        self.__set_cell_state(position_to_make, 0)

    def is_free(self, coords: Iterable) -> bool:
        item = self.__field
        for i in range(len(coords)-1):
            item = item[coords[i]]
        return item[coords[-1]] == 0

    def get_cell_within_borders(self, point):
        return tuple((v % DIMENTION for v in point))

    def get_available_cells(self, curr_pos):
        available_cells = []
        for axis in range(axis_count()):
            for dir in [Direction.DIRECTION_BACKWARD, Direction.DIRECTION_FORWARD]:
                curr_moved = self.get_cell_within_borders(move_cell(curr_pos, Axis(axis), dir))
                if self.is_free(curr_moved):
                    available_cells.append(curr_moved)

        return available_cells
    
    def define_start_position(self):
        start_position = tuple([int(DIMENTION / 2) for i in range(axis_count())])
        while not self.is_free(start_position):
            available_cells = self.get_available_cells(start_position)
            if len(available_cells) != 0:
                return available_cells[np.random.randint(0, len(available_cells))]
            random_axis = np.random.randint(0, Axis.AXIS_COUNT.value)
            random_direction = np.random.randint(0, Direction.DIRECTION_COUNT.value)
            start_position = move_cell(start_position, Axis(random_axis), Direction.DIRECTION_FORWARD if random_direction == 0 else Direction.DIRECTION_BACKWARD)
        return start_position