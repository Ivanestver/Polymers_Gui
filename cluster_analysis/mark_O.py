from cluster_analysis.function import Direction, get_direction
from alg.polymer_lib import Monomer, MonomerType
from polymer_view import GlobulaView
from cluster_analysis.cluster import Cluster

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

def mark_clusters(current_globula: GlobulaView, avg: float):
    clusters = list[Cluster]()
    for pol in current_globula:
        for monomer in pol:
            if monomer.is_of_type(MonomerType.Owise):
                cluster = Cluster(monomer)
                cluster.make()
                if not cluster.fits_constraint_avg(avg):
                    cluster.destroy_cluster()
                else:
                    clusters.append(cluster)

    for cluster in clusters:
        cluster.set_cluster_type(MonomerType.Owise)