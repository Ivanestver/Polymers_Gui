from PyQt6.QtGui import QVector3D

class Space:
    up_vector = QVector3D(0, 1, 0)
    down_vector = -1 * up_vector
    
    left_vector = QVector3D(1, 0, 0)
    right_vector = -1 * left_vector

    forward_vector = QVector3D(0, 0, 1)
    backword_vector = QVector3D(0, 0, -1)

    space_center = QVector3D(0, 0, 0)
    space_dimention = 0
    
    def point_within_borders(point: tuple):
        return all([0 <= c <= Space.space_dimention for c in point])

    def set_space_dimention(space_dimention):
        Space.space_dimention = space_dimention
        Space.space_center = QVector3D(space_dimention / 2, space_dimention / 2, space_dimention / 2)
