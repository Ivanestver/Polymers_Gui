from uis.ui_choose_axis import Ui_DlgChooseAxis
from PyQt6.QtWidgets import QDialog
from alg.config import Axis

class DlgChooseAxis(QDialog):
    def __init__(self, current_axis: Axis, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.ui = Ui_DlgChooseAxis()
        self.ui.setupUi(self)
        self.axis_translation = {
            Axis.X_AXIS: "Ось X",
            Axis.Y_AXIS: "Ось Y",
            Axis.Z_AXIS: "Ось Z"
        }
        self.axis_permitted = {
            Axis.X_AXIS: [Axis.Y_AXIS, Axis.Z_AXIS],
            Axis.Y_AXIS: [Axis.X_AXIS, Axis.Z_AXIS],
            Axis.Z_AXIS: [Axis.X_AXIS, Axis.Y_AXIS]
        }
        self.current_axis = current_axis

        l = list([self.axis_translation[axis] for axis in self.axis_permitted[self.current_axis]])
        self.ui.cmbAxis.addItems(l)

    def get_chosen_axis(self):
        return self.axis_permitted[self.current_axis][self.ui.cmbAxis.currentIndex()]