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

    def __get_next_current_position(self, potential_configs: list[Polymer], U_current: float):
        while True:
            choice = np.random.randint(0, len(potential_configs))
            next_config = potential_configs[choice]
            deltaU = U_current - next_config.calc_energy()
            if deltaU > 0:
                return next_config.back()
            r = np.random.uniform(0.0, 1.0, 1)
            if r < self._accept_threshold:
                return next_config.back()

    def __get_next_current_position_cristall(self, potential_configs: list[Polymer], U_current: float):
        deltas_of_potential_configs = np.array([U_current - pol.calc_energy() for pol in potential_configs])
        if np.any(deltas_of_potential_configs > 0):
            max_values_indices = np.argwhere(deltas_of_potential_configs == np.amax(deltas_of_potential_configs)).flatten()
            choice = np.random.randint(0, len(max_values_indices))
            return potential_configs[max_values_indices[choice]].back()
        else:
            while True:
                choice = np.random.randint(0, len(potential_configs))
                r = np.random.uniform(0.0, 1.0, 1)
                if r < self._accept_threshold:
                    return potential_configs[choice].back()

    def __run_alg_simultaneously(self, polymers: list[Polymer], field: Field):
        finished_polimers = []
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
                current_position = self.__get_next_current_position(potential_configs, polymer.calc_energy())
                    
                polymer.add_monomer(current_position)
                field.make_filled(current_position)
                persentage = polymer.len() / self._max_monomers_count * 100
                int_persentage = int(persentage)
                if (persentage - int_persentage < 0.1):
                    print(f'\t\t{polymer.name()}. Done: {int_persentage}%/100%')

                print(f'{polymer.name()}\'s monomers count: {polymer.len()}')
                if polymer.len() == self._max_monomers_count:
                    finished_polimers.append(polymer)
                    blacklist.append(i)

        return finished_polimers

    def calc_simultaneously(self) -> list[Polymer]:
        field = Field(self._sphere_radius)
        polymers = [Polymer(field) for i in range(self._polymers_count)]
        finished_polimers = self.__run_alg_simultaneously(polymers, field)
        return finished_polimers
    
    def build_more_simultaneously(self, globula: list[Polymer]):
        if len(globula) == 0:
            return None

        polymers = list[Polymer]()
        field = globula[0].field()
        if field.is_busy():
            return None

        for i in range(self._polymers_count):
            polymers.append(Polymer(field))
        
        built_polimers = self.__run_alg_simultaneously(polymers, field)
        return globula + built_polimers

    def calc_as_cristall(self):
        is_stepping_back = False
        finished_polimers = []
        field = Field(self._sphere_radius)
        while len(finished_polimers) != self._polymers_count:
            print(f'\tfinished_polymers count = {len(finished_polimers)}')
            current_position = field.define_start_position()
            field.make_filled(current_position)
            polymer = Polymer(field)
            polymer.add_monomer(current_position)
            polymer_cache = []
            while polymer.len() != self._max_monomers_count:
                available_cells = field.get_available_cells(current_position)
                if len(available_cells) == 0:
                    if not is_stepping_back:
                        is_stepping_back = True
                    if polymer.make_step_back():
                        current_position = polymer.back()
                        continue
                    else:
                        break
                    
                if is_stepping_back:
                    is_stepping_back = False

                continuations = self.__get_continuations(len(available_cells), available_cells)
                potential_configs = [self.__get_next_config(polymer, next_step) for next_step in continuations]
                current_position = self.__get_next_current_position_cristall(potential_configs, polymer.calc_energy())

                polymer.add_monomer(current_position)
                polymer_cache.append(current_position)
                field.make_filled(current_position)
                persentage = polymer.len() / self._max_monomers_count * 100
                int_persentage = int(persentage)
                if (persentage - int_persentage < 0.1):
                    print(f'\t\tDone: {int_persentage}%/100%')
            print(f'Monomers count: {polymer.len()}')
            if polymer.len() == self._max_monomers_count:
                finished_polimers.append(polymer)
            else:
                print(f'It is less than required. Roll back and repeat...')
                for cell in polymer_cache:
                    field.make_free(cell)
        
        return finished_polimers