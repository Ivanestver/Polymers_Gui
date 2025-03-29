from enum import IntEnum
from alg.config import Side

class MonomerType(IntEnum):
    Undefined = -1
    Usual = 0
    Owise = 1
    Nwise = 2
    Fwise = 3

def monomer_type_to_literal(type: MonomerType):
    if type == MonomerType.Undefined:
        raise Exception("There's no letter for Undefined monomer type")
    elif type == MonomerType.Usual:
        return 'C'
    elif type == MonomerType.Owise:
        return 'O'
    elif type == MonomerType.Nwise:
        return 'N'
    elif type == MonomerType.Fwise:
        return 'F'
    else:
        raise Exception(f"There's no letter for type {type}")

class Monomer:
    def __init__(self, coords: list[int], type: MonomerType = MonomerType.Undefined):
        self._coords = coords
        self._type = type
        self._prev_monomer = None
        self._next_monomer = None
        self.sides = { side: None for side in Side }

    def is_of_type(self, type: MonomerType):
        return self.type == type

    def is_not_of_type(self, type: MonomerType):
        return self.type != type

    def get_sibling(self, side: Side):
        return self.sides[side]

    @property
    def coords(self):
        return self._coords

    @property
    def type(self):
        return self._type
    
    @type.setter
    def type(self, type):
        self._type = type

    @property
    def prev_monomer(self):
        return self._prev_monomer
    
    @prev_monomer.setter
    def prev_monomer(self, prev):
        self._prev_monomer = prev

    @property
    def next_monomer(self):
        return self._next_monomer
    
    @next_monomer.setter
    def next_monomer(self, next):
        self._next_monomer = next
    
    def __iter__(self):
        return iter(self._coords)

    def __len__(self):
        return len(self._coords)

    def __getitem__(self, i):
        return self._coords[i]