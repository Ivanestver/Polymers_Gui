from PyQt6.QtGui import QVector3D
from PyQt6.Qt3DExtras import Qt3DWindow, QPhongMaterial, QOrbitCameraController, QSphereMesh
from PyQt6.Qt3DCore import QEntity, QTransform
from alg.polymer_lib import Polymer

class Window3D(Qt3DWindow):
    def __init__(self):
        super().__init__()

        self.rootEntity = QEntity()
        self.material = QPhongMaterial(self.rootEntity)

        self.setRootEntity(self.rootEntity)
        self.__configure_camera()

    def __configure_camera(self):
        camera = self.camera()
        camera.setPosition(QVector3D(0, 0, 30))
        camera.setViewCenter(QVector3D(0, 0, 0))

        camController = QOrbitCameraController(self.rootEntity)
        camController.setLinearSpeed(100)
        camController.setLookSpeed(360)
        camController.setCamera(camera)
    
    def add_polymer(self, polymer: Polymer):
        for monomer in polymer:
            sphereMesh = QSphereMesh()
            sphereMesh.setRadius(0.3)

            sphereTransform = QTransform()
            sphereTransform.setTranslation(QVector3D(monomer[0] - 100, monomer[1] - 100, monomer[2] - 100))

            sphereEntity = QEntity(self.rootEntity)
            sphereEntity.addComponent(sphereMesh)
            sphereEntity.addComponent(sphereTransform)
            sphereEntity.addComponent(self.material)
