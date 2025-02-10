from alg.config import Axis
from polymer_view import PolymerView, GlobulaView

class SaveToFile:
    def __init__(self, file_label, file_extension):
        self.file_label = file_label
        self.file_extension = file_extension
        self._finished_polimers_labels = ['C', 'N', 'H', 'O', 'F', 'Na', 'Mg', 'Al', 'P', 'S']
        self.contents = ""
    
    def convert(self, globula: GlobulaView):
        self._build_contents(globula)
        return f'{globula.name}_{self.file_label}.{self.file_extension}', self.contents

    def _build_contents(self, globula: GlobulaView):
        pass

    def add_string(self, string: str = ""):
        self.contents += f'{string}\n'
    
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

    def _build_contents(self, globula: GlobulaView):
        self.add_string("LAMMPS data file via write_data, version 24 Dec 2020, timestep = 40000000")
        self.add_string()

        atoms_count = self.__get_atoms_count(globula)
        self.add_string(f'{atoms_count} atoms')
        self.add_string(f'{globula.len()} atom types')
        self.add_string(f'{atoms_count - 1} bonds')
        self.add_string(f'1 bond types')
        self.add_string()