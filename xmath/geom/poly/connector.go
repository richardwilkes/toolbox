package poly

import "github.com/richardwilkes/toolbox/xmath/geom"

type connector struct {
	open   []chain
	closed []chain
	op     op
}

func (c *connector) add(s segment) {
	for i := range c.open {
		chn := &c.open[i]
		if !chn.linkSegment(s) {
			continue
		}
		if chn.closed {
			if len(chn.points) == 2 {
				chn.closed = false
				return
			}
			c.closed = append(c.closed, c.open[i])
			c.open = append(c.open[:i], c.open[i+1:]...)
			return
		}
		k := len(c.open)
		for j := i + 1; j < k; j++ {
			if chn.linkChain(&c.open[j]) {
				c.open = append(c.open[:j], c.open[j+1:]...)
				return
			}
		}
		return
	}
	c.open = append(c.open, chain{points: []geom.Point{s.start, s.end}})
}

func (c *connector) toPolygon() Polygon {
	var chn []chain
	if c.op == clipLineOp {
		chn = c.open
	} else {
		chn = c.closed
	}
	poly := Polygon{}
	for i := range chn {
		poly.Add(append(Contour{}, chn[i].points...))
	}
	return poly
}
