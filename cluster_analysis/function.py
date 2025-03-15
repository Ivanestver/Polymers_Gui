from enum import IntEnum
from operator import sub

class Direction(IntEnum):
    Left = 0 * 100 + 0 * 10 + 1,
    Right = 0 * 100 + 0 * 10 + -1,
    Up = 0 * 100 + 1 * 10 + 0,
    Down = 0 * 100 + -1 * 10 + 0,
    Forward = 1 * 100 + 0 * 10 + 0,
    Backward = -1 * 100 + 0 * 10 + 0

def get_direction(prev: tuple, curr: tuple) -> Direction:
    tuples_sub = tuple(map(sub, curr, prev))
    return Direction(sum(tuples_sub[i] * (10 ** i) for i in range(len(tuples_sub))))