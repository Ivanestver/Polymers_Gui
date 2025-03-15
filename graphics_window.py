from uis.ui_graphics import Ui_DlgGraphics
from PyQt6.QtWidgets import QDialog
from polymer_view import GlobulaView
from math import ceil, pi
from alg.common_funcs import distance
from chart_window import BarChartWindow

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
        self.__add_dist_to_mass_center_distribution()
        self.__add_dist_g_r()

    def __add_end_to_end_distance_distribution(self):
        distribution = {}
        for pol in self.endToEndDistanceSettings.globula:
            start_monomer, end_monomer = pol.get_start_end_monomers()
            dist = distance(start_monomer, end_monomer)
            def add_to_dist(key):
                if key in distribution.keys():
                    distribution[key] += 1
                else:
                    distribution[key] = 1

            part = dist - int(dist)
            if part < 0.25:
                add_to_dist(int(dist))
            elif 0.25 <= part < 0.5:
                add_to_dist(int(dist) + 0.5)
            elif 0.5 <= part < 0.75:
                add_to_dist(int(dist) + 0.75)
            else:
                add_to_dist(ceil(dist))
            
        distribution = sorted(distribution.items(), key=lambda item: item[0])
        chartWindow = BarChartWindow(distribution, "Распределение межконцевых расстояний")
        self.ui.graphicsLayout.addWidget(chartWindow, 0, 0)

    def __add_polymer_length_distribution(self):
        distribution = {}
        for pol in self.endToEndDistanceSettings.globula:
            if pol.len() not in distribution.keys():
                distribution[pol.len()] = 1
            else:
                distribution[pol.len()] += 1

        points = [(n, v) for n, v in distribution.items()]
        points = sorted(points, key=lambda item: item[0])
        chartWindow = BarChartWindow(points, "Распределение длин полимеров")
        self.ui.graphicsLayout.addWidget(chartWindow, 0, 1)

    def __add_dist_to_mass_center_distribution(self):
        distribution = {i: 0 for i in range(1, self.endToEndDistanceSettings.L)}
        mass_center = self.endToEndDistanceSettings.globula.get_mass_center()
        for polymer in self.endToEndDistanceSettings.globula:
            for monomer in polymer:
                rj = int(ceil(distance(monomer, mass_center)))
                distribution[rj] += 1

        chartWindow = BarChartWindow(sorted(distribution.items(), key=lambda item: item[0]), "Распределение расстояний от точки до центра масс")
        self.ui.graphicsLayout.addWidget(chartWindow, 1, 0)
    
    def __add_dist_g_r(self):
        distribution = {i: 0 for i in range(1, self.endToEndDistanceSettings.L)}
        mass_center = self.endToEndDistanceSettings.globula.get_mass_center()
        for polymer in self.endToEndDistanceSettings.globula:
            for monomer in polymer:
                rj = int(ceil(distance(monomer, mass_center)))
                distribution[rj] += 1

        for k in distribution.keys():
            distribution[k] /= 2 * pi * k * k

        chartWindow = BarChartWindow(sorted(distribution.items(), key=lambda item: item[0]), "g(r)")
        self.ui.graphicsLayout.addWidget(chartWindow, 1, 1)

    def __create_distance_distribution_function(self):
        # Считаем плотность
        N = sum([pol.len() for pol in self.endToEndDistanceSettings.globula])
        ro = N / ((4.0 / 3) * pi * (self.endToEndDistanceSettings.L ** 3))

        # Строим распределение расстояний от центра масс к точке
        distribution = {}
        mass_center = self.endToEndDistanceSettings.globula.get_mass_center()
        for polymer in self.endToEndDistanceSettings.globula:
            for monomer in polymer:
                dist = int(ceil(distance(monomer, mass_center)))
                if dist not in distribution.keys():
                    distribution[dist] = 1
                else:
                    distribution[dist] += 1

                    
        def distance_distribution_function(rj):
            dist = int(ceil(rj))
            # return (1 / ro) * (distribution[dist] / (2 * pi * (rj ** 2)))
            return distribution[dist]
        
        return distance_distribution_function
    