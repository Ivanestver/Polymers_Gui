from PyQt6.QtCharts import QChartView, QChart, QLineSeries
from PyQt6.QtGui import QPainter

class ChartWindow(QChartView):
    def __init__(self, points: list[tuple[int, float]]):
        super().__init__()
        self.__chart = QChart()
        self.setChart(self.__chart)

        series = QLineSeries()
        for point in points:
            series.append(*point)
        self.__chart.addSeries(series)
        self.__chart.createDefaultAxes()
        self.setRenderHint(QPainter.RenderHint.Antialiasing)

    def __setup_legend(self):
        pass