import enum
from space import Space

class Axis(enum.Enum):
    X_AXIS = 0
    Y_AXIS = 1
    Z_AXIS = 2
    AXIS_COUNT = 3

class Direction(enum.Enum):
    DIRECTION_FORWARD = 1
    DIRECTION_BACKWARD = -1
    DIRECTION_COUNT = 2

def axis_count():
    return Axis.AXIS_COUNT.value

max_monomers_count = 5

def move_cell(cell, axis: Axis, direction: Direction):
    return tuple((cell[i] if i != axis.value else cell[i] + direction.value) % Space.space_dimention for i in range(len(cell)))

def move_cell_no_borders(cell, axis: Axis, direction: Direction):
    return tuple((cell[i] if i != axis.value else cell[i] + direction.value) for i in range(len(cell)))