import enum

class Axis(enum.Enum):
    X_AXIS = 0
    Y_AXIS = 1
    Z_AXIS = 2
    AXIS_COUNT = 3

class MoveDirection(enum.Enum):
    DIRECTION_FORWARD = 1
    DIRECTION_BACKWARD = -1
    DIRECTION_COUNT = 2

class Side(enum.Enum):
    Undefined = -10
    Forward = (Axis.X_AXIS.value + 1) * MoveDirection.DIRECTION_FORWARD.value
    Backward = (Axis.X_AXIS.value + 1) * MoveDirection.DIRECTION_BACKWARD.value
    Left = (Axis.Y_AXIS.value + 1) * MoveDirection.DIRECTION_FORWARD.value
    Right = (Axis.Y_AXIS.value + 1) * MoveDirection.DIRECTION_BACKWARD.value
    Up = (Axis.Z_AXIS.value + 1) * MoveDirection.DIRECTION_FORWARD.value
    Down = (Axis.Z_AXIS.value + 1) * MoveDirection.DIRECTION_BACKWARD.value

def get_side(axis: Axis, move_direction: MoveDirection) -> Side:
    return Side((axis.value + 1) * move_direction.value)

def get_reversed_side(side: Side) -> Side:
    if side == Side.Forward:
        return Side.Backward
    elif side == Side.Backward:
        return Side.Forward
    elif side == Side.Up:
        return Side.Down
    elif side == Side.Down:
        return Side.Up
    elif side == Side.Left:
        return Side.Right
    else:
        return Side.Left

def axis_count():
    return Axis.AXIS_COUNT.value

max_monomers_count = 5