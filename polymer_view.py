from random import random
from PyQt6.QtGui import QVector3D, QQuaternion, QColor
from PyQt6.Qt3DCore import QEntity, QTransform
from PyQt6.Qt3DExtras import QPhongMaterial, QSphereMesh, QCylinderMesh
from PyQt6.Qt3DRender import QObjectPicker
from space import Space
from alg.polymer_lib import Polymer

class PolymerView:
    def __init__(self, polymer: Polymer, rootEntity: QEntity):
        self.transform = QTransform()
        self.material = QPhongMaterial()
        self.material.setAmbient(QColor.fromRgb(int(random() * 255), int(random() * 255), int(random() * 255)))
        self.transform.setTranslation(Space.global_zero)
        self.name = polymer.name()

        self.entity = QEntity(rootEntity)
        self.entity.addComponent(self.transform)
        for i in range(polymer.len() - 1):
            curr_monomer = polymer.get_monomer_by_idx(i)
            self.__add_monomer__(curr_monomer)
            next_monomer = polymer.get_monomer_by_idx(i + 1)
            self.__add_monomer__(next_monomer)
            self.__add_connection__(curr_monomer, next_monomer)
        
        self.entity.addComponent(self.material)
        self.objectPicker = QObjectPicker()
        self.entity.addComponent(self.objectPicker)

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