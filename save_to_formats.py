from alg.config import Axis, ConnectionType, get_all_sides
from polymer_view import PolymerView, GlobulaView
from space import Space
from alg.polymer_lib import MonomerType
from alg.monomer_lib import monomer_type_to_literal
from math import sqrt

def bond_type_mass(type: ConnectionType):
    if type == ConnectionType.TypeOne:
        return 1
    elif type == ConnectionType.TypeTwo:
        return sqrt(2)
    elif type == ConnectionType.TypeThree:
        return sqrt(3)
    else:
        assert False, "Unknown connection type"

class SaveToFile:
    def __init__(self, file_label, file_extension):
        self.file_label = file_label
        self.file_extension = file_extension
        self._finished_polimers_labels = ['H', 'He', 'Li', 'Be', 'B', 'C', 'N', 'O', 'F', 'Ne', 'Na', 'Mg', 'Al', 'Si', 'P', 'S', 'Cl', 'Ar', 'K', 'Ca', 'Sc', 'Ti', 'V', 'Cr', 'Mn', 'Fe', 'Co', 'Ni', 'Cu', 'Zn', 'Ga', 'Ge', 'As', 'Se', 'Br', 'Kr', 'Rb', 'Sr', 'Y', 'Zr', 'Nb', 'Mo', 'Tc', 'Ru', 'Rh', 'Pd', 'Ag', 'Cd', 'In', 'Sn', 'Sb', 'Te', 'I', 'Xe']
        self.contents = ""
    
    def convert(self, globula: GlobulaView):
        self._build_contents(globula)
        return f'{globula.name}_{self.file_label}.{self.file_extension}', self.contents

    def _build_contents(self, globula: GlobulaView):
        pass

    def add_string(self, string: str):
        self.contents += f'{string}\n'
    
    def add_new_line(self):
        self.add_string("")
    
class SaveToMol(SaveToFile):
    def __init__(self):
        super().__init__("mol", "mol2")

    def convert_monomer(self, polymer: PolymerView, element_name: str, number: int):
        contents = ""
        for i, monomer in enumerate(polymer):
            label = element_name
            if i == 0:
                label = 'Fe'
            if i == polymer.len() - 1:
                label = 'Cu'
            x: float = monomer[Axis.X_AXIS.value]
            y: float = monomer[Axis.Y_AXIS.value]
            z: float = monomer[Axis.Z_AXIS.value]
            self.add_string(f"{polymer.len() * number + i + 1} {label} {x:.4f} {y:.4f} {z:.4f} {label}")

        return contents
    
    def _build_contents(self, globula: GlobulaView):
        self.add_string("@<TRIPOS>MOLECULE")
        self.add_string("*****")
        self.add_string(f" {sum([pl.len() for pl in globula])} {sum([pl.len() for pl in globula]) - globula.len()} 0 0 0")
        self.add_string("SMALL")
        self.add_string("GASTEIGER\n\n\r")

        self.add_string("@<TRIPOS>ATOM")
        for pol_number, pol in enumerate(globula):
            self.convert_monomer(pol, self._finished_polimers_labels[pol_number], pol_number)

        self.add_string("@<TRIPOS>BOND")
        addition = 0
        for pol_number, pol in enumerate(globula):
            for i in range(pol.len() - 1):
                self.add_string(f"{addition + i + 1} {addition + i + 1} {addition + i + 2} 1")
            addition += pol.len()

        
class SaveToLammps(SaveToFile):
    def __init__(self):
        super().__init__("lammps", "data")

    def __get_atoms_count(self, globula: GlobulaView):
        atoms_count = 0
        for pol in globula:
            for mon in pol:
                if mon.type != MonomerType.Undefined:
                    atoms_count += 1
        return atoms_count

    def __get_bonds_count(self, globula: GlobulaView):
        bonds_count = 0
        for pol in globula:
            start_mon, end_mon = pol.get_start_end_monomers()
            curr_mon = start_mon
            while curr_mon != end_mon:
                if curr_mon.type != MonomerType.Undefined \
                    and curr_mon.next_monomer != None \
                    and curr_mon.next_monomer.type != MonomerType.Undefined:
                        bonds_count += 1
                curr_mon = curr_mon.next_monomer

        return bonds_count

    def __get_monomer_types(self, globula: GlobulaView):
        types_set = set[MonomerType]()
        for pol in globula:
            for monomer in pol:
                types_set.add(monomer.type)

        try:
            types_set.remove(MonomerType.Undefined)
        except KeyError as e:
            return types_set
        
        return types_set

    def __get_bond_types(self, globula: GlobulaView):
        types_set = set[ConnectionType]()
        all_sides = get_all_sides()
        bonds_count = 0
        for pol in globula:
            for mon in pol:
                for side in all_sides:
                    type = mon.get_type_of_connection_with_side(side)
                    sibling = mon.get_sibling(side)
                    if type != ConnectionType.TypeUndefined and sibling is not None and sibling.number != -1:
                        types_set.add(type)
                        bonds_count += 1
                        
        return types_set, int(bonds_count / 2)
                

    def _build_contents(self, globula: GlobulaView):
        self.add_string("LAMMPS data file via write_data, version 24 Dec 2020, timestep = 40000000")
        self.add_new_line()

        monomer_types_set = self.__get_monomer_types(globula)
        bond_types_set, bonds_count = self.__get_bond_types(globula)
        self.add_string(f'{self.__get_atoms_count(globula)} atoms')
        self.add_string(f'{len(monomer_types_set)} atom types')
        self.add_string(f'{bonds_count} bonds')
        self.add_string(f'{len(bond_types_set)} bond types')
        self.add_new_line()

        self.add_string(f'0 {Space.space_dimention} xlo xhi')
        self.add_string(f'0 {Space.space_dimention} ylo yhi')
        self.add_string(f'0 {Space.space_dimention} zlo zhi')
        self.add_new_line()
        
        self.add_string("Masses")
        self.add_new_line()

        d_types = { type: i + 1 for i, type in enumerate(monomer_types_set)}

        for type, i in d_types.items():
            if type == MonomerType.Undefined:
                continue
            try:
                self.add_string(f'{i} 1 # {monomer_type_to_literal(type)}')
            except Exception as err:
                print(err.with_traceback(None))
                return

        self.add_new_line()

        self.add_string('Bond Coeffs # harmonic')
        self.add_new_line()

        for i, bond_type in enumerate(bond_types_set):
            self.add_string(f'{i + 1} 10000 # {bond_type_mass(bond_type)}')
        self.add_new_line()

        self.add_string('Atoms # full')
        self.add_new_line()

        for pol_number, pol in enumerate(globula):
            for monomer in pol:
                if monomer.type == MonomerType.Undefined:
                    continue
                self.add_string(f'{monomer.number} 1 {d_types[monomer.type]} 0.00000 {monomer[Axis.X_AXIS.value]} {monomer[Axis.Y_AXIS.value]} {monomer[Axis.Z_AXIS.value]} 0 0 0')
        self.add_new_line()

        self.add_string('Bonds')
        self.add_new_line()

        bond_number = 1
        used_pairs = set()
        all_sides = get_all_sides()
        for pol in globula:
            for mon in pol:
                for side in all_sides:
                    other_mon= mon.get_sibling(side)
                    if other_mon is not None and other_mon.number != -1:
                        pair = tuple(sorted((mon.number, other_mon.number)))
                        if pair not in used_pairs:
                            self.add_string(f'{bond_number} 1 {mon.number} {other_mon.number}')
                            bond_number += 1
                            used_pairs.add(pair)

        self.add_new_line()
        self.add_string('Angles')
        self.add_new_line()

        angle_number = 1
        for pol in globula:
            for i in range(1, pol.len() - 1):
                self.add_string(f'{angle_number} 1 {i - 1} {i} {i + 1}')
                angle_number += 1