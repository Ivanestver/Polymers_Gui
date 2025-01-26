from uis.ui_stats import Ui_DlgStats
from PyQt6.QtWidgets import QDialog, QListWidgetItem, QMessageBox, QLabel
from PyQt6.QtCore import Qt
from polymer_view import GlobulaView, PolymerView
import math

class DlgStats(QDialog):
    def __init__(self, globula: GlobulaView, *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.ui = Ui_DlgStats()
        self.ui.setupUi(self)
        self.setWindowTitle(f'Статистика полимеров глобулы {globula.name}')
        self.ui.polymersListWidget.itemClicked.connect(self.__on_item_clicked_handler)
        self.POLYMER_ROLE = Qt.ItemDataRole.UserRole + 1

        self.__show_polymers_in_globula(globula)

    def __show_polymers_in_globula(self, globula: GlobulaView):
        for polymer in globula:
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
        start_monomer, end_monomer = polymer.get_start_end_monomers()
        length = math.sqrt(sum([(e - s) ** 2 for s, e in zip(start_monomer, end_monomer)]))
        self.ui.stats_layout.addRow("Расстояние между началом и концом", QLabel(str(length)))