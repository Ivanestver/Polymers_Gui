# Protein folding
# Here the nPERMss is implemented
import numpy as np
import math
import os

from alg.config import Axis
from alg.field_lib import Field
from alg.polymer_lib import Polymer
from space import Space

class CalcAlg:
    def __init__(self, globuls_count: int = 1, polymers_count: int = 1, accept_threshold: float = 0.1, max_monomers_count = 10, sphere_radius = 1):
        self._globuls_count: int = globuls_count
        self._polymers_count: int = polymers_count
        self._accept_threshold: float = accept_threshold
        self._finished_polimers_labels = ['C', 'N', 'H', 'O', 'F', 'Na', 'Mg', 'Al', 'P', 'S']
        self._max_monomers_count = max_monomers_count
        self._sphere_radius = sphere_radius

    def __get_continuations(self, k_free, available_cells):
        chosen_continuations_idxs = np.arange(k_free)
        np.random.shuffle(chosen_continuations_idxs)
        return np.array([available_cells[i] for i in chosen_continuations_idxs])

    def __get_next_config(self, curr_config: Polymer, continuation):
        config_copy = curr_config.copy()
        config_copy.add_monomer(tuple(continuation))
        return config_copy


    def __save_polymer_mol_into_file(self, path_to_save_dir, file_number, finished_polimers: list[Polymer]):
        if not os.path.isdir(path_to_save_dir):
            os.mkdir(path_to_save_dir)

        contents = ""
        contents += "@<TRIPOS>MOLECULE\n*****\n"
        contents += f" {sum([pl.len() for pl in finished_polimers])} {sum([pl.len() for pl in finished_polimers]) - len(finished_polimers)} 0 0 0\n"
        contents += "SMALL\n"
        contents += "GASTEIGER\n\n\r\n"

        contents += "@<TRIPOS>ATOM\n"
        for pol_number, pol in enumerate(finished_polimers):
            min_width, max_width, min_height, max_height = pol.get_min_max_width_height()
            if min_width == max_width or min_height == max_height:
                return
            for i, monomer in enumerate(pol):
                label = self._finished_polimers_labels[pol_number]
                if i == 0:
                    label = 'Fe'
                if i == pol.len() - 1:
                    label = 'Cu'
                x: float = monomer[Axis.X_AXIS.value]
                y: float = monomer[Axis.Y_AXIS.value]
                z: float = monomer[Axis.Z_AXIS.value]
                contents += f"{self._max_monomers_count * pol_number + i + 1} {label} {x:.4f} {y:.4f} {z:.4f} {label}\n"

        contents += "@<TRIPOS>BOND\n"
        for pol_number, pol in enumerate(finished_polimers):
            addition = self._max_monomers_count * pol_number
            for i in range(pol.len() - 1):
                contents += f"{addition + i + 1} {addition + i + 1} {addition + i + 2} 1\n"

        file_name = f'{path_to_save_dir}/polymer_{file_number}_{self._accept_threshold}.mol2'
        with open(file_name, 'x') as f:
            f.write(contents)
            print(f"Файлы сохранены в {file_name}")



    def calc(self, path_to_save_dir: str = None) -> list[Polymer]:
        for epoch in range(self._globuls_count):
            print(f'epoch = {epoch}')
            finished_polimers = []
            field = Field(self._sphere_radius)
            polymers = [Polymer(field) for i in range(self._polymers_count)]
            polymers_cache = [[] for i in range(self._polymers_count)]
            for p in polymers:
                new_pos = field.define_start_position()
                p.add_monomer(new_pos)
                field.make_filled(new_pos)

            blacklist = []
            while len(finished_polimers) != len(polymers):
                for i in range(len(polymers)):
                    if i in blacklist:
                        continue
                    polymer = polymers[i]
                    polymer_cache = polymers_cache[i]
                    if polymer.len() == self._max_monomers_count:
                        continue
                    
                    current_position = polymer.back()
                    available_cells = field.get_available_cells(current_position)
                    if  len(available_cells) == 0:
                        finished_polimers.append(polymer)
                        blacklist.append(i)
                        continue
                        
                    continuations = self.__get_continuations(len(available_cells), available_cells)
                    potential_configs = [self.__get_next_config(polymer, next_step) for next_step in continuations]
                    U_current = polymer.calc_energy()
                    while True:
                        choice = np.random.randint(0, len(potential_configs))
                        next_config = potential_configs[choice]
                        deltaU = U_current - next_config.calc_energy()
                        if deltaU > 0:
                            current_position = next_config.back()
                            break
                        r = np.random.uniform(0.0, 1.0, 1)
                        if r < self._accept_threshold:
                            current_position = next_config.back()
                            break
                        
                    polymer.add_monomer(current_position)
                    polymer_cache.append(current_position)
                    field.make_filled(current_position)
                    persentage = polymer.len() / self._max_monomers_count * 100
                    int_persentage = int(persentage)
                    if (persentage - int_persentage < 0.1):
                        print(f'\t\t{polymer.name()}. Done: {int_persentage}%/100%')

                    print(f'{polymer.name()}\'s monomers count: {polymer.len()}')
                    if polymer.len() == self._max_monomers_count:
                        finished_polimers.append(polymer)
                        blacklist.append(i)

            if path_to_save_dir is not None:
                self.__save_polymer_mol_into_file(epoch, finished_polimers)
            return finished_polimers