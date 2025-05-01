import sys
from PyQt6.QtWidgets import QApplication, QDialog
from main_window import MainWindow
from start_settings_window import StartSettingsWindow
from alg.common_funcs import get_side_by_monomers, Side

if __name__ == "__main__":
    app = QApplication(sys.argv)
    startSettingsWindow = StartSettingsWindow()
    responce = startSettingsWindow.exec()
    if responce == QDialog.DialogCode.Rejected:
        sys.exit()
    else:
        main_window = MainWindow(startSettingsWindow.get_space_dimention())
        main_window.show()
        app.exec()
        sys.exit()