package poly

import "github.com/richardwilkes/toolbox/xmath/geom"

type chain struct {
	points []geom.Point
	closed bool
}

func (c *chain) linkSegment(s segment) bool {
	front := c.points[0]
	back := c.points[len(c.points)-1]
	if s.start == front {
		if s.end == back {
			c.closed = true
		} else {
			c.points = append([]geom.Point{s.end}, c.points...)
		}
		return true
	}
	if s.end == back {
		if s.start == front {
			c.closed = true
		} else {
			c.points = append(c.points, s.start)
		}
		return true
	}
	if s.end == front {
		if s.start == back {
			c.closed = true
		} else {
			c.points = append([]geom.Point{s.start}, c.points...)
		}
		return true
	}
	if s.start == back {
		if s.end == front {
			c.closed = true
		} else {
			c.points = append(c.points, s.end)
		}
		return true
	}
	return false
}

func (c *chain) linkChain(other *chain) bool {
	back := c.points[len(c.points)-1]
	otherFront := other.points[0]
	if otherFront == back {
		c.points = append(c.points, other.points[1:]...)
		other.points = nil
		return true
	}
	front := c.points[0]
	otherBack := other.points[len(other.points)-1]
	if otherBack == front {
		c.points = append(other.points, c.points[1:]...)
		other.points = nil
		return true
	}
	if otherFront == front {
		c.points = append(reverse(other.points), c.points[1:]...)
		other.points = nil
		return true
	}
	if otherBack == back {
		c.points = append(c.points[:len(c.points)-1], reverse(other.points)...)
		other.points = nil
		return true
	}
	return false
}

func reverse(list []geom.Point) []geom.Point {
	other := make([]geom.Point, len(list))
	for i := range list {
		other[len(list)-i-1] = list[i]
	}
	return other
}
