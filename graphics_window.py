from uis.ui_graphics import Ui_DlgGraphics
from PyQt6.QtWidgets import QDialog
from polymer_view import GlobulaView
from math import ceil
from common_funcs import distance
from chart_window import LineChartWindow, BarChartWindow

class Settings:
    def __init__(self, globula: GlobulaView):
        self.globula = globula

class EndToEndDistanceSettings(Settings):
    def __init__(self, globula: GlobulaView, L: int):
        super().__init__(globula)
        self.L = L

class DlgGraphicsWindow(QDialog):
    def __init__(self, endToEndDistanceSettings: EndToEndDistanceSettings, *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.ui = Ui_DlgGraphics()
        self.ui.setupUi(self)
        self.endToEndDistanceSettings = endToEndDistanceSettings
        self.__add_end_to_end_distance_distribution()
        self.__add_polymer_length_distribution()

    def __add_end_to_end_distance_distribution(self):
        distribution = self.__get_distance_distribution()
        distribution = sorted(distribution.items(), key=lambda item: item[0])
        chartWindow = LineChartWindow(distribution, "Распределение межконцевых расстояний")
        self.ui.graphicsLayout.addWidget(chartWindow)

    def __add_polymer_length_distribution(self):
        distribution = {}
        for pol in self.endToEndDistanceSettings.globula:
            if pol.len() not in distribution.keys():
                distribution[pol.len()] = 1
            else:
                distribution[pol.len()] += 1

        chartWindow = BarChartWindow([(n, v) for n, v in distribution.items()], "Распределение длин полимеров")
        self.ui.graphicsLayout.addWidget(chartWindow)

    def __get_distance_distribution(self):
        distribution = {}
        mass_center = self.endToEndDistanceSettings.globula.get_mass_center()
        for polymer in self.endToEndDistanceSettings.globula:
            for monomer in polymer:
                dist = int(ceil(distance(monomer, mass_center)))
                if dist not in distribution.keys():
                    distribution[dist] = 1
                else:
                    distribution[dist] += 1
        
        return distribution
    