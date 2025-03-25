from enum import IntEnum
from operator import sub
from alg.config import Side

class Direction(IntEnum):
    Left = 0 * 100 + 0 * 10 + 1,
    Right = 0 * 100 + 0 * 10 + -1,
    Forward = 0 * 100 + 1 * 10 + 0,
    Backward = 0 * 100 + -1 * 10 + 0
    Up = 1 * 100 + 0 * 10 + 0,
    Down = -1 * 100 + 0 * 10 + 0,

def get_direction(prev: tuple, curr: tuple) -> Direction:
    tuples_sub = tuple(map(sub, curr, prev))
    s = 0.0
    for i in range(len(tuples_sub)):
        value = tuples_sub[-(i+1)] * (10 ** i)
        s += value
    return Direction(s)

def translate_direction_to_side(direction: Direction) -> Side:
    return {
        Direction.Left: Side.Left,
        Direction.Right: Side.Right,
        Direction.Up: Side.Up,
        Direction.Down: Side.Down,
        Direction.Forward: Side.Forward,
        Direction.Backward: Side.Backward,
    }[direction]

def get_side(prev: tuple, curr: tuple) -> Side:
    return translate_direction_to_side(get_direction(prev, curr))
