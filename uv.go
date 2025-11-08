package shog

type UV struct {
	X      int
	Y      int
	Square int
}

func NewUV(x, y int) UV {
	return UV{
		X:      x,
		Y:      y,
		Square: x * y,
	}
}

func (uv *UV) Zero() bool {
	return uv.X == 0 && uv.Y == 0
}
