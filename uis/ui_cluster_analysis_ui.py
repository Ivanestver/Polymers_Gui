# Form implementation generated from reading ui file '.\cluster_analysis_ui.ui'
#
# Created by: PyQt6 UI code generator 6.8.0
#
# WARNING: Any manual changes made to this file will be lost when pyuic6 is
# run again.  Do not edit this file unless you know what you are doing.


from PyQt6 import QtCore, QtGui, QtWidgets


class Ui_DlgClusterAnalysis(object):
    def setupUi(self, DlgClusterAnalysis):
        DlgClusterAnalysis.setObjectName("DlgClusterAnalysis")
        DlgClusterAnalysis.resize(682, 474)
        self.verticalLayout_3 = QtWidgets.QVBoxLayout(DlgClusterAnalysis)
        self.verticalLayout_3.setObjectName("verticalLayout_3")
        self.gridLayout_2 = QtWidgets.QGridLayout()
        self.gridLayout_2.setObjectName("gridLayout_2")
        self.verticalLayout_2 = QtWidgets.QVBoxLayout()
        self.verticalLayout_2.setObjectName("verticalLayout_2")
        self.label_2 = QtWidgets.QLabel(parent=DlgClusterAnalysis)
        self.label_2.setObjectName("label_2")
        self.verticalLayout_2.addWidget(self.label_2)
        self.trunkMolMassLayout = QtWidgets.QVBoxLayout()
        self.trunkMolMassLayout.setObjectName("trunkMolMassLayout")
        self.verticalLayout_2.addLayout(self.trunkMolMassLayout)
        self.verticalLayout_2.setStretch(0, 1)
        self.verticalLayout_2.setStretch(1, 3)
        self.gridLayout_2.addLayout(self.verticalLayout_2, 1, 0, 1, 1)
        self.verticalLayout = QtWidgets.QVBoxLayout()
        self.verticalLayout.setObjectName("verticalLayout")
        self.label = QtWidgets.QLabel(parent=DlgClusterAnalysis)
        self.label.setObjectName("label")
        self.verticalLayout.addWidget(self.label)
        self.molMassLayout = QtWidgets.QVBoxLayout()
        self.molMassLayout.setObjectName("molMassLayout")
        self.verticalLayout.addLayout(self.molMassLayout)
        self.verticalLayout.setStretch(0, 1)
        self.verticalLayout.setStretch(1, 3)
        self.gridLayout_2.addLayout(self.verticalLayout, 0, 0, 1, 1)
        self.gridLayout_3 = QtWidgets.QGridLayout()
        self.gridLayout_3.setObjectName("gridLayout_3")
        self.label_4 = QtWidgets.QLabel(parent=DlgClusterAnalysis)
        sizePolicy = QtWidgets.QSizePolicy(QtWidgets.QSizePolicy.Policy.Preferred, QtWidgets.QSizePolicy.Policy.Minimum)
        sizePolicy.setHorizontalStretch(0)
        sizePolicy.setVerticalStretch(0)
        sizePolicy.setHeightForWidth(self.label_4.sizePolicy().hasHeightForWidth())
        self.label_4.setSizePolicy(sizePolicy)
        self.label_4.setObjectName("label_4")
        self.gridLayout_3.addWidget(self.label_4, 0, 0, 1, 1)
        self.applyLambdaBtn = QtWidgets.QPushButton(parent=DlgClusterAnalysis)
        self.applyLambdaBtn.setObjectName("applyLambdaBtn")
        self.gridLayout_3.addWidget(self.applyLambdaBtn, 4, 1, 1, 1)
        self.label_3 = QtWidgets.QLabel(parent=DlgClusterAnalysis)
        self.label_3.setObjectName("label_3")
        self.gridLayout_3.addWidget(self.label_3, 2, 0, 1, 1)
        self.lambdaSpinBox = QtWidgets.QDoubleSpinBox(parent=DlgClusterAnalysis)
        self.lambdaSpinBox.setDecimals(6)
        self.lambdaSpinBox.setMaximum(9999.99)
        self.lambdaSpinBox.setObjectName("lambdaSpinBox")
        self.gridLayout_3.addWidget(self.lambdaSpinBox, 2, 1, 1, 1)
        self.tunedMolMassLayout = QtWidgets.QVBoxLayout()
        self.tunedMolMassLayout.setObjectName("tunedMolMassLayout")
        self.gridLayout_3.addLayout(self.tunedMolMassLayout, 3, 0, 1, 2)
        self.label_5 = QtWidgets.QLabel(parent=DlgClusterAnalysis)
        sizePolicy = QtWidgets.QSizePolicy(QtWidgets.QSizePolicy.Policy.Preferred, QtWidgets.QSizePolicy.Policy.Minimum)
        sizePolicy.setHorizontalStretch(0)
        sizePolicy.setVerticalStretch(0)
        sizePolicy.setHeightForWidth(self.label_5.sizePolicy().hasHeightForWidth())
        self.label_5.setSizePolicy(sizePolicy)
        self.label_5.setObjectName("label_5")
        self.gridLayout_3.addWidget(self.label_5, 0, 1, 1, 1)
        self.label_6 = QtWidgets.QLabel(parent=DlgClusterAnalysis)
        sizePolicy = QtWidgets.QSizePolicy(QtWidgets.QSizePolicy.Policy.Preferred, QtWidgets.QSizePolicy.Policy.Minimum)
        sizePolicy.setHorizontalStretch(0)
        sizePolicy.setVerticalStretch(0)
        sizePolicy.setHeightForWidth(self.label_6.sizePolicy().hasHeightForWidth())
        self.label_6.setSizePolicy(sizePolicy)
        self.label_6.setObjectName("label_6")
        self.gridLayout_3.addWidget(self.label_6, 1, 0, 1, 1)
        self.aSpinBox = QtWidgets.QDoubleSpinBox(parent=DlgClusterAnalysis)
        self.aSpinBox.setDecimals(6)
        self.aSpinBox.setObjectName("aSpinBox")
        self.gridLayout_3.addWidget(self.aSpinBox, 1, 1, 1, 1)
        self.gridLayout_3.setRowStretch(0, 1)
        self.gridLayout_3.setRowStretch(1, 1)
        self.gridLayout_3.setRowStretch(2, 1)
        self.gridLayout_3.setRowStretch(3, 3)
        self.gridLayout_3.setRowStretch(4, 1)
        self.gridLayout_2.addLayout(self.gridLayout_3, 0, 1, 2, 1)
        self.verticalLayout_3.addLayout(self.gridLayout_2)

        self.retranslateUi(DlgClusterAnalysis)
        QtCore.QMetaObject.connectSlotsByName(DlgClusterAnalysis)

    def retranslateUi(self, DlgClusterAnalysis):
        _translate = QtCore.QCoreApplication.translate
        DlgClusterAnalysis.setWindowTitle(_translate("DlgClusterAnalysis", "Dialog"))
        self.label_2.setText(_translate("DlgClusterAnalysis", "Молекулярно-массовое распределение (обрезанный)"))
        self.label.setText(_translate("DlgClusterAnalysis", "Молекулярно-массовое распределение"))
        self.label_4.setText(_translate("DlgClusterAnalysis", "Вид уравнения:"))
        self.applyLambdaBtn.setText(_translate("DlgClusterAnalysis", "Применить"))
        self.label_3.setText(_translate("DlgClusterAnalysis", "Значение λ"))
        self.label_5.setText(_translate("DlgClusterAnalysis", "y=a*exp(-λx)"))
        self.label_6.setText(_translate("DlgClusterAnalysis", "Значение a"))
