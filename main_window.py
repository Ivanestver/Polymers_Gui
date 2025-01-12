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

        self.ui.mainLayout.addWidget(container)
        self.ui.calculateBtn.clicked.connect(self.on_calc_btn_clicked)

    def on_calc_btn_clicked(self):
        pass