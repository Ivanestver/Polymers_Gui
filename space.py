from PyQt6.QtGui import QVector3D

class Space:
    up_vector = QVector3D(0, 1, 0)
    down_vector = -1 * up_vector
    
    left_vector = QVector3D(1, 0, 0)
    right_vector = -1 * left_vector

    forward_vector = QVector3D(0, 0, 1)
    backword_vector = QVector3D(0, 0, -1)

    global_zero = QVector3D(0, 0, 0)
    space_dimention = 0