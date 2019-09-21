package poly

import "github.com/richardwilkes/toolbox/xmath/geom"

type edgeType int

const (
	normalEdge edgeType = iota
	nonContributingEdge
	sameTransitionEdge
	differentTransitionEdge
)

type edge struct {
	pt       geom.Point
	other    *edge
	edgeType edgeType
	subject  bool
	left     bool
	inOut    bool
	inside   bool
}

func (e *edge) above(pt geom.Point) bool {
	return !e.below(pt)
}

func (e *edge) below(pt geom.Point) bool {
	if e.left {
		return signedArea(e.pt, e.other.pt, pt) > 0
	}
	return signedArea(e.other.pt, e.pt, pt) > 0
}

func (e *edge) equals(other *edge) bool {
	return *e == *other
}

func (e *edge) segment() segment {
	return segment{e.pt, e.other.pt}
}

func (e *edge) isValidDirection() bool {
	if e.left {
		return e.pt.X < e.other.pt.X || (e.pt.X == e.other.pt.X && e.pt.Y < e.other.pt.Y)
	}
	return e.other.pt.X < e.pt.X || (e.other.pt.X == e.pt.X && e.other.pt.Y < e.pt.Y)
}

func (e *edge) less(other *edge) bool {
	switch {
	case e.pt.X != other.pt.X:
		return e.pt.X > other.pt.X
	case e.pt.Y != other.pt.Y:
		return e.pt.Y > other.pt.Y
	case e.left != other.left:
		return e.left
	default:
		return e.above(other.other.pt)
	}
}

func signedArea(pt0, pt1, pt2 geom.Point) float64 {
	return (pt0.X-pt2.X)*(pt1.Y-pt2.Y) - (pt1.X-pt2.X)*(pt0.Y-pt2.Y)
}
