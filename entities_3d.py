from PyQt6.QtGui import QColor, QVector3D

AtomColors = {
    "C": QColor.black()
}

def get_color_by_label(label: str):
    if label not in AtomColors.keys():
        raise Exception(f"There is no such a label: {label}")
    return AtomColors.get(label)

class Atom:
    def __init__(self, radius: float, label: str, position: QVector3D):
        self._radius = radius
        self._label = label
        self._position = position
    
    @property
    def radius(self):
        return self._radius

    @property
    def label(self):
        return self.label

    @property
    def position(self):
        return self._position
    
    @property
    def color(self):
        try:
            return get_color_by_label(self._label)
        except Exception as err:
            print(err.args)
            return QColor.black()

class Polymer:
    def __init__(self, name: str = ""):
        self._name = name
        self._atoms = list[Atom]()
    
    @property
    def name(self):
        return self._name

    def len(self):
        return len(self._atoms)