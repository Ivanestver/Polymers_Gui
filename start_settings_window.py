from uis.ui_start_settings_window import Ui_StartSettingsWindow
from PyQt6.QtWidgets import QDialog

class StartSettingsWindow(QDialog):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.ui = Ui_StartSettingsWindow()
        self.ui.setupUi(self)

    def get_space_dimention(self):
        return self.ui.spaceDimentionSpinBox.value()