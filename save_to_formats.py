from alg.config import Axis
from polymer_view import PolymerView, GlobulaView
from space import Space
from alg.polymer_lib import MonomerType
from alg.monomer_lib import monomer_type_to_literal

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
            atoms_count += pol.len()
        return atoms_count

    def __get_bonds_count(self, globula: GlobulaView):
        return sum([pol.len() - 1 for pol in globula])

    def __get_monomer_types(self, globula: GlobulaView):
        types_set = set[MonomerType]()
        for pol in globula:
            for monomer in pol:
                types_set.add(monomer.type)

        return types_set

    def _build_contents(self, globula: GlobulaView):
        self.add_string("LAMMPS data file via write_data, version 24 Dec 2020, timestep = 40000000")
        self.add_new_line()

        types_set = self.__get_monomer_types(globula)
        self.add_string(f'{self.__get_atoms_count(globula)} atoms')
        self.add_string(f'{len(types_set)} atom types')
        self.add_string(f'{self.__get_bonds_count(globula)} bonds')
        self.add_string(f'1 bond types')
        self.add_new_line()

        self.add_string(f'0 {Space.space_dimention} xlo xhi')
        self.add_string(f'0 {Space.space_dimention} ylo yhi')
        self.add_string(f'0 {Space.space_dimention} zlo zhi')
        self.add_new_line()
        
        self.add_string("Masses")
        self.add_new_line()

        for i, type in enumerate(types_set):
            try:
                self.add_string(f'{i + 1} 1 # {monomer_type_to_literal(type)}')
            except Exception as err:
                print(err.with_traceback())
                return

        self.add_new_line()

        self.add_string('Bond Coeffs # harmonic')
        self.add_new_line()

        self.add_string(f'1 10000 # 1')
        self.add_new_line()

        self.add_string('Atoms # molecular')
        self.add_new_line()

        monomer_number = 1
        for pol_number, pol in enumerate(globula):
            for monomer in pol:
                self.add_string(f'{monomer_number} 1 {1 if monomer.type == MonomerType.Usual else 2} 0.00000 {monomer[Axis.X_AXIS.value]} {monomer[Axis.Y_AXIS.value]} {monomer[Axis.Z_AXIS.value]} 0 0 0')
                monomer_number += 1
        self.add_new_line()

        self.add_string('Bonds')
        self.add_new_line()

        bond_number = 1
        monomer_number = 1
        for pol in globula:
            for i in range(pol.len() - 1):
                self.add_string(f'{bond_number} 1 {monomer_number} {monomer_number + 1}')
                bond_number += 1
                monomer_number += 1
            monomer_number += 1

        angle_number = 1
        for pol in globula:
            for i in range(1, pol.len() - 1):
                self.add_string(f'{angle_number} 1 {i - 1} {i} {i + 1}')
                angle_number += 1