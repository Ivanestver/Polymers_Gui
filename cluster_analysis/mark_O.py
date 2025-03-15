from cluster_analysis.function import Direction, get_direction
from alg.polymer_lib import Monomer, MonomerType
from polymer_view import GlobulaView

def mark_O_part(current_globula: GlobulaView):
    for polymer in current_globula:
        previous_monomer = polymer[0]
        first_monomer_in_row = polymer[0]
        current_direction = Direction.Forward
        monomers_in_row = 1
        for i in range(1, polymer.len()):
            current_monomer = polymer[i]
            direction = get_direction(previous_monomer, current_monomer)
            if direction == current_direction:
                monomers_in_row += 1
            else:
                if monomers_in_row >= 6: # Сейчас считаем размер таким
                    # Помечаем полученную цепочку кислородом
                    m: Monomer = first_monomer_in_row
                    while m != previous_monomer:
                        m.type = MonomerType.Owise
                        m = m.next_monomer
                    m.type = MonomerType.Owise

                monomers_in_row = 2
                current_direction = direction
                first_monomer_in_row = previous_monomer
            previous_monomer = current_monomer