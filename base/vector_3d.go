package base

type Vector3D struct {
	X int64
	Y int64
	Z int64
}

func InvalidVector() Vector3D {
	return Vector3D{
		X: -1,
		Y: -1,
		Z: -1,
	}
}

func (this *Vector3D) IsInvalid() bool {
	return this.X == -1 && this.Y == -1 && this.Z == -1
}

func VectorsAreEqual(left, right *Vector3D) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return left.X == right.X &&
		left.Y == right.Y &&
		left.Z == right.Z
}

func (vector *Vector3D) Add(other *Vector3D) {
	vector.X += other.X
	vector.Y += other.Y
	vector.Z += other.Z
}

func AddVec(left *Vector3D, right *Vector3D) *Vector3D {
	return &Vector3D{
		X: left.X + right.X,
		Y: left.Y + right.Y,
		Z: left.Z + right.Z,
	}
}
