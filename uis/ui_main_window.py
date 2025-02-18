# Form implementation generated from reading ui file '.\main_window.ui'
#
# Created by: PyQt6 UI code generator 6.8.0
#
# WARNING: Any manual changes made to this file will be lost when pyuic6 is
# run again.  Do not edit this file unless you know what you are doing.


from PyQt6 import QtCore, QtGui, QtWidgets


class Ui_MainWindow(object):
    def setupUi(self, MainWindow):
        MainWindow.setObjectName("MainWindow")
        MainWindow.resize(926, 628)
        MainWindow.setContextMenuPolicy(QtCore.Qt.ContextMenuPolicy.NoContextMenu)
        self.horizontalLayout = QtWidgets.QHBoxLayout(MainWindow)
        self.horizontalLayout.setObjectName("horizontalLayout")
        self.mainLayout = QtWidgets.QVBoxLayout()
        self.mainLayout.setObjectName("mainLayout")
        self.horizontalLayout.addLayout(self.mainLayout)
        self.formLayout = QtWidgets.QFormLayout()
        self.formLayout.setObjectName("formLayout")
        self.label = QtWidgets.QLabel(parent=MainWindow)
        self.label.setObjectName("label")
        self.formLayout.setWidget(0, QtWidgets.QFormLayout.ItemRole.LabelRole, self.label)
        self.filesCountSpinBox = QtWidgets.QSpinBox(parent=MainWindow)
        self.filesCountSpinBox.setMinimum(1)
        self.filesCountSpinBox.setObjectName("filesCountSpinBox")
        self.formLayout.setWidget(0, QtWidgets.QFormLayout.ItemRole.FieldRole, self.filesCountSpinBox)
        self.label_2 = QtWidgets.QLabel(parent=MainWindow)
        self.label_2.setObjectName("label_2")
        self.formLayout.setWidget(1, QtWidgets.QFormLayout.ItemRole.LabelRole, self.label_2)
        self.polymersCountSpinBox = QtWidgets.QSpinBox(parent=MainWindow)
        self.polymersCountSpinBox.setMinimum(1)
        self.polymersCountSpinBox.setObjectName("polymersCountSpinBox")
        self.formLayout.setWidget(1, QtWidgets.QFormLayout.ItemRole.FieldRole, self.polymersCountSpinBox)
        self.label_3 = QtWidgets.QLabel(parent=MainWindow)
        self.label_3.setObjectName("label_3")
        self.formLayout.setWidget(2, QtWidgets.QFormLayout.ItemRole.LabelRole, self.label_3)
        self.thresholdSpinBox = QtWidgets.QDoubleSpinBox(parent=MainWindow)
        self.thresholdSpinBox.setDecimals(4)
        self.thresholdSpinBox.setMinimum(0.001)
        self.thresholdSpinBox.setObjectName("thresholdSpinBox")
        self.formLayout.setWidget(2, QtWidgets.QFormLayout.ItemRole.FieldRole, self.thresholdSpinBox)
        self.label_4 = QtWidgets.QLabel(parent=MainWindow)
        self.label_4.setObjectName("label_4")
        self.formLayout.setWidget(3, QtWidgets.QFormLayout.ItemRole.LabelRole, self.label_4)
        self.monomersCountSpinBox = QtWidgets.QSpinBox(parent=MainWindow)
        self.monomersCountSpinBox.setMinimum(10)
        self.monomersCountSpinBox.setMaximum(2048)
        self.monomersCountSpinBox.setObjectName("monomersCountSpinBox")
        self.formLayout.setWidget(3, QtWidgets.QFormLayout.ItemRole.FieldRole, self.monomersCountSpinBox)
        self.sphereRadiusText = QtWidgets.QLabel(parent=MainWindow)
        self.sphereRadiusText.setObjectName("sphereRadiusText")
        self.formLayout.setWidget(4, QtWidgets.QFormLayout.ItemRole.LabelRole, self.sphereRadiusText)
        self.radiusSphereSpinBox = QtWidgets.QSpinBox(parent=MainWindow)
        self.radiusSphereSpinBox.setObjectName("radiusSphereSpinBox")
        self.formLayout.setWidget(4, QtWidgets.QFormLayout.ItemRole.FieldRole, self.radiusSphereSpinBox)
        self.calculateBtn = QtWidgets.QPushButton(parent=MainWindow)
        self.calculateBtn.setObjectName("calculateBtn")
        self.formLayout.setWidget(6, QtWidgets.QFormLayout.ItemRole.FieldRole, self.calculateBtn)
        self.gridLayout = QtWidgets.QGridLayout()
        self.gridLayout.setObjectName("gridLayout")
        self.btnUp = QtWidgets.QPushButton(parent=MainWindow)
        self.btnUp.setObjectName("btnUp")
        self.gridLayout.addWidget(self.btnUp, 0, 1, 1, 1)
        self.btnLeft = QtWidgets.QPushButton(parent=MainWindow)
        self.btnLeft.setObjectName("btnLeft")
        self.gridLayout.addWidget(self.btnLeft, 1, 0, 1, 1)
        self.btnDown = QtWidgets.QPushButton(parent=MainWindow)
        self.btnDown.setObjectName("btnDown")
        self.gridLayout.addWidget(self.btnDown, 1, 1, 1, 1)
        self.btnRight = QtWidgets.QPushButton(parent=MainWindow)
        self.btnRight.setObjectName("btnRight")
        self.gridLayout.addWidget(self.btnRight, 1, 2, 1, 1)
        self.formLayout.setLayout(7, QtWidgets.QFormLayout.ItemRole.SpanningRole, self.gridLayout)
        self.polymersListWidget = QtWidgets.QListWidget(parent=MainWindow)
        self.polymersListWidget.setContextMenuPolicy(QtCore.Qt.ContextMenuPolicy.ActionsContextMenu)
        self.polymersListWidget.setObjectName("polymersListWidget")
        self.formLayout.setWidget(8, QtWidgets.QFormLayout.ItemRole.SpanningRole, self.polymersListWidget)
        self.buildMoreBtn = QtWidgets.QPushButton(parent=MainWindow)
        self.buildMoreBtn.setEnabled(False)
        self.buildMoreBtn.setObjectName("buildMoreBtn")
        self.formLayout.setWidget(6, QtWidgets.QFormLayout.ItemRole.LabelRole, self.buildMoreBtn)
        self.toShowResultCheckBox = QtWidgets.QCheckBox(parent=MainWindow)
        self.toShowResultCheckBox.setChecked(True)
        self.toShowResultCheckBox.setTristate(False)
        self.toShowResultCheckBox.setObjectName("toShowResultCheckBox")
        self.formLayout.setWidget(5, QtWidgets.QFormLayout.ItemRole.LabelRole, self.toShowResultCheckBox)
        self.buildCristallCheckBox = QtWidgets.QCheckBox(parent=MainWindow)
        self.buildCristallCheckBox.setObjectName("buildCristallCheckBox")
        self.formLayout.setWidget(5, QtWidgets.QFormLayout.ItemRole.FieldRole, self.buildCristallCheckBox)
        self.horizontalLayout.addLayout(self.formLayout)
        self.horizontalLayout.setStretch(0, 5)

        self.retranslateUi(MainWindow)
        QtCore.QMetaObject.connectSlotsByName(MainWindow)

    def retranslateUi(self, MainWindow):
        _translate = QtCore.QCoreApplication.translate
        MainWindow.setWindowTitle(_translate("MainWindow", "Dialog"))
        self.label.setText(_translate("MainWindow", "Количество файлов"))
        self.label_2.setText(_translate("MainWindow", "Количество полимеров"))
        self.label_3.setText(_translate("MainWindow", "Значение порога"))
        self.label_4.setText(_translate("MainWindow", "Количество мономеров"))
        self.sphereRadiusText.setText(_translate("MainWindow", "TextLabel"))
        self.calculateBtn.setText(_translate("MainWindow", "Рассчитать"))
        self.btnUp.setText(_translate("MainWindow", "^"))
        self.btnLeft.setText(_translate("MainWindow", "<"))
        self.btnDown.setText(_translate("MainWindow", "v"))
        self.btnRight.setText(_translate("MainWindow", ">"))
        self.buildMoreBtn.setText(_translate("MainWindow", "Достроить"))
        self.toShowResultCheckBox.setText(_translate("MainWindow", "Отобразить результат"))
        self.buildCristallCheckBox.setText(_translate("MainWindow", "Строить кристалл"))
