package shog

type UV struct {
	X int
	Y int
}

func NewUV(x, y int) UV {
	return UV{
		X: x,
		Y: y,
	}
}

func (uv *UV) Zero() bool {
	return uv.X == 0 && uv.Y == 0
}

func (uv *UV) Square() int {
	return uv.X * uv.Y
}
