# Form implementation generated from reading ui file '.\choose_axis.ui'
#
# Created by: PyQt6 UI code generator 6.8.0
#
# WARNING: Any manual changes made to this file will be lost when pyuic6 is
# run again.  Do not edit this file unless you know what you are doing.


from PyQt6 import QtCore, QtGui, QtWidgets


class Ui_DlgChooseAxis(object):
    def setupUi(self, DlgChooseAxis):
        DlgChooseAxis.setObjectName("DlgChooseAxis")
        DlgChooseAxis.resize(94, 67)
        self.verticalLayout = QtWidgets.QVBoxLayout(DlgChooseAxis)
        self.verticalLayout.setObjectName("verticalLayout")
        self.cmbAxis = QtWidgets.QComboBox(parent=DlgChooseAxis)
        self.cmbAxis.setObjectName("cmbAxis")
        self.verticalLayout.addWidget(self.cmbAxis)
        self.btnAccept = QtWidgets.QPushButton(parent=DlgChooseAxis)
        self.btnAccept.setObjectName("btnAccept")
        self.verticalLayout.addWidget(self.btnAccept)

        self.retranslateUi(DlgChooseAxis)
        self.btnAccept.clicked.connect(DlgChooseAxis.accept) # type: ignore
        QtCore.QMetaObject.connectSlotsByName(DlgChooseAxis)

    def retranslateUi(self, DlgChooseAxis):
        _translate = QtCore.QCoreApplication.translate
        DlgChooseAxis.setWindowTitle(_translate("DlgChooseAxis", "Dialog"))
        self.btnAccept.setText(_translate("DlgChooseAxis", "OK"))
