from uis.ui_cluster_analysis_ui import Ui_DlgClusterAnalysis
from PyQt6.QtWidgets import QDialog
from polymer_view import GlobulaView
from chart_window import BarChartWindow
import math

class DlgClusterAnalysis(QDialog):
    def __init__(self, globula: GlobulaView, *args):
        super().__init__(*args)

        self.ui = Ui_DlgClusterAnalysis()
        self.ui.setupUi(self)
        self.globula = globula
        self.division = 20

        self.distro = self.__get_distro()
        self.__show_mol_mass_distribution(self.distro)
        self.__trunk_distro(self.distro)
        self.__show_trunk_mol_mass_distribution(self.distro)
        self.__define_lambda()
        self.__show_tuned_mol_mass_distribution()
        self.ui.aSpinBox.valueChanged.connect(self.__show_tuned_mol_mass_distribution)
        self.ui.lambdaSpinBox.valueChanged.connect(self.__show_tuned_mol_mass_distribution)

    def __get_distro(self):
        clusters_x, clusters_y, clusters_z = self.globula.common_clusters()
        max_mass = max([
            max([cluster.size for cluster in clusters_x]) if clusters_x.len() != 0 else 0,
            max([cluster.size for cluster in clusters_y]) if clusters_y.len() != 0 else 0,
            max([cluster.size for cluster in clusters_z]) if clusters_z.len() != 0 else 0,
        ])

        self.H = int(max_mass / self.division)
        distro = { i * self.H: 0 for i in range(1, self.division if self.H * self.division == max_mass else self.division + 1) }
        for cluster_axis in [clusters_x, clusters_y, clusters_z]:
            for cluster in cluster_axis:
                for key in distro.keys():
                    if key - self.H < cluster.size <= key:
                        distro[key] += 1

        keys_to_remove = sorted([key for key in distro.keys() if distro[key] == 0], reverse=True)
        for key in keys_to_remove:
            del distro[key]
        
        return distro

    def __show_mol_mass_distribution(self, distro: dict[int, int]):
        chart = BarChartWindow([(key, value) for key, value in distro.items()], "Молекулярно-массовое распределение")
        self.ui.molMassLayout.addWidget(chart)

    def __show_trunk_mol_mass_distribution(self, distro: dict[int, int]):
        chart = BarChartWindow([(key, value) for key, value in distro.items()], "Молекулярно-массовое распределение (усечённое)")
        self.ui.trunkMolMassLayout.addWidget(chart)

    def __trunk_distro(self, distro: dict[int, int]):
        keys = [key for key in distro.keys()]
        for i in range(1, len(keys)):
            if distro[keys[i - 1]] < distro[keys[i]]:
                distro[keys[i]] = distro[keys[i - 1]]

    def __define_lambda(self):
        xs = list([key for key in self.distro.keys()])
        xs_square = list((x * x for x in xs))
        zs = list([math.log(y) for _, y in self.distro.items()])
        xzs = list([x * z for x, z in zip(xs, zs)])
        N = len(xs)

        m = (sum(xzs) - (sum(xs) * sum(zs) / N)) / (sum(xs_square) - ((sum(xs) ** 2) / N))
        assert m < 0, f"m must be negative: {m}"
        n = (sum(zs) - m * sum(xs)) / N

        self.ui.aSpinBox.setValue(math.exp(n))
        self.ui.lambdaSpinBox.setValue(-m)

    def __calc_distribution_value(self, input_value):
        a = self.ui.aSpinBox.value()
        lamb = self.ui.lambdaSpinBox.value()
        return a * math.exp(-lamb * input_value)

    def __show_tuned_mol_mass_distribution(self):
        for i in range(len(self.ui.tunedMolMassLayout)):
            item = self.ui.tunedMolMassLayout.takeAt(i)
            del item
            
        chart = BarChartWindow([(key, self.__calc_distribution_value(key)) for key in self.distro.keys()])
        self.ui.tunedMolMassLayout.addWidget(chart)