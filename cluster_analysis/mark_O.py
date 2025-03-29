from alg.polymer_lib import Monomer, MonomerType
from polymer_view import GlobulaView
from alg.config import Side, Axis

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

def mark_clusters(current_globula: GlobulaView, avg: float, axis: Axis):
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

            if axis == Axis.X_AXIS and main_direction in [Side.Forward, Side.Backward]:
                fill_clusters([Side.Left, Side.Down, Side.Right])
                fill_clusters([Side.Left, Side.Up, Side.Right])
                fill_clusters([Side.Right, Side.Down, Side.Left])
                fill_clusters([Side.Right, Side.Up, Side.Left])
            elif axis == Axis.Y_AXIS and main_direction in [Side.Left, Side.Right]:
                fill_clusters([Side.Backward, Side.Down, Side.Forward])
                fill_clusters([Side.Backward, Side.Up, Side.Forward])
                fill_clusters([Side.Forward, Side.Down, Side.Backward])
                fill_clusters([Side.Forward, Side.Up, Side.Backward])
            elif axis == Axis.Z_AXIS and main_direction in [Side.Up, Side.Down]:
                fill_clusters([Side.Left, Side.Forward, Side.Right])
                fill_clusters([Side.Left, Side.Backward, Side.Right])
                fill_clusters([Side.Right, Side.Forward, Side.Left])
                fill_clusters([Side.Right, Side.Backward, Side.Left])
            else:
                continue
            
    axis_type_dict = {
        Axis.X_AXIS: MonomerType.Owise,
        Axis.Y_AXIS: MonomerType.Nwise,
        Axis.Z_AXIS: MonomerType.Fwise
    }
    for monomer in clusters:
        monomer.type = axis_type_dict[axis]
