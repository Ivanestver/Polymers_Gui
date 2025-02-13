from PyQt6.QtCharts import QChartView, QChart, QLineSeries, QBarSet, QBarSeries, QBarCategoryAxis, QValueAxis
from PyQt6.QtGui import QPainter
from PyQt6.QtCore import Qt

class LineChartWindow(QChartView):
    def __init__(self, points: list[tuple[int, float]], title: str = ""):
        super().__init__()
        self.__chart = QChart()
        self.setChart(self.__chart)
        self.__chart.setTitle(title)

        series = QLineSeries()
        for point in points:
            series.append(*point)
        self.__chart.addSeries(series)
        self.__chart.createDefaultAxes()
        self.setRenderHint(QPainter.RenderHint.Antialiasing)

class BarChartWindow(QChartView):
    def __init__(self, points: list[tuple[int, float]], title: str = ""):
        super().__init__()
        self.__chart = QChart()
        self.setChart(self.__chart)
        self.__chart.setTitle(title)

        # Добавляем точки
        barSet = QBarSet(title)
        y_values = list([point[1] for point in points])
        barSet.append(y_values)

        # Создаём series
        series = QBarSeries()
        series.append(barSet)
        self.__chart.addSeries(series)

        axisX = QBarCategoryAxis()
        x_values = list([str(point[0]) for point in points])
        axisX.append(x_values)
        self.__chart.addAxis(axisX, Qt.AlignmentFlag.AlignBottom)
        series.attachAxis(axisX)

        # Добавляем ось Y
        axisY = QValueAxis()
        y_values.append(0)
        axisY.setRange(min(y_values), max(y_values))
        self.__chart.addAxis(axisY, Qt.AlignmentFlag.AlignLeft)
        series.attachAxis(axisY)
            
        self.setRenderHint(QPainter.RenderHint.Antialiasing)