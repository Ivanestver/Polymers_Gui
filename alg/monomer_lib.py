from alg.config import Side, MonomerType, ConnectionType, get_reversed_side, get_movement_sides, get_surface_diagonal_sides, get_cube_diagonal_sides
from alg.common_funcs import get_side_by_monomers

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
    elif type == MonomerType.Clwise:
        return 'Cl'
    else:
        raise Exception(f"There's no letter for type {type}")

class Connection:
    def __init__(self, mon1, mon2, type: ConnectionType):
        self.__monomers = list([mon1, mon2])
        self.__type = type

    def get_other_side(self, curr_mon):
        idx = self.__monomers.index(curr_mon)
        if idx < 0:
            raise Exception("Connection is corrupted")
        return self.__monomers[(idx + 1) % len(self.__monomers)]

    @property
    def type(self):
        return self.__type

    @type.setter
    def type(self, other_type):
        self.__type = other_type

class Monomer:
    def __init__(self, coords: list[int], type: MonomerType = MonomerType.Undefined):
        self._coords = coords
        self._type = type
        self._prev_monomer = None
        self._next_monomer = None
        self.sides: dict[Side, Connection] = { side: None for side in Side }
        self.__number = -1

    def is_of_type(self, type: MonomerType):
        return self.type == type

    def is_not_of_type(self, type: MonomerType):
        return self.type != type

    def get_sibling(self, side: Side):
        conn = self.sides[side]
        return conn.get_other_side(self) if conn is not None else None

    def get_side_of_sibling(self, other_mon):
        for side, conn in self.sides.items():
            if conn is not None and other_mon == conn.get_other_side(self):
                return side
        return Side.Undefined

    def get_type_of_connection_with_side(self, side: Side):
        conn = self.sides[side]
        return conn.type if conn is not None else ConnectionType.TypeUndefined

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

    @property
    def number(self):
        return self.__number
    
    @number.setter
    def number(self, value):
        self.__number = value
    
    def __iter__(self):
        return iter(self._coords)

    def __len__(self):
        return len(self._coords)

    def __getitem__(self, i):
        return self._coords[i]

        
def make_connection(mon1: Monomer, mon2: Monomer, type: ConnectionType):
    if mon1 is None or mon2 is None:
        return

    side = get_side_by_monomers(mon1.coords, mon2.coords)
    conn = mon1.sides[side]
    reversed_side = get_reversed_side(side)
    if conn is None:
        new_conn = Connection(mon1, mon2, type)
        mon1.sides[side] = new_conn
        mon2.sides[reversed_side] = new_conn
    else:
        other_side = conn.get_other_side(mon1)
        if mon2 != other_side:
            raise Exception("Two monomers occupy the same location")
        conn.type = type

def break_connection(mon1: Monomer, mon2: Monomer, side: Side):
    if side in get_movement_sides():
        return

    reversed_side = get_reversed_side(side)
    conn = mon1.sides[side]
    if conn != mon2.sides[reversed_side]:
        raise Exception("Two differenct connections between the same monomers exist")
    mon1.sides[side] = None
    mon2.sides[reversed_side] = None

def get_connection_type(side_one: Monomer, side_two: Monomer, side: Side = None):
    if side is None:
        side = get_side_by_monomers(side_one, side_two)
    if side in get_movement_sides():
        return ConnectionType.TypeOne
    elif side in get_surface_diagonal_sides():
        return ConnectionType.TypeTwo
    elif side in get_cube_diagonal_sides():
        return ConnectionType.TypeThree
    else:
        return ConnectionType.TypeUndefined