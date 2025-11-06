package shog

type PannelBorder struct {
	Horizontal  Symbol
	Virtical    Symbol
	TopLeft     Symbol
	TopRight    Symbol
	BottomLeft  Symbol
	BottomRight Symbol
}

type PannelBorderOption func(pb *PannelBorder)

func NewPannelBorder(options ...PannelBorderOption) PannelBorder {
	pb := PannelBorder{
		Horizontal:  B_H,
		Virtical:    B_V,
		TopLeft:     B_TLEFT,
		TopRight:    B_TRIGHT,
		BottomLeft:  B_BLEFT,
		BottomRight: B_BRIGHT,
	}
	for _, option := range options {
		option(&pb)
	}
	return pb
}

func WithCustomPannelBorder(h, v, tl, tr, bl, br Symbol) PannelBorderOption {
	return func(pb *PannelBorder) {
		pb.Horizontal = h
		pb.Virtical = v
		pb.TopLeft = tl
		pb.TopRight = tr
		pb.BottomLeft = bl
		pb.BottomRight = br
	}
}
