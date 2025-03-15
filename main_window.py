import os
from PyQt6.QtCore import Qt
from PyQt6.QtWidgets import QDialog, QWidget, QMessageBox, QListWidgetItem, QFileDialog
from polymer_view import GlobulaView
from uis.ui_main_window import Ui_MainWindow
import PyQt6.QtCore as core
from window_3d import Window3D
from alg.calc_alg import CalcAlg
from alg.polymer_lib import Monomer, MonomerType
from space import Space
from PyQt6.QtGui import QVector3D, QColor, QAction
from PyQt6.Qt3DRender import QPickEvent
from PyQt6.Qt3DExtras import QPhongMaterial
from stats_window import DlgStats, StatsInput
from save_to_formats import SaveToFile, SaveToMol, SaveToLammps
from operator import sub
from enum import IntEnum

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
        self.ui.buildMoreBtn.clicked.connect(self.on_build_more_btn_clicked)

        self.ui.btnUp.clicked.connect(self.on_btn_up_clicked)
        self.ui.btnLeft.clicked.connect(self.on_btn_left_clicked)
        self.ui.btnDown.clicked.connect(self.on_btn_down_clicked)
        self.ui.btnRight.clicked.connect(self.on_btn_right_clicked)
        self.setup_context_menu()

        self.ui.sphereRadiusText.setText(f"Радиус сферы (не больше {Space.space_dimention})")
        self.ui.radiusSphereSpinBox.setMaximum(Space.space_dimention)

        self.ui.polymersListWidget.itemDoubleClicked.connect(self.__on_globula_list_item_double_clicked)
        self.ui.polymersListWidget.itemClicked.connect(self.__on_polymer_item_list_clicked)
        self.GLOBULA_ROLE = core.Qt.ItemDataRole.UserRole + 1

    def setup_context_menu(self):
        def add_action(name, handler):
            action = QAction(name, self)
            action.triggered.connect(handler)
            self.ui.polymersListWidget.addAction(action)
        
        add_action("Save to mol", self.__on_save_to_file_action_triggered)
        add_action("Save to lammps", self.__on_save_to_lammps_btn_clicked)
        add_action("Mark O", self.__mark_O_part)
        add_action("Remove", self.__on_remove_globula_triggered)

    def __on_save_to_file_action_triggered(self):
        self.__save_to_file("*.mol2", SaveToMol())
    
    def __on_save_to_lammps_btn_clicked(self):
        self.__save_to_file("*.data", SaveToLammps())

    def __save_to_file(self, file_filters: str, saveToFileObj: SaveToFile):
        curr_globula_number = self.ui.polymersListWidget.currentRow()
        if curr_globula_number < 0 or curr_globula_number >= self.view.get_globulas_count():
            return

        globula = self.view.get_globula(curr_globula_number)
        path_to_save_dir = "./slices"

        if not os.path.isdir(path_to_save_dir):
            os.mkdir(path_to_save_dir)

        file_name, file_contents = saveToFileObj.convert(globula)
        full_file_name = f'{path_to_save_dir}/{file_name}'
        real_full_file_name, _ = QFileDialog.getSaveFileName(self, "Сохранить", full_file_name, file_filters)
        if len(real_full_file_name) == 0:
            return
        
        if os.path.exists(real_full_file_name):
            os.remove(real_full_file_name)

        with open(real_full_file_name, 'w', encoding='utf-8') as file:
            file.write(file_contents)

        real_full_file_name_split = real_full_file_name.split('/')
        real_file_name = real_full_file_name_split[-1]
        real_full_file_name_split.remove(real_file_name)
        file_name = real_file_name.split('.')[0]
        json_file_name = '/'.join(real_full_file_name_split) + f'/{file_name}.json'
        if os.path.exists(json_file_name):
            os.remove(json_file_name)

        with open(json_file_name, 'w', encoding='utf-8') as file:
            file.write(self.__get_modelling_info_as_json(globula))

        QMessageBox.information(self, "Information", f"The selected globula was saved to {real_full_file_name}")

    def __on_remove_globula_triggered(self):
        current_globula = self.__get_current_globula()
        self.view.remove_globula(current_globula)
        self.ui.polymersListWidget.takeItem(self.ui.polymersListWidget.currentRow())

    def __mark_O_part(self):
        class Direction(IntEnum):
            Up = 0 * 100 + 0 * 10 + 1,
            Down = 0 * 100 + 0 * 10 + -1,
            Left = 0 * 100 + 1 * 10 + 0,
            Right = 0 * 100 + -1 * 10 + 0,
            Forward = 1 * 100 + 0 * 10 + 0,
            Backward = -1 * 100 + 0 * 10 + 0

        def get_direction(prev: tuple, curr: tuple) -> Direction:
            tuples_sub = tuple(map(sub, curr, prev))
            return Direction(sum(tuples_sub[i] * (10 ** i) for i in range(len(tuples_sub))))
        
        current_globula = self.__get_current_globula()
        for polymer in current_globula:
            previous_monomer = polymer[0]
            first_monomer_in_row = polymer[0]
            current_direction = Direction.Up
            monomers_in_row = 1
            for i in range(1, polymer.len()):
                current_monomer = polymer[i]
                direction = get_direction(previous_monomer, current_monomer)
                if direction == current_direction:
                    monomers_in_row += 1
                else:
                    if monomers_in_row >= self.ui.chainLengthSpinBox.value():
                        # Помечаем полученную цепочку кислородом
                        m: Monomer = first_monomer_in_row
                        while m != previous_monomer:
                            m.type = MonomerType.Owise
                            m = m.next_monomer
                        m.type = MonomerType.Owise

                    monomers_in_row = 2
                    current_direction = direction
                    first_monomer_in_row = previous_monomer
                previous_monomer = current_monomer

        self.repaint()

    def __get_modelling_info_as_json(self, globula: GlobulaView):
        space_dimention, polymers_count, threshold, monomers_count, radius_sphere = self.__get_modelling_parameters()
        json_dict = {
            'space_dimention': {
                'value': space_dimention,
                'name': "Название"
            },
            'polymers_count': {
                'value': polymers_count,
                'name': "Количество полимеров"
            },
            'threshold': {
                'value': threshold,
                'name': "Порог принятия"
            },
            'max_monomers_count': {
                'value': monomers_count,
                'name': "Максимальное количество мономеров"
            },
            'radius_sphere': {
                'value': radius_sphere,
                'name': "Радиус сферы"
            },
            'globula': {
                'value': globula.to_json(),
                'name': "Глобула"
            }
        }
        json = core.QJsonDocument(json_dict)
        return json.toJson().data().decode("utf-8")
        
    def __get_modelling_parameters(self):
        return (Space.space_dimention,
                self.ui.polymersCountSpinBox.value(),
                self.ui.thresholdSpinBox.value(),
                self.ui.monomersCountSpinBox.value(),
                self.ui.radiusSphereSpinBox.value())


    def on_calc_btn_clicked(self):
        globuls_count = self.ui.filesCountSpinBox.value()
        polymers_count = self.ui.polymersCountSpinBox.value()
        accept_threshold = self.ui.thresholdSpinBox.value()
        monomers_count = self.ui.monomersCountSpinBox.value()
        sphere_radius = self.ui.radiusSphereSpinBox.value()
        calcAlg = CalcAlg(globuls_count, polymers_count, accept_threshold, monomers_count, sphere_radius)
        if self.ui.buildCristallCheckBox.isChecked():
            finished_polymers = calcAlg.calc_as_cristall()
        else:
            finished_polymers = calcAlg.calc_simultaneously()
        if len(finished_polymers) == 0:
            QMessageBox.warning(self, "Внимание", "Не удалось построить многочлены. Проверьте настройки и повторите попытку")
            return

        new_globula = self.view.add_globula(finished_polymers, self.ui.toShowResultCheckBox.isChecked())
        globula_item = QListWidgetItem(new_globula.name)
        globula_item.setData(self.GLOBULA_ROLE, new_globula)
        self.ui.polymersListWidget.addItem(globula_item)
    
    def on_build_more_btn_clicked(self):
        globula = self.__get_current_globula()
        if globula == None:
            QMessageBox.warning(self, "Внимание", "Выберите глобулу")
            return

        polymers = list()
        for pol in globula:
            polymers.append(pol.polymer)
        globuls_count = self.ui.filesCountSpinBox.value()
        polymers_count = self.ui.polymersCountSpinBox.value()
        accept_threshold = self.ui.thresholdSpinBox.value()
        monomers_count = self.ui.monomersCountSpinBox.value()
        sphere_radius = self.ui.radiusSphereSpinBox.value()
        calcAlg = CalcAlg(globuls_count, polymers_count, accept_threshold, monomers_count, sphere_radius)
        if self.ui.buildCristallCheckBox.isChecked():
            finished_polymers = calcAlg.build_more_as_cristall(polymers)
        else:
            finished_polymers = calcAlg.build_more_simultaneously(polymers)
        if finished_polymers == None or len(finished_polymers) == 0:
            QMessageBox.warning(self, "Внимание", "Не удалось построить глобулу. Проверьте настройки и повторите попытку")
            return

        self.view.remove_globula(globula)
        new_globula = self.view.add_globula(finished_polymers, self.ui.toShowResultCheckBox.isChecked())
        globula_item = self.ui.polymersListWidget.currentItem()
        globula_item.setText(new_globula.name)
        globula_item.setData(self.GLOBULA_ROLE, new_globula)

    def on_picker_clicked(self, e: QPickEvent):
        material = QPhongMaterial()
        material.setAmbient(QColor.fromRgb(0, 0, 0))
        e.entity().addComponent(material)
        pass

    def __get_current_globula(self):
        curr_polymer_number = self.ui.polymersListWidget.currentRow()
        if curr_polymer_number < 0 or curr_polymer_number >= self.view.get_globulas_count():
            return None

        return self.view.get_globula(curr_polymer_number)

    def on_btn_up_clicked(self):
        globula = self.__get_current_globula() 
        if globula is not None:
            globula.setRotationX(globula.rotationX() + 10)

    def on_btn_down_clicked(self):
        globula = self.__get_current_globula() 
        if globula is not None:
            globula.setRotationX(globula.rotationX() - 10)

    def on_btn_left_clicked(self):
        globula = self.__get_current_globula() 
        if globula is not None:
            globula.setRotationY(globula.rotationY() + 10)

    def on_btn_right_clicked(self):
        globula = self.__get_current_globula() 
        if globula is not None:
            globula.setRotationY(globula.rotationY() - 10)

    def __on_globula_list_item_double_clicked(self, item: QListWidgetItem):
        globula = item.data(self.GLOBULA_ROLE)
        if globula is None:
            QMessageBox.critical(self, "Ошибка", "Отсутствует информация по указанной глобуле")
            return

        inputData = StatsInput()
        inputData.L = self.ui.radiusSphereSpinBox.value()
        inputData.globula = globula

        dlg = DlgStats(inputData)
        dlg.exec()

    def __on_polymer_item_list_clicked(self, listItem: QListWidgetItem):
        self.ui.buildMoreBtn.setEnabled(True)
