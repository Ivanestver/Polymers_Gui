from uis.ui_radius_to_cut import Ui_DlgRadiusToCutWindow
from PyQt6.QtWidgets import QDialog

class DlgRadiusToCutWindow(QDialog):
    def __init__(self, max_radius, *args):
        super().__init__(*args)

        self.ui = Ui_DlgRadiusToCutWindow()
        self.ui.setupUi(self)

        self.max_radius = max_radius
        self.__set_msg()
        self.__set_spin_box()

    def __set_msg(self):
        self.ui.msgLabel.setText(f"Желаемый радиус (максимально: {self.max_radius})")

    def __set_spin_box(self):
        self.ui.radiusSpinBox.setValue(self.max_radius)
        self.ui.radiusSpinBox.setMaximum(self.max_radius)

    def get_desired_radius(self):
        return self.ui.radiusSpinBox.value()