from math import sqrt
from typing import Iterable
from alg.config import Side, get_movement_sides

def distance(point1: Iterable, point2: Iterable):
    if len(point1) != len(point2):
        raise ArithmeticError(f"The input points' dimentions must be equal (len of point1 is {len(point1)}, len of point2 is {len(point2)})")
    return sqrt(sum([(e - s) ** 2 for s, e in zip(point1, point2)])) 

def get_side_by_monomers(from_: Iterable, to_: Iterable):
    assert len(from_) == len(to_), "In get_side_by_monomers: len(from_) != len(to_)"
    if len(from_) == 0 or len(to_) == 0:
        return Side.Undefined

    result = [t - f for f, t in zip(from_, to_)]
    if result[0] == 0:
        if result[1] == 0:
            if result[2] == 0:
                return Side.Undefined
            elif result[2] > 0:
                return Side.Up
            else:
                return Side.Down
        elif result[1] > 0:
            if result[2] == 0:
                return Side.Left
            elif result[2] > 0:
                return Side.LeftUp
            else:
                return Side.LeftDown
        elif result[1] < 0:
            if result[2] == 0:
                return Side.Right
            elif result[2] > 0:
                return Side.RightUp
            else:
                return Side.RightDown
    elif result[0] > 0:
        if result[1] == 0:
            if result[2] == 0:
                return Side.Forward
            elif result[2] > 0:
                return Side.UpForward
            else:
                return Side.DownForward
        elif result[1] > 0:
            if result[2] == 0:
                return Side.LeftForward
            elif result[2] > 0:
                return Side.UpLeftForward
            else:
                return Side.DownLeftForward
        elif result[1] < 0:
            if result[2] == 0:
                return Side.RightForward
            elif result[2] > 0:
                return Side.UpRightForward
            else:
                return Side.DownRightForward
    else: # result[0] < 0
        if result[1] == 0:
            if result[2] == 0:
                return Side.Backward
            elif result[2] > 0:
                return Side.UpBackward
            else:
                return Side.DownBackward
        elif result[1] > 0:
            if result[2] == 0:
                return Side.LeftBackward
            elif result[2] > 0:
                return Side.UpLeftBackward
            else:
                return Side.DownLeftBackward
        elif result[1] < 0:
            if result[2] == 0:
                return Side.RightBackward
            elif result[2] > 0:
                return Side.UpRightBackward
            else:
                return Side.DownRightBackward