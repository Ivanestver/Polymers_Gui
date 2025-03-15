from PyQt6.QtGui import QVector3D, QQuaternion, QColor
from PyQt6.Qt3DExtras import Qt3DWindow, QPhongMaterial, QOrbitCameraController, QCylinderMesh
from PyQt6.Qt3DCore import QEntity, QTransform
from alg.polymer_lib import Polymer
from space import Space
from polymer_view import GlobulaView

class Window3D(Qt3DWindow):
    
    def __init__(self):
        super().__init__()

        self.rootEntity = QEntity()
        self.material = QPhongMaterial(self.rootEntity)

        self.setRootEntity(self.rootEntity)
        self.globulas = list[GlobulaView]()
        self.__configure_camera()
        self.__configure_coordinate_axises()

    def __configure_camera(self):
        camera = self.camera()
        camera.setPosition(Space.global_zero)
        camera.setViewCenter(Space.global_zero - QVector3D(0, 0, Space.global_zero.z()))

        camController = QOrbitCameraController(self.rootEntity)
        camController.setLinearSpeed(100)
        camController.setLookSpeed(0)
        camController.setCamera(camera)

    def __configure_coordinate_axises(self):
        def setup_axis(rotation: QQuaternion, material: QPhongMaterial):
            mesh = QCylinderMesh()
            mesh.setRadius(0.02)
            mesh.setLength(1000)
        
            translation = QTransform()
            translation.setTranslation(Space.global_zero)
            translation.setRotation(rotation)

            entity = QEntity(self.rootEntity)
            entity.addComponent(mesh)
            entity.addComponent(translation)
            entity.addComponent(material)

        # set up the X axis
        xMaterial = QPhongMaterial()
        xMaterial.setAmbient(QColor.fromRgb(255, 0, 0))
        setup_axis(
            QQuaternion.fromAxisAndAngle(Space.forward_vector, 90.0),
            xMaterial
        )
        # set up the Y axis
        yMaterial = QPhongMaterial()
        yMaterial.setAmbient(QColor.fromRgb(0, 255, 0))
        setup_axis(
            QQuaternion.fromAxisAndAngle(Space.left_vector, 00.0),
            yMaterial
        )
        # set up the Z axis
        zMaterial = QPhongMaterial()
        zMaterial.setAmbient(QColor.fromRgb(0, 0, 255))
        setup_axis(
            QQuaternion.fromAxisAndAngle(Space.left_vector, 90.0),
            zMaterial
        )
    
    def clear_scene(self):
        nodes: list = self.rootEntity.childNodes()
        for i in range(len(nodes) - 1, 4, -1):
            if nodes[i] is QEntity:
                nodes.remove(nodes[i])

    def add_globula(self, polymers: list[Polymer], to_show: bool = True):
        new_globula = GlobulaView(f'Globula {len(self.globulas)}', polymers, self.rootEntity if to_show else None)
        self.globulas.append(new_globula)
        return new_globula

    def remove_globula(self, globula: GlobulaView):
        if globula not in self.globulas:
            return
        
        if globula.entity != None:
            globula.entity.setParent(None)
        self.globulas.remove(globula)
    
    def get_globula(self, i) -> GlobulaView:
        return self.globulas[i]
    
    def get_globulas_count(self) -> int:
        return len(self.globulas)