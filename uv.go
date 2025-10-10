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
