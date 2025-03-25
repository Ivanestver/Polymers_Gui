from cluster_analysis.function import Direction, get_direction
from alg.polymer_lib import Monomer, MonomerType
from polymer_view import GlobulaView
from alg.config import Side

def mark_O_part(current_globula: GlobulaView):
    min_chain_length = 3
    for polymer in current_globula:
        previous_monomer = polymer[0]
        first_monomer_in_row = polymer[0]
        current_monomer = polymer[1]
        current_direction = get_direction(previous_monomer, current_monomer)
        previous_monomer = current_monomer
        monomers_in_row = 2
        for i in range(2, polymer.len()):
            current_monomer = polymer[i]
            direction = get_direction(previous_monomer, current_monomer)
            if direction == current_direction:
                monomers_in_row += 1
            else:
                if monomers_in_row >= min_chain_length: # Сейчас считаем размер таким
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
        
        # Осталось обработать последний в цепочке полимер
        if monomers_in_row >= min_chain_length: # Сейчас считаем размер таким
            # Помечаем полученную цепочку кислородом
            m: Monomer = first_monomer_in_row
            while m != previous_monomer:
                m.type = MonomerType.Owise
                m = m.next_monomer
            m.type = MonomerType.Owise

def get_direction(monomer: Monomer):
    next_monomer = monomer.next_monomer
    if next_monomer is None:
        return Side.Undefined
    
    for side, side_monomer in monomer.sides.items():
        if side_monomer is None:
            continue
        if side_monomer == next_monomer:
            return side

    return Side.Undefined

def do_traverse(main_direction: Side, directions: list[Side], curr_monomer: Monomer, pot_cluster: list[Monomer]):
    if len(directions) == 0:
        return
    
    next_monomer: Monomer = curr_monomer.sides[directions[0]]
    if next_monomer is None or next_monomer.is_of_type(MonomerType.Undefined):
        pot_cluster.clear()
        return
    
    next_in_main_direction_monomer: Monomer = next_monomer.sides[main_direction]
    if next_in_main_direction_monomer is None:
        pot_cluster.clear()
        return

    next_next_monomer = next_monomer.next_monomer
    prev_next_monomer = next_monomer.prev_monomer
    if next_next_monomer is None and prev_next_monomer is None:
        pot_cluster.clear()
        return

    if next_next_monomer == next_in_main_direction_monomer or prev_next_monomer == next_in_main_direction_monomer:
        pot_cluster.append(next_monomer)
        pot_cluster.append(next_in_main_direction_monomer)
        do_traverse(main_direction, directions[1:], next_monomer, pot_cluster)
    else:
        pot_cluster.clear()

def mark_clusters(current_globula: GlobulaView, avg: float, take_horizontal: bool = False):
    clusters = set[Monomer]()
    def fill_clusters(directions):
        pot_cluster: list[Monomer] = [monomer, monomer.sides[main_direction]]
        do_traverse(main_direction, directions, monomer, pot_cluster)
        if len(pot_cluster) == 8:
            clusters.update(pot_cluster)
    for pol in current_globula:
        for monomer in pol:
            main_direction = get_direction(monomer)
            if main_direction == Side.Undefined:
                continue
            if take_horizontal:
                if main_direction in [Side.Up, Side.Down]:
                    continue

                if main_direction in [Side.Left, Side.Right]:
                    fill_clusters([Side.Backward, Side.Down, Side.Forward])
                    fill_clusters([Side.Backward, Side.Up, Side.Forward])
                    fill_clusters([Side.Forward, Side.Down, Side.Backward])
                    fill_clusters([Side.Forward, Side.Up, Side.Backward])
                else:
                    fill_clusters([Side.Left, Side.Down, Side.Right])
                    fill_clusters([Side.Left, Side.Up, Side.Right])
                    fill_clusters([Side.Right, Side.Down, Side.Left])
                    fill_clusters([Side.Right, Side.Up, Side.Left])
            else:
                if main_direction not in [Side.Up, Side.Down]:
                    continue

                fill_clusters([Side.Left, Side.Forward, Side.Right])
                fill_clusters([Side.Left, Side.Backward, Side.Right])
                fill_clusters([Side.Right, Side.Forward, Side.Left])
                fill_clusters([Side.Right, Side.Backward, Side.Left])

    for monomer in clusters:
        monomer.type = MonomerType.Owise if take_horizontal else MonomerType.Nwise
