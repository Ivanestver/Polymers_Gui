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
    UpLeftForward = 1000
    UpForward = UpLeftForward + 1
    UpRightForward = UpForward + 1
    LeftForward = UpRightForward + 1
    RightForward = LeftForward + 1
    DownLeftForward = RightForward + 1
    DownForward = DownLeftForward + 1
    DownRightForward = DownForward + 1
    LeftUp = DownRightForward + 1
    LeftDown = LeftUp + 1
    RightUp = LeftDown + 1
    RightDown = RightUp + 1
    UpLeftBackward = RightDown + 1
    UpBackward = UpLeftBackward + 1
    UpRightBackward = UpBackward + 1
    LeftBackward = UpRightBackward + 1
    RightBackward = LeftBackward + 1
    DownLeftBackward = RightBackward + 1
    DownBackward = DownLeftBackward + 1
    DownRightBackward = DownBackward + 1

class MonomerType(enum.IntEnum):
    Undefined = -1
    Usual = 0
    Owise = 1
    Nwise = 2
    Fwise = 3,
    Clwise = 4

class ConnectionType(enum.IntEnum):
    TypeUndefined = 0
    TypeOne = 1
    TypeTwo = 2
    TypeThree = 3
    ConnectionTypeCount = 3

def get_movement_sides():
    return [Side.Left, Side.Right, Side.Forward, Side.Backward, Side.Up, Side.Down]
    
def get_surface_diagonal_sides():
    return [Side.UpForward, Side.UpBackward, Side.DownForward, Side.DownBackward, Side.LeftUp, Side.LeftDown, Side.RightUp, Side.RightDown, Side.LeftForward, Side.LeftBackward, Side.RightForward, Side.RightBackward]

def get_cube_diagonal_sides():
    return [Side.UpLeftForward, Side.UpLeftBackward, Side.UpRightForward, Side.UpRightBackward, Side.DownLeftForward, Side.DownLeftBackward, Side.DownRightForward, Side.DownRightBackward]

def get_additional_sides():
    return get_surface_diagonal_sides() + get_cube_diagonal_sides()

def get_all_sides():
    return get_movement_sides() + get_additional_sides()

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
    elif side == Side.Right:
        return Side.Left
    elif side == Side.UpLeftForward:
        return Side.DownRightBackward
    elif side == Side.UpForward:
        return Side.DownBackward
    elif side == Side.UpRightForward:
        return Side.DownLeftBackward
    elif side == Side.LeftForward:
        return Side.RightBackward
    elif side == Side.RightForward:
        return Side.LeftBackward
    elif side == Side.DownLeftForward:
        return Side.UpRightBackward
    elif side == Side.DownForward:
        return Side.UpBackward
    elif side == Side.DownRightForward:
        return Side.UpLeftBackward
    elif side == Side.LeftUp:
        return Side.RightDown
    elif side == Side.LeftDown:
        return Side.RightUp
    elif side == Side.RightUp:
        return Side.LeftDown
    elif side == Side.RightDown:
        return Side.LeftUp
    elif side == Side.UpLeftBackward:
        return Side.DownRightForward
    elif side == Side.UpBackward:
        return Side.DownForward
    elif side == Side.UpRightBackward:
        return Side.DownLeftForward
    elif side == Side.LeftBackward:
        return Side.RightForward
    elif side == Side.RightBackward:
        return Side.LeftForward
    elif side == Side.DownLeftBackward:
        return Side.UpRightForward
    elif side == Side.DownBackward:
        return Side.UpForward
    elif side == Side.DownRightBackward:
        return Side.UpLeftForward
    else:
        assert(False), "Unappropriate side"

def axis_count():
    return Axis.AXIS_COUNT.value

def get_axis_color(axis: Axis):
    return {
        Axis.X_AXIS: MonomerType.Owise,
        Axis.Y_AXIS: MonomerType.Nwise,
        Axis.Z_AXIS: MonomerType.Fwise
    }[axis]

max_monomers_count = 5