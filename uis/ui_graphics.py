# Form implementation generated from reading ui file '.\graphics.ui'
#
# Created by: PyQt6 UI code generator 6.8.0
#
# WARNING: Any manual changes made to this file will be lost when pyuic6 is
# run again.  Do not edit this file unless you know what you are doing.


from PyQt6 import QtCore, QtGui, QtWidgets


class Ui_DlgGraphics(object):
    def setupUi(self, DlgGraphics):
        DlgGraphics.setObjectName("DlgGraphics")
        DlgGraphics.resize(400, 300)
        self.verticalLayout = QtWidgets.QVBoxLayout(DlgGraphics)
        self.verticalLayout.setObjectName("verticalLayout")
        self.graphicsLayout = QtWidgets.QGridLayout()
        self.graphicsLayout.setObjectName("graphicsLayout")
        self.verticalLayout.addLayout(self.graphicsLayout)
        self.buttonBox = QtWidgets.QDialogButtonBox(parent=DlgGraphics)
        self.buttonBox.setOrientation(QtCore.Qt.Orientation.Horizontal)
        self.buttonBox.setStandardButtons(QtWidgets.QDialogButtonBox.StandardButton.Cancel|QtWidgets.QDialogButtonBox.StandardButton.Ok)
        self.buttonBox.setObjectName("buttonBox")
        self.verticalLayout.addWidget(self.buttonBox)

        self.retranslateUi(DlgGraphics)
        self.buttonBox.accepted.connect(DlgGraphics.accept) # type: ignore
        self.buttonBox.rejected.connect(DlgGraphics.reject) # type: ignore
        QtCore.QMetaObject.connectSlotsByName(DlgGraphics)

    def retranslateUi(self, DlgGraphics):
        _translate = QtCore.QCoreApplication.translate
        DlgGraphics.setWindowTitle(_translate("DlgGraphics", "Dialog"))
