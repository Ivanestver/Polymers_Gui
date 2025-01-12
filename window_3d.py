from PyQt6.QtGui import QVector3D
from PyQt6.Qt3DExtras import Qt3DWindow, QPhongMaterial, QOrbitCameraController
from PyQt6.Qt3DCore import QEntity
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
        camera.setPosition(QVector3D(0, 0, 40))
        camera.setViewCenter(QVector3D(0, 0, 0))

        camController = QOrbitCameraController(self.rootEntity)
        camController.setLinearSpeed(100)
        camController.setLookSpeed(360)
        camController.setCamera(camera)
    
    def addPolymer(self, polymer: Polymer):
        pass
