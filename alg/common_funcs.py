from math import sqrt
from typing import Iterable

def distance(point1: Iterable, point2: Iterable):
    if len(point1) != len(point2):
        raise ArithmeticError(f"The input points' dimentions must be equal (len of point1 is {len(point1)}, len of point2 is {len(point2)})")
    return sqrt(sum([(e - s) ** 2 for s, e in zip(point1, point2)])) 