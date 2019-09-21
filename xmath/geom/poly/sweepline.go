package poly

type sweep []*edge

func (s *sweep) insert(item *edge) int {
	length := len(*s)
	if length == 0 {
		*s = append(*s, item)
		return 0
	}

	*s = append(*s, &edge{})
	i := length - 1
	for i >= 0 && compareSegment(item, (*s)[i]) {
		(*s)[i+1] = (*s)[i]
		i--
	}
	(*s)[i+1] = item
	return i + 1
}

func (s *sweep) remove(e *edge) {
	for i, el := range *s {
		if el.equals(e) {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}

func compareSegment(e1, e2 *edge) bool {
	switch {
	case e1 == e2:
		return false
	case signedArea(e1.pt, e1.other.pt, e2.pt) != 0 || signedArea(e1.pt, e1.other.pt, e2.other.pt) != 0:
		if e1.pt == e2.pt {
			return e1.below(e2.other.pt)
		}
		if e1.less(e2) {
			return e2.above(e1.pt)
		}
		return e1.below(e2.pt)
	case e1.pt == e2.pt:
		return false
	default:
		return e1.less(e2)
	}
}
