from PyQt6.QtWidgets import QDialog, QWidget
from uis.ui_main_window import Ui_MainWindow
import PyQt6.QtCore as core
from window_3d import Window3D
from alg.polymers_copy import CalcAlg
from space import Space
from PyQt6.QtGui import QVector3D

class MainWindow(QDialog):
    def __init__(self, space_dimention, *args, **kwargs):
        super().__init__(*args, **kwargs)

        Space.space_dimention = space_dimention
        Space.global_zero = QVector3D(space_dimention / 2, space_dimention / 2, space_dimention / 2)

        self.ui = Ui_MainWindow()
        self.ui.setupUi(self)

        self.view = Window3D()
        container = QWidget.createWindowContainer(self.view)
        screenSize = self.view.screen().size()
        container.setMinimumSize(core.QSize(200, 100))
        container.setMaximumSize(screenSize)

        self.ui.mainLayout.addWidget(container)
        self.ui.calculateBtn.clicked.connect(self.on_calc_btn_clicked)

        self.ui.btnUp.clicked.connect(self.on_btn_up_clicked)
        self.ui.btnLeft.clicked.connect(self.on_btn_left_clicked)
        self.ui.btnDown.clicked.connect(self.on_btn_down_clicked)
        self.ui.btnRight.clicked.connect(self.on_btn_right_clicked)

    def on_calc_btn_clicked(self):
        globuls_count = self.ui.filesCountSpinBox.value()
        polymers_count = self.ui.polymersCountSpinBox.value()
        accept_threshold = self.ui.thresholdSpinBox.value()
        monomers_count = self.ui.monomersCountSpinBox.value()
        calcAlg = CalcAlg(globuls_count, polymers_count, accept_threshold, monomers_count)
        finished_polymers = calcAlg.calc()
        for polymer in finished_polymers:
            self.view.add_polymer(polymer)

    def on_btn_up_clicked(self):
        polymer = self.view.get_polymer(0)
        polymer.setRotationX(polymer.rotationX() + 10)

    def on_btn_down_clicked(self):
        polymer = self.view.get_polymer(0)
        polymer.setRotationX(polymer.rotationX() - 10)

    def on_btn_left_clicked(self):
        polymer = self.view.get_polymer(0)
        polymer.setRotationY(polymer.rotationY() + 10)

    def on_btn_right_clicked(self):
        polymer = self.view.get_polymer(0)
        polymer.setRotationY(polymer.rotationY() - 10)