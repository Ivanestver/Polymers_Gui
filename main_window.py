from PyQt6.QtWidgets import QDialog, QWidget
from PyQt6.QtGui import QVector3D
from uis.ui_main_window import Ui_MainWindow
import PyQt6.Qt3DCore as core3d
from PyQt6.Qt3DExtras import QSphereMesh, QPhongMaterial, QOrbitCameraController
import PyQt6.QtCore as core
from window_3d import Window3D

class MainWindow(QDialog):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.ui = Ui_MainWindow()
        self.ui.setupUi(self)

        self.view = Window3D()
        container = QWidget.createWindowContainer(self.view)
        screenSize = self.view.screen().size()
        container.setMinimumSize(core.QSize(200, 100))
        container.setMaximumSize(screenSize)

        self.ui.horizontalLayout.addWidget(container)

        self.rootEntity = core3d.QEntity()
        material = QPhongMaterial(self.rootEntity)

        sphereMesh = QSphereMesh()
        sphereMesh.setRadius(3)
        sphereMesh.setGenerateTangents(True)

        sphereTransform = core3d.QTransform()
        sphereTransform.setTranslation(QVector3D(0, 0, 0))

        sphereEntity = core3d.QEntity(self.rootEntity)
        sphereEntity.addComponent(sphereMesh)
        sphereEntity.addComponent(sphereTransform)
        sphereEntity.addComponent(material)

        self.view.setRootEntity(self.rootEntity)
        camera = self.view.camera()
        #camera.lens().setPerspectiveProjection(45.0, 16.0/9.0, 0.1, 1000.0)
        camera.setPosition(QVector3D(0, 0, 40))
        camera.setViewCenter(QVector3D(0, 0, 0))

        camController = QOrbitCameraController(self.rootEntity)
        camController.setLinearSpeed(100)
        camController.setLookSpeed(360)
        camController.setCamera(camera)