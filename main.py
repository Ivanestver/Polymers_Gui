import sys
from PyQt6.QtWidgets import QApplication, QDialog
from main_window import MainWindow
from start_settings_window import StartSettingsWindow
from alg.common_funcs import get_side_by_monomers, Side

if __name__ == "__main__":

    assert get_side_by_monomers([0, 0, 0], [1, 0, 0]) == Side.Forward, "You fucked up"
    assert get_side_by_monomers([1, 0, 0], [0, 0, 0]) == Side.Backward, "You fucked up"
    assert get_side_by_monomers([0, 0, 0], [0, 1, 0]) == Side.Left, "You fucked up"
    assert get_side_by_monomers([0, 1, 0], [0, 0, 0]) == Side.Right, "You fucked up"
    assert get_side_by_monomers([0, 0, 0], [0, 0, 1]) == Side.Up, "You fucked up"
    assert get_side_by_monomers([0, 0, 1], [0, 0, 0]) == Side.Down, "You fucked up"
    assert get_side_by_monomers([0, 0, 0], [1, 1, 1]) == Side.UpLeftForward, "You fucked up"
    assert get_side_by_monomers([0, 0, 0], [1, 0, 1]) == Side.UpForward, "You fucked up"
    assert get_side_by_monomers([0, 1, 0], [1, 0, 1]) == Side.UpRightForward, "You fucked up"
    assert get_side_by_monomers([0, 0, 0], [1, 1, 0]) == Side.LeftForward, "You fucked up"
    assert get_side_by_monomers([0, 1, 0], [1, 0, 0]) == Side.RightForward, "You fucked up"
    assert get_side_by_monomers([0, 0, 1], [1, 1, 0]) == Side.DownLeftForward, "You fucked up"
    assert get_side_by_monomers([0, 0, 1], [1, 0, 0]) == Side.DownForward, "You fucked up"
    assert get_side_by_monomers([0, 1, 1], [1, 0, 0]) == Side.DownRightForward, "You fucked up"
    assert get_side_by_monomers([0, 0, 0], [0, 1, 1]) == Side.LeftUp, "You fucked up"
    assert get_side_by_monomers([0, 0, 1], [0, 1, 0]) == Side.LeftDown, "You fucked up"
    assert get_side_by_monomers([0, 1, 0], [0, 0, 1]) == Side.RightUp, "You fucked up"
    assert get_side_by_monomers([0, 1, 1], [0, 0, 0]) == Side.RightDown, "You fucked up"
    assert get_side_by_monomers([1, 0, 0], [0, 1, 1]) == Side.UpLeftBackward, "You fucked up"
    assert get_side_by_monomers([1, 0, 0], [0, 0, 1]) == Side.UpBackward, "You fucked up"
    assert get_side_by_monomers([1, 1, 0], [0, 0, 1]) == Side.UpRightBackward, "You fucked up"
    assert get_side_by_monomers([1, 0, 0], [0, 1, 0]) == Side.LeftBackward, "You fucked up"
    assert get_side_by_monomers([1, 1, 0], [0, 0, 0]) == Side.RightBackward, "You fucked up"
    assert get_side_by_monomers([1, 0, 1], [0, 1, 0]) == Side.DownLeftBackward, "You fucked up"
    assert get_side_by_monomers([1, 0, 1], [0, 0, 0]) == Side.DownBackward, "You fucked up"
    assert get_side_by_monomers([1, 1, 1], [0, 0, 0]) == Side.DownRightBackward, "You fucked up"
    
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