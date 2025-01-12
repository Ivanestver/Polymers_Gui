from PyQt6.Qt3DExtras import Qt3DWindow, QPhongMaterial
from PyQt6.Qt3DCore import QEntity
from alg.polymer_lib import Polymer

class Window3D(Qt3DWindow):
    def __init__(self):
        super().__init__()

        self.rootEntity = QEntity()
        self.material = QPhongMaterial(self.rootEntity)

        self.setRootEntity(self.rootEntity)
    
    def addPolymer(self, polymer: Polymer):
        pass
