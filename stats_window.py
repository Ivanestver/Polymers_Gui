from uis.ui_stats import Ui_DlgStats
from PyQt6.QtWidgets import QDialog, QListWidgetItem, QMessageBox, QLabel
from PyQt6.QtCore import Qt
from polymer_view import GlobulaView, PolymerView
from math import sqrt, ceil
from common_funcs import distance
from graphics_window import DlgGraphicsWindow, EndToEndDistanceSettings

class StatsInput:
    def __init__(self):
        self.L = 0
        self.globula: GlobulaView = None

class DlgStats(QDialog):
    def __init__(self, inputData: StatsInput, *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.ui = Ui_DlgStats()
        self.ui.setupUi(self)
        self.setWindowTitle(f'Статистика полимеров глобулы {inputData.globula.name}')
        self.ui.polymersListWidget.itemClicked.connect(self.__on_item_clicked_handler)
        self.POLYMER_ROLE = Qt.ItemDataRole.UserRole + 1
        self.L = inputData.L
        self.globula = inputData.globula

        self.__show_polymers_in_globula()
        self.__show_globula_stats()
        self.ui.graphicsBtn.clicked.connect(self.__show_chart)

    def __show_polymers_in_globula(self):
        for polymer in self.globula:
            polymer_item = QListWidgetItem(polymer.name)
            polymer_item.setData(self.POLYMER_ROLE, polymer)
            self.ui.polymersListWidget.addItem(polymer_item)

    def __on_item_clicked_handler(self, item: QListWidgetItem):
        polymer = item.data(self.POLYMER_ROLE)
        if polymer is None:
            QMessageBox.critical(self, "Ошибка", "Выбранный элемент не содержит информацию о полимере")
            return
            
        self.__show_statistics(polymer)

    def __show_statistics(self, polymer: PolymerView):
        # Очищаем окно статистики
        while (self.ui.stats_layout.rowCount() > 0):
            self.ui.stats_layout.removeRow(self.ui.stats_layout.rowCount() - 1)

        # Показываем длину цепи
        self.ui.stats_layout.addRow("Длина цепи", QLabel(str(polymer.len())))

        # Показываем расстояние между началом и концом
        self.ui.stats_layout.addRow("Расстояние между началом и концом", QLabel(str(round(self.__get_start_end_monomer_distance(polymer), 2))))

    def __show_chart(self):
        graphics_ui = DlgGraphicsWindow(EndToEndDistanceSettings(self.globula, self.L), self)
        graphics_ui.exec()

    def __get_distance_distribution(self):
        distribution = {i: 0 for i in range(self.L)}
        mass_center = self.globula.get_mass_center()
        for polymer in self.globula:
            for monomer in polymer:
                distribution[int(ceil(distance(monomer, mass_center)))] += 1
        
        return distribution
    
    def __get_start_end_monomer_distance(self, polymer: PolymerView):
        start_monomer, end_monomer = polymer.get_start_end_monomers()
        return distance(start_monomer, end_monomer)

    def __show_globula_stats(self):
        # Считаем центр масс
        mass_center = self.globula.get_mass_center()
        self.ui.massCenterLabel.setText(f'({round(mass_center[0], 2)}, {round(mass_center[1], 2)}, {round(mass_center[2], 2)})')

        distribution = self.__get_distance_distribution()
        N = sum([v for i, v in distribution.items()])

        # Рассчитываем матожидание
        mean_value = sum([i * (v / N) for i, v in distribution.items()])
        self.ui.meanValueLabel.setText(str(round(mean_value, 2)))

        # Рассчитываем дисперсию
        dispersity = sum([i * i * (v / N) for i, v in distribution.items()]) - mean_value ** 2
        self.ui.dispersityLabel.setText(str(round(dispersity, 2)))

        # Рассчитываем стандартное отклонение
        stdDeviation = sqrt(dispersity)
        self.ui.stdDeviationLabel.setText(str(round(stdDeviation, 2)))