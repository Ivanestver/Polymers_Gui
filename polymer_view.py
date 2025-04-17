import math
from PyQt6.QtCore import QJsonValue
from PyQt6.QtGui import QVector3D, QQuaternion, QColor
from PyQt6.Qt3DCore import QEntity, QTransform
from PyQt6.Qt3DExtras import QPhongMaterial, QSphereMesh, QCylinderMesh
from PyQt6.Qt3DRender import QObjectPicker
from space import Space
from alg.polymer_lib import Polymer, Monomer, MonomerType
from cluster_analysis.cluster import intersect_clusters, Cluster, join_clusters
from alg.config import Axis, get_axis_color, Side
from alg.field_lib import Field

class Entity:
    def __init__(self, name: str, rootEntity: QEntity) -> None:
        self._name = name
        self.entity = None
        if rootEntity is not None:
            self.transform = QTransform()

            self.usual_material = QPhongMaterial()
            #self.material.setAmbient(QColor.fromRgb(int(random() * 255), int(random() * 255), int(random() * 255)))
            self.usual_material.setAmbient(QColor.fromRgb(0, 0, 0))

            self.Owise_material = QPhongMaterial()
            self.Owise_material.setAmbient(QColor.fromRgb(255, 0, 0))

            self.entity = QEntity(rootEntity)
            self.entity.addComponent(self.transform)
            self.entity.addComponent(self.usual_material)

    @property
    def name(self):
        return self._name
        
    def rotationX(self):
        return self.transform.rotationX()

    def rotationY(self):
        return self.transform.rotationY()

    def rotationZ(self):
        return self.transform.rotationZ()

    def setRotationX(self, angle: float):
        self.transform.setRotationX(angle)

    def setRotationY(self, angle: float):
        self.transform.setRotationY(angle)

    def setRotationZ(self, angle: float):
        self.transform.setRotationZ(angle)

class PolymerView(Entity):
    def __init__(self, polymer: Polymer, rootEntity: QEntity):
        super().__init__(polymer.name(), rootEntity)
        if rootEntity is not None:
            self.transform.setTranslation(QVector3D(0, 0, 0))
            self.objectPicker = QObjectPicker()
            self.entity.addComponent(self.objectPicker)

            for i in range(polymer.len() - 1):
                curr_monomer = polymer.get_monomer_by_idx(i)
                self.__add_monomer__(curr_monomer)
                next_monomer = polymer.get_monomer_by_idx(i + 1)
                self.__add_monomer__(next_monomer)
                self.__add_connection__(curr_monomer, next_monomer)
        
        self.polymer = polymer

    def len(self):
        return self.polymer.len()

    def __iter__(self):
        return iter(self.polymer)

    def __getitem__(self, i):
        return self.polymer.get_monomer_by_idx(i)

    def get_start_end_monomers(self):
        if self.polymer.len() == 0:
            raise Exception("The PolymerView cannot contain an empty polymer")
        
        return (self.polymer.get_monomer_by_idx(0), self.polymer.get_monomer_by_idx(self.polymer.len() - 1))

    def __add_monomer__(self, monomer: Monomer):
        sphereMesh = QSphereMesh()
        sphereMesh.setRadius(0.3)

        sphereTransform = QTransform()
        sphereTransform.setTranslation(
            QVector3D(
                monomer[0], 
                monomer[1], 
                monomer[2]
            ) - Space.space_center
        )

        sphereEntity = QEntity(self.entity)
        sphereEntity.addComponent(sphereMesh)
        sphereEntity.addComponent(sphereTransform)
        sphereEntity.addComponent(self.__get_material_by_monomer_type(monomer.type))
    
    def __add_connection__(self, monomer1, monomer2):
        cylinderMesh = QCylinderMesh()
        cylinderMesh.setRadius(0.1)
        cylinderMesh.setLength(1)
        
        cylinderTransform = QTransform()
        cylinderTransform.setTranslation(QVector3D(
            (monomer1[0] + monomer2[0]) / 2,
            (monomer1[1] + monomer2[1]) / 2,
            (monomer1[2] + monomer2[2]) / 2
            ) - Space.space_center)
        rotation = self.__rotate_connection__(monomer1, monomer2)
        if rotation is not None:
            cylinderTransform.setRotation(rotation)

        cylinderEntity = QEntity(self.entity)
        cylinderEntity.addComponent(cylinderMesh)
        cylinderEntity.addComponent(cylinderTransform)
        cylinderEntity.addComponent(self.usual_material)

    def __get_material_by_monomer_type(self, type: MonomerType):
        if type == MonomerType.Usual:
            return self.usual_material
        elif type == MonomerType.Owise:
            return self.Owise_material
        else:
            raise Exception(f"Wrong monomer type: {type.value}")

    def __rotate_connection__(self, monomer1, monomer2):
        if monomer1[0] == monomer2[0] and monomer1[2] == monomer2[2]:
            return None
        elif monomer1[0] == monomer2[0]:
            return QQuaternion.fromAxisAndAngle(Space.left_vector, 90.0)
        else:
            return QQuaternion.fromAxisAndAngle(Space.forward_vector, 90.0)
        
    def __monomer_to_json(self, monomer: Monomer):
        return QJsonValue({
            'coords': QJsonValue([c for c in monomer.coords]),
            'type': QJsonValue(monomer.type),
            'prev_monomer': QJsonValue([c for c in monomer.prev_monomer.coords] if monomer.prev_monomer is not None else []),
            'next_monomer': QJsonValue([c for c in monomer.next_monomer.coords] if monomer.next_monomer is not None else [])
        })

    def __polymer_to_json(self, polymer: Polymer):
        return QJsonValue({
            'number': polymer.number(),
            'monomers': QJsonValue([self.__monomer_to_json(mon) for mon in polymer])
        })

    def to_json(self):
        start_point, end_point = self.get_start_end_monomers()
        length = math.sqrt(sum([(e - s) ** 2 for s, e in zip(start_point, end_point)]))
        json_dict = {
            'name': {
                'value': self.name,
                'name': 'Название'
            },
            'monomers_count': {
                'value': self.len(),
                'name': "Количество мономеров"
            },
            'polymer': self.__polymer_to_json(self.polymer),
            'statistics': {
                'end-to-end length': {
                    'value': length,
                    'name': "Межконцевое расстояние"
                }
            }
        }
        return QJsonValue(json_dict)

    def from_json(json: QJsonValue, field: Field):
        pol = Polymer(field, json['polymer']['number'].toInt())
        monomers = json['polymer']['monomers'].toArray()
        for monomer in monomers:
            coords = monomer['coords'].toArray()
            coords = list([v.toInt() for v in coords])
            next_monomer = monomer['next_monomer'].toArray()
            next_monomer = list([v.toInt() for v in next_monomer])
            prev_monomer = monomer['prev_monomer'].toArray()
            prev_monomer = list([v.toInt() for v in prev_monomer])
            t = monomer['type'].toInt()
            m = field.get_monomer_by_coords(coords)
            m.type = t
            m.next_monomer = field.get_monomer_by_coords(next_monomer) if len(next_monomer) != 0 else None
            m.prev_monomer = field.get_monomer_by_coords(prev_monomer) if len(prev_monomer) != 0 else None
            pol._polymer.append(m)

        return PolymerView(pol, None)

class ClusterView:
    def __init__(self, globula, avg_tuple: tuple[int, Axis]):
        self._avg, self._axis = avg_tuple
        self._clusters = find_clusters(globula, self._axis, avg_tuple)

    @property
    def clusters(self):
        return self._clusters

    def colorize(self, reset=False):
        for cluster in self._clusters:
            cluster.set_type_of_monomers(MonomerType.Usual if reset else get_axis_color(self._axis))

    def __iter__(self):
        return iter(self._clusters)

    def len(self):
        return len(self._clusters)

    def sort(self, lamb, reverse=True):
        self._clusters = sorted(self._clusters, key=lambda c: lamb(c), reverse=reverse)

    def get_cluster_by_number(self, num: int):
        return self._clusters[num]

class GlobulaView(Entity):
    _finished_polimers_labels = ['C', 'N', 'H', 'O', 'F', 'Na', 'Mg', 'Al', 'P', 'S']
    
    def __init__(self, name: str, polymers: list[Polymer], rootEntity: QEntity) -> None:
        super().__init__(name, rootEntity)
        if rootEntity is not None:
            self.transform.setTranslation(Space.space_center)

        self.__polymers = [PolymerView(polymer, self.entity) for polymer in polymers]
        self.__x_clusters: ClusterView = None
        self.__y_clusters: ClusterView = None
        self.__z_clusters: ClusterView = None
        self.__common_cluster_done = False

    def __iter__(self):
        return iter(self.__polymers)

    def len(self):
        return len(self.__polymers)
    
    def to_json(self):
        json_dict = {
            'polymers_count': {
                'value': len(self.__polymers),
                'name': "Количество полимеров"
            },
            'polymers': {
                'value': QJsonValue([polymer.to_json() for polymer in self.__polymers]),
                'name': "Полимеры"
            }
        }

        return QJsonValue(json_dict)

    def from_json(json: QJsonValue, sphere_radius: int):
        polymers_json = json['value']['polymers']['value'].toArray()
        new_globula = GlobulaView(json['name'].toString(), list(), None)
        field = Field(sphere_radius)
        for polymer_json in polymers_json:
            new_globula.__polymers.append(PolymerView.from_json(polymer_json, field))
        return new_globula

    def get_mass_center(self):
        center = (0, 0, 0)
        N = 0
        for polymer in self.__polymers:
            for monomer in polymer:
                center = tuple(c + m for c, m in zip(center, monomer))
            N += polymer.len()
        
        return tuple((i / N for i in center)) if N > 0 else (0, 0, 0)

    def x_clusters(self, avg: int = 0):
        if self.__x_clusters is None:
            self.__x_clusters = ClusterView(self, (avg, Axis.X_AXIS))

        return self.__x_clusters

    def y_clusters(self, avg: int = 0):
        if self.__y_clusters is None:
            self.__y_clusters = ClusterView(self, (avg, Axis.Y_AXIS))

        return self.__y_clusters

    def z_clusters(self, avg: int = 0):
        if self.__z_clusters is None:
            self.__z_clusters = ClusterView(self, (avg, Axis.Z_AXIS))

        return self.__z_clusters

    def common_clusters(self):
        if not self.__common_cluster_done:
            clusters_X = self.x_clusters().clusters
            clusters_Y = self.y_clusters().clusters
            clusters_Z = self.z_clusters().clusters

            clusters_X_Y = list([intersect_clusters(cluster_X, cluster_Y) for cluster_X in clusters_X for cluster_Y in clusters_Y])
            clusters_X_Y = list([mon for cluster in clusters_X_Y for mon in cluster])

            clusters_X_Z = list([intersect_clusters(cluster_X, cluster_Z) for cluster_X in clusters_X for cluster_Z in clusters_Z])
            clusters_X_Z = list([mon for cluster in clusters_X_Z for mon in cluster])

            clusters_Y_Z = list([intersect_clusters(cluster_Y, cluster_Z) for cluster_Y in clusters_Y for cluster_Z in clusters_Z])
            clusters_Y_Z = list([mon for cluster in clusters_Y_Z for mon in cluster])

            for mon in clusters_X_Y:
                for cluster in clusters_X:
                    if cluster.remove_monomer(mon):
                        break
                for cluster in clusters_Y:
                    if cluster.remove_monomer(mon):
                        break

            for mon in clusters_X_Z:
                for cluster in clusters_X:
                    if cluster.remove_monomer(mon):
                        break
                for cluster in clusters_Z:
                    if cluster.remove_monomer(mon):
                        break

            for mon in clusters_Y_Z:
                for cluster in clusters_Y:
                    if cluster.remove_monomer(mon):
                        break
                for cluster in clusters_Z:
                    if cluster.remove_monomer(mon):
                        break

        return (self.__x_clusters, self.__y_clusters, self.__z_clusters)
    
    def biggest_common_clusters(self, desired_percentage):
        clusters_X, clusters_Y, clusters_Z = self.common_clusters()
        total_monomers_count = sum([pol.len() for pol in self.__polymers])

        clusters_X.sort(lambda c1: c1.size, True)
        clusters_Y.sort(lambda c1: c1.size, True)
        clusters_Z.sort(lambda c1: c1.size, True)

        chosen_clusters: list[ClusterView] = []
        
        def chosen_clusters_monomers_count():
            return sum([c.size for c in chosen_clusters]) if len(chosen_clusters) > 0 else 0

        def test(testing: ClusterView, first: ClusterView, second: ClusterView):
            if testing.len() == 0:
                return False
            
            if first.len() != 0 and second.len() != 0:
                return testing.get_cluster_by_number(0).size >= first.get_cluster_by_number(0).size and testing.get_cluster_by_number(0).size >= second.get_cluster_by_number(0).size
            elif first.len() != 0:
                return testing.get_cluster_by_number(0).size >= first.get_cluster_by_number(0).size
            elif second.len() != 0:
                return testing.get_cluster_by_number(0).size >= second.get_cluster_by_number(0).size
            else:
                return True
        
        percentage = 0
        while percentage <= desired_percentage:
            if test(clusters_X, clusters_Y, clusters_Z):
                chosen_clusters.append(clusters_X.get_cluster_by_number(0))
                clusters_X = clusters_X[1:]
            elif test(clusters_Y, clusters_X, clusters_Z):
                chosen_clusters.append(clusters_Y.get_cluster_by_number(0))
                clusters_Y = clusters_Y[1:]
            elif test(clusters_Z, clusters_X, clusters_Y):
                chosen_clusters.append(clusters_Z.get_cluster_by_number(0))
                clusters_Z = clusters_Z[1:]
            else:
                assert False, "It couldn't happen"

            percentage = chosen_clusters_monomers_count() / total_monomers_count * 100

        return chosen_clusters

    def reset(self):
        for pol in self.__polymers:
            for mon in pol:
                mon.type = MonomerType.Usual

        self.__x_clusters = None
        self.__y_clusters = None
        self.__z_clusters = None
        self.__common_cluster_done = False

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

def find_clusters(current_globula: GlobulaView, axis: Axis, avg_tuple: tuple[int, Axis] = None):
    clusters = list[Cluster]()
    def fill_clusters(directions, main_direction):
        pot_cluster: list[Monomer] = [monomer, monomer.sides[main_direction]]
        do_traverse(main_direction, directions, monomer, pot_cluster)
        if len(pot_cluster) == 8:
            clusters.append(Cluster(pot_cluster, main_direction, axis))
    for pol in current_globula:
        for monomer in pol:
            main_direction = get_direction(monomer)
            if main_direction == Side.Undefined or monomer.is_of_type(MonomerType.Undefined):
                continue

            old_size = len(clusters)
            if axis == Axis.X_AXIS and main_direction in [Side.Forward, Side.Backward]:
                fill_clusters([Side.Left, Side.Down, Side.Right], main_direction)
                fill_clusters([Side.Left, Side.Up, Side.Right], main_direction)
                fill_clusters([Side.Right, Side.Down, Side.Left], main_direction)
                fill_clusters([Side.Right, Side.Up, Side.Left], main_direction)
            elif axis == Axis.Y_AXIS and main_direction in [Side.Left, Side.Right]:
                fill_clusters([Side.Backward, Side.Down, Side.Forward], main_direction)
                fill_clusters([Side.Backward, Side.Up, Side.Forward], main_direction)
                fill_clusters([Side.Forward, Side.Down, Side.Backward], main_direction)
                fill_clusters([Side.Forward, Side.Up, Side.Backward], main_direction)
            elif axis == Axis.Z_AXIS and main_direction in [Side.Up, Side.Down]:
                fill_clusters([Side.Left, Side.Forward, Side.Right], main_direction)
                fill_clusters([Side.Left, Side.Backward, Side.Right], main_direction)
                fill_clusters([Side.Right, Side.Forward, Side.Left], main_direction)
                fill_clusters([Side.Right, Side.Backward, Side.Left], main_direction)
            else:
                continue

            if old_size != len(clusters):
                clusters = gather_clusters(clusters, avg_tuple, old_size)

    # Although we seem to have all the clusters already joined
    # it is possible that some clusters can be joined.
    # Since here we don't have a lot of them, this step shouldn't be
    # too time-consuming
    return gather_clusters(clusters, avg_tuple, None)

def get_monomers_count_in_clusters(clusters: list[Cluster]):
    s = set[Monomer]()
    for cluster in clusters:
        s.update(set(cluster.monomers))

    return len(s)

def gather_clusters(clusters: list[Cluster], avg_tuple : tuple[int, Axis] = None, start_from: int = None):
    if avg_tuple is None:
        return clusters

    avg, avg_axis = avg_tuple
    while True:
        clusters_new = list[Cluster]()
        # Build a dict where for each i we map a list of js that can be joined with i
        d = dict[int, list[int]]()
        # If start_from is None, then we need to compare all the cluster between each other
        # to do the deep check of clusters possible to join
        if start_from is None:
            for i in range(len(clusters) - 1):
                for j in range(i + 1, len(clusters)):
                    joined_cluster = join_clusters(clusters[i], clusters[j])
                    if joined_cluster is not None:
                        if i not in d.keys():
                            d[i] = list[int]()
                            
                        d[i].append(j)
        # The algorithm of highlighing clusters determine that all the new clusters are added
        # into the end of the clusters list, therefore, it is assumed it isn't able to join clusters
        # before start_from. So we need to check only the new ones
        else:
            for i in range(start_from, len(clusters)):
                for j in range(start_from):
                    joined_cluster = join_clusters(clusters[i], clusters[j])
                    if joined_cluster is not None:
                        if i not in d.keys():
                            d[i] = list[int]()
                            
                        d[i].append(j)

        # If there's nothing to join, stop the algorithm
        if len(d) == 0:
            break

        # Now join the extracted clusters recursively
        # that means that if we, e.g, joined clusters 1 and 2 into a 1-2 cluster
        # and cluster 2 and 3 into a 2-3 cluster, therefore, the joined ones have a common set of monomers
        # which are in the cluster 2 and can be joined as well into a 1-2-3 cluster
        def join_joined_clusters(current_node: int):
            joined_clusters_for_current_node = list[int]()
            # If nothing to join, come back
            if current_node not in d.keys() or len(d[current_node]) == 0:
                return joined_clusters_for_current_node
            
            for value in d[current_node]:
                # Add the current cluster as the connection
                joined_clusters_for_current_node.append(value)
                # Recursively find others ready to be joined
                received = join_joined_clusters(value)
                for rec_value in received:
                    if rec_value not in joined_clusters_for_current_node:
                        joined_clusters_for_current_node.append(rec_value)
                
            # All clusters for this cluster are to become a part of another list,
            # therefore, clear the current to avoid issues
            d[current_node].clear()
            return joined_clusters_for_current_node
        
        # Join the clusters. In fact, find the numbers belonging to the same cluster
        for key in d.keys():
            d[key] = join_joined_clusters(key) 
        
        # Do actual join
        joined_clusters_indexes = set()
        for key in d.keys():
            joined_clusters_indexes.add(key)
            joined_clusters_indexes.update(d[key])

            if len(d[key]) == 0:
                continue

            curr_cluster = clusters[key]
            for value in d[key]:
                c = join_clusters(curr_cluster, clusters[value])
                if c is None:
                    continue
                curr_cluster = c
            clusters_new.append(curr_cluster)

        # Those clusters not been touched must be moved as well
        indexes_not_joined = set([i for i in range(len(clusters)) if i not in joined_clusters_indexes])
        for i in indexes_not_joined:
            clusters_new.append(clusters[i])
        
        clusters = clusters_new

    # return those which satisfy the condition
    return list(set([c for c in clusters if c.get_avg_length_by_axis(avg_axis) >= avg]))