import os
from random import random
from PyQt6.QtGui import QVector3D, QQuaternion, QColor
from PyQt6.Qt3DCore import QEntity, QTransform
from PyQt6.Qt3DExtras import QPhongMaterial, QSphereMesh, QCylinderMesh
from PyQt6.Qt3DRender import QObjectPicker
from alg.config import Axis
from space import Space
from alg.polymer_lib import Polymer

class Entity:
    def __init__(self, name: str, rootEntity: QEntity) -> None:
        self.transform = QTransform()

        self.material = QPhongMaterial()
        self.material.setAmbient(QColor.fromRgb(int(random() * 255), int(random() * 255), int(random() * 255)))
        self.name = name

        self.entity = QEntity(rootEntity)
        self.entity.addComponent(self.transform)
        self.entity.addComponent(self.material)
        
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
        self.transform.setTranslation(QVector3D(0, 0, 0))

        for i in range(polymer.len() - 1):
            curr_monomer = polymer.get_monomer_by_idx(i)
            self.__add_monomer__(curr_monomer)
            next_monomer = polymer.get_monomer_by_idx(i + 1)
            self.__add_monomer__(next_monomer)
            self.__add_connection__(curr_monomer, next_monomer)
        
        self.objectPicker = QObjectPicker()
        self.entity.addComponent(self.objectPicker)
        self.polymer = polymer

    def len(self):
        return self.polymer.len()

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

    def turn_to_mol(self, element_name: str, number: int):
        min_width, max_width, min_height, max_height = self.polymer.get_min_max_width_height()
        if min_width == max_width or min_height == max_height:
            return
        contents = ""
        for i, monomer in enumerate(self.polymer):
            label = element_name
            if i == 0:
                label = 'Fe'
            if i == self.polymer.len() - 1:
                label = 'Cu'
            x: float = monomer[Axis.X_AXIS.value]
            y: float = monomer[Axis.Y_AXIS.value]
            z: float = monomer[Axis.Z_AXIS.value]
            contents += f"{self.len() * number + i + 1} {label} {x:.4f} {y:.4f} {z:.4f} {label}\n"

        return contents

class GlobulaView(Entity):
    _finished_polimers_labels = ['C', 'N', 'H', 'O', 'F', 'Na', 'Mg', 'Al', 'P', 'S']
    
    def __init__(self, name: str, polymers: list[Polymer], rootEntity: QEntity) -> None:
        super().__init__(name, rootEntity)
        self.transform.setTranslation(Space.global_zero)

        self.polymers = [PolymerView(polymer, self.entity) for polymer in polymers]
    
    def turn_to_mol(self):

        contents = ""
        contents += "@<TRIPOS>MOLECULE\n*****\n"
        contents += f" {sum([pl.len() for pl in self.polymers])} {sum([pl.len() for pl in self.polymers]) - len(self.polymers)} 0 0 0\n"
        contents += "SMALL\n"
        contents += "GASTEIGER\n\n\r\n"

        contents += "@<TRIPOS>ATOM\n"
        for pol_number, pol in enumerate(self.polymers):
            contents += pol.turn_to_mol(self._finished_polimers_labels[pol_number], pol_number)

        contents += "@<TRIPOS>BOND\n"
        addition = 0
        for pol_number, pol in enumerate(self.polymers):
            for i in range(pol.len() - 1):
                contents += f"{addition + i + 1} {addition + i + 1} {addition + i + 2} 1\n"
            addition += pol.len()

        return f'{self.name}.mol2', contents