from PyQt6.QtWidgets import QDialog, QWidget
from uis.ui_main_window import Ui_MainWindow
import PyQt6.QtCore as core
from window_3d import Window3D
from alg.calc_alg import CalcAlg
from space import Space
from PyQt6.QtGui import QVector3D, QColor
from PyQt6.Qt3DRender import QPickEvent
from PyQt6.Qt3DExtras import QPhongMaterial

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

        self.ui.sphereRadiusText.setText(f"Радиус сферы (не больше {Space.space_dimention})")
        self.ui.radiusSphereSpinBox.setMaximum(Space.space_dimention)

    def on_calc_btn_clicked(self):
        globuls_count = self.ui.filesCountSpinBox.value()
        polymers_count = self.ui.polymersCountSpinBox.value()
        accept_threshold = self.ui.thresholdSpinBox.value()
        monomers_count = self.ui.monomersCountSpinBox.value()
        sphere_radius = self.ui.radiusSphereSpinBox.value()
        calcAlg = CalcAlg(globuls_count, polymers_count, accept_threshold, monomers_count, sphere_radius)
        finished_polymers = calcAlg.calc()
        new_globula = self.view.add_globula(finished_polymers, self.on_picker_clicked)
        self.ui.polymersListWidget.addItem(new_globula.name)
        # for polymer in finished_polymers:
        #     self.view.add_polymer(polymer, self.on_picker_clicked)
        #     self.ui.polymersListWidget.addItem(QListWidgetItem(polymer.name()))

    def on_picker_clicked(self, e: QPickEvent):
        material = QPhongMaterial()
        material.setAmbient(QColor.fromRgb(0, 0, 0))
        e.entity().addComponent(material)
        pass

    def on_btn_up_clicked(self):
        curr_polymer_number = self.ui.polymersListWidget.currentRow()
        if curr_polymer_number < 0 or curr_polymer_number >= self.view.get_globulas_count():
            return

        polymer = self.view.get_globula(curr_polymer_number)
        polymer.setRotationX(polymer.rotationX() + 10)

    def on_btn_down_clicked(self):
        curr_polymer_number = self.ui.polymersListWidget.currentRow()
        if curr_polymer_number < 0 or curr_polymer_number >= self.view.get_globulas_count():
            return

        polymer = self.view.get_globula(curr_polymer_number)
        polymer.setRotationX(polymer.rotationX() - 10)

    def on_btn_left_clicked(self):
        curr_polymer_number = self.ui.polymersListWidget.currentRow()
        if curr_polymer_number < 0 or curr_polymer_number >= self.view.get_globulas_count():
            return

        polymer = self.view.get_globula(curr_polymer_number)
        polymer.setRotationY(polymer.rotationY() + 10)

    def on_btn_right_clicked(self):
        curr_polymer_number = self.ui.polymersListWidget.currentRow()
        if curr_polymer_number < 0 or curr_polymer_number >= self.view.get_globulas_count():
            return

        polymer = self.view.get_globula(curr_polymer_number)
        polymer.setRotationY(polymer.rotationY() - 10)