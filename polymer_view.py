import math
import os
from random import random
from PyQt6.QtCore import QJsonValue, QJsonDocument
from PyQt6.QtGui import QVector3D, QQuaternion, QColor
from PyQt6.Qt3DCore import QEntity, QTransform
from PyQt6.Qt3DExtras import QPhongMaterial, QSphereMesh, QCylinderMesh
from PyQt6.Qt3DRender import QObjectPicker
from alg.config import Axis
from space import Space
from alg.polymer_lib import Polymer

class Entity:
    def __init__(self, name: str, rootEntity: QEntity) -> None:
        self._name = name
        self.entity = None
        if rootEntity is not None:
            self.transform = QTransform()

            self.material = QPhongMaterial()
            self.material.setAmbient(QColor.fromRgb(int(random() * 255), int(random() * 255), int(random() * 255)))

            self.entity = QEntity(rootEntity)
            self.entity.addComponent(self.transform)
            self.entity.addComponent(self.material)

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

    def get_start_end_monomers(self):
        if self.polymer.len() == 0:
            raise Exception("The PolymerView cannot contain an empty polymer")
        
        return (self.polymer.get_monomer_by_idx(0), self.polymer.get_monomer_by_idx(self.polymer.len() - 1))

    def __add_monomer__(self, monomer):
        sphereMesh = QSphereMesh()
        sphereMesh.setRadius(0.3)

        sphereTransform = QTransform()
        sphereTransform.setTranslation(
            QVector3D(
                monomer[0], 
                monomer[1], 
                monomer[2]
            ) - Space.global_zero
        )

        sphereEntity = QEntity(self.entity)
        sphereEntity.addComponent(sphereMesh)
        sphereEntity.addComponent(sphereTransform)
        sphereEntity.addComponent(self.material)
    
    def __add_connection__(self, monomer1, monomer2):
        cylinderMesh = QCylinderMesh()
        cylinderMesh.setRadius(0.1)
        cylinderMesh.setLength(1)
        
        cylinderTransform = QTransform()
        cylinderTransform.setTranslation(QVector3D(
            (monomer1[0] + monomer2[0]) / 2,
            (monomer1[1] + monomer2[1]) / 2,
            (monomer1[2] + monomer2[2]) / 2
            ) - Space.global_zero)
        rotation = self.__rotate_connection__(monomer1, monomer2)
        if rotation is not None:
            cylinderTransform.setRotation(rotation)

        cylinderEntity = QEntity(self.entity)
        cylinderEntity.addComponent(cylinderMesh)
        cylinderEntity.addComponent(cylinderTransform)
        cylinderEntity.addComponent(self.material)

    def __rotate_connection__(self, monomer1, monomer2):
        if monomer1[0] == monomer2[0] and monomer1[2] == monomer2[2]:
            return None
        elif monomer1[0] == monomer2[0]:
            return QQuaternion.fromAxisAndAngle(Space.left_vector, 90.0)
        else:
            return QQuaternion.fromAxisAndAngle(Space.forward_vector, 90.0)

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
            'statistics': {
                'end-to-end length': {
                    'value': length,
                    'name': "Межконцевое расстояние"
                }
            }
        }
        return QJsonValue(json_dict)

class GlobulaView(Entity):
    _finished_polimers_labels = ['C', 'N', 'H', 'O', 'F', 'Na', 'Mg', 'Al', 'P', 'S']
    
    def __init__(self, name: str, polymers: list[Polymer], rootEntity: QEntity) -> None:
        super().__init__(name, rootEntity)
        if rootEntity is not None:
            self.transform.setTranslation(Space.global_zero)

        self.polymers = [PolymerView(polymer, self.entity) for polymer in polymers]

    def __iter__(self):
        return iter(self.polymers)

    def len(self):
        return len(self.polymers)
    
    def to_json(self):
        json_dict = {
            'polymers_count': {
                'value': len(self.polymers),
                'name': "Количество полимеров"
            },
            'polymers': {
                'value': QJsonValue([polymer.to_json() for polymer in self.polymers]),
                'name': "Полимеры"
            }
        }

        return QJsonValue(json_dict)
    
    def get_mass_center(self):
        center = (0, 0, 0)
        N = 0
        for polymer in self.polymers:
            for monomer in polymer:
                center = tuple(c + m for c, m in zip(center, monomer))
            N += polymer.len()
        
        return tuple((i / N for i in center)) if N > 0 else (0, 0, 0)