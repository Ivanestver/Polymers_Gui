import os
from PyQt6.QtWidgets import QMainWindow, QDialog, QWidget, QMessageBox, QListWidgetItem, QFileDialog
from polymer_view import GlobulaView
from uis.ui_main_window import Ui_MainWindow
import PyQt6.QtCore as core
from window_3d import Window3D
from alg.calc_alg import CalcAlg
from space import Space
from PyQt6.QtGui import QVector3D, QColor, QAction
from PyQt6.Qt3DRender import QPickEvent
from PyQt6.Qt3DExtras import QPhongMaterial
from stats_window import DlgStats, StatsInput
from save_to_formats import SaveToFile, SaveToMol, SaveToLammps
from alg.polymer_lib import MonomerType
from alg.config import Axis
from choose_axis_dlg import DlgChooseAxis
from cluster_analysis_ui import DlgClusterAnalysis
from radius_to_cut_window import DlgRadiusToCutWindow
from alg.common_funcs import distance

class MainWindow(QMainWindow):
    def __init__(self, space_dimention, *args, **kwargs):
        super().__init__(*args, **kwargs)

        Space.set_space_dimention(space_dimention)

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
        self.ui.loadFromDataFileBtn.triggered.connect(self.__load_from_data_file)

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
        add_action("Mark X clusters", self.__mark_x_axis)
        add_action("Mark Y clusters", self.__mark_y_axis)
        add_action("Mark Z clusters", self.__mark_z_axis)
        add_action("Mark common_clusters", self.__mark_common_clusters)
        add_action("Mark biggest clusters", self.__mark_biggest_clusters)
        add_action("Cluster analysis", self.__open_cluster_analysis_window)
        add_action("Mark as C", self.__mark_carbon)
        add_action("Cut by the desired raduis", self.__on_cut_by_radius)
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

        full_file_name = f'{path_to_save_dir}/Globula'
        real_full_file_name, _ = QFileDialog.getSaveFileName(self, "Сохранить", full_file_name, file_filters)
        if len(real_full_file_name) == 0:
            return

        file_name, file_contents = saveToFileObj.convert(globula)
        
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

    def __mark_axis(self, axis: Axis):
        current_globula = self.__get_current_globula()
        avg = self.ui.chainLengthSpinBox.value()
        if avg > 0:
            dlg = DlgChooseAxis(axis, self)
            if dlg.exec() == QDialog.DialogCode.Accepted:
                axis = dlg.get_chosen_axis()

        if axis == Axis.X_AXIS:
            clusters = current_globula.x_clusters(avg)
        elif axis == Axis.Y_AXIS:
            clusters = current_globula.y_clusters(avg)
        else:
            clusters = current_globula.z_clusters(avg)
        clusters.colorize()

    def __mark_x_axis(self):
        self.__mark_axis(Axis.X_AXIS)
    
    def __mark_y_axis(self):
        self.__mark_axis(Axis.Y_AXIS)

    def __mark_z_axis(self):
        self.__mark_axis(Axis.Z_AXIS)

    def __mark_common_clusters(self):
        current_globula = self.__get_current_globula()
        common_clusters = current_globula.common_clusters()
        for clusters in common_clusters:
            clusters.colorize()
        
    def __mark_biggest_clusters(self):
        current_globula = self.__get_current_globula()
        for percentage in range(20, 51, 10):
            clusters_percent = current_globula.biggest_common_clusters(percentage)
            for clusterView in clusters_percent:
                clusterView.colorize()
            QMessageBox.information(self, "Сохранение", f"Происходит сохранение со степенью кристалличности {percentage}")
            self.__on_save_to_lammps_btn_clicked()
            for clusterView in clusters_percent:
                clusterView.colorize(reset=True)

    def __open_cluster_analysis_window(self):
        current_globula = self.__get_current_globula()
        dlg = DlgClusterAnalysis(current_globula, self)
        dlg.exec()

    def __mark_carbon(self):
        current_globula = self.__get_current_globula()
        current_globula.reset()

    def __load_from_data_file(self):
        file_path = QFileDialog.getOpenFileName(self, "Выбрать частицу", "./slices", "*.json")
        if len(file_path[0]) == 0:
            return

        with open(file_path[0], 'r', encoding='utf-8') as file:
            contents = file.read()
            self.__read_from_json(contents)
            

    def __read_from_json(self, json_str: str):
        json = core.QByteArray(json_str.encode())
        doc = core.QJsonDocument.fromJson(json)
        Space.set_space_dimention(doc['space_dimention']['value'].toInt())
        self.ui.thresholdSpinBox.setValue(doc['threshold']['value'].toInt())
        self.ui.monomersCountSpinBox.setValue(doc['max_monomers_count']['value'].toInt())
        self.ui.polymersCountSpinBox.setValue(doc['polymers_count']['value'].toInt())
        self.ui.radiusSphereSpinBox.setMaximum(Space.space_dimention)
        radius_sphere = doc['radius_sphere']['value'].toInt()
        self.ui.radiusSphereSpinBox.setValue(radius_sphere)
        self.ui.sphereRadiusText.setText(f"Радиус сферы (не больше {radius_sphere})")
        new_globula: GlobulaView = GlobulaView.from_json(doc['globula'], self.ui.radiusSphereSpinBox.value())
        self.view.add_ready_made_globula(new_globula)
        globula_item = QListWidgetItem(new_globula.name)
        globula_item.setData(self.GLOBULA_ROLE, new_globula)
        self.ui.polymersListWidget.addItem(globula_item)

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

    def __on_cut_by_radius(self):
        dlg = DlgRadiusToCutWindow(self.ui.radiusSphereSpinBox.value(), self)
        if dlg.exec() == QDialog.DialogCode.Rejected:
            return

        desired_radius = dlg.get_desired_radius()
        current_globula = self.__get_current_globula()
        current_globula.reset()
        for pol in current_globula:
            for mon in pol:
                dist = distance(mon.coords, [Space.space_center.x(), Space.space_center.y(), Space.space_center.z()])
                if dist > desired_radius:
                    mon.type = MonomerType.Undefined

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
