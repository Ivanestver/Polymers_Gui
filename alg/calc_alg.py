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

    def __run_alg(self, polymers: list[Polymer], field: Field):
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

    def calc(self) -> list[Polymer]:
        field = Field(self._sphere_radius)
        polymers = [Polymer(field) for i in range(self._polymers_count)]
        finished_polimers = self.__run_alg(polymers, field)
        return finished_polimers
    
    def build_more(self, globula: list[Polymer]):
        if len(globula) == 0:
            return None

        polymers = list[Polymer]()
        field = globula[0].field()
        if field.is_busy():
            return None

        for i in range(self._polymers_count):
            polymers.append(Polymer(field))
        
        built_polimers = self.__run_alg(polymers, field)
        return globula + built_polimers