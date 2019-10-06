package poly

type scanBeamTree struct {
	root    *scanBeamNode
	entries int
}

type scanBeamNode struct {
	y    float64
	less *scanBeamNode
	more *scanBeamNode
}

func (sbt *scanBeamTree) add(y float64) {
	sbt.addToScanBeamTreeAt(&sbt.root, y)
}

func (sbt *scanBeamTree) addToScanBeamTreeAt(node **scanBeamNode, y float64) {
	switch {
	case *node == nil:
		*node = &scanBeamNode{y: y}
		sbt.entries++
	case (*node).y > y:
		sbt.addToScanBeamTreeAt(&(*node).less, y)
	case (*node).y < y:
		sbt.addToScanBeamTreeAt(&(*node).more, y)
	default:
	}
}

func (sbt *scanBeamTree) buildScanBeamTable() []float64 {
	table := make([]float64, sbt.entries)
	if sbt.root != nil {
		sbt.root.buildScanBeamTableEntries(0, table)
	}
	return table
}

func (sbn *scanBeamNode) buildScanBeamTableEntries(index int, table []float64) int {
	if sbn.less != nil {
		index = sbn.less.buildScanBeamTableEntries(index, table)
	}
	table[index] = sbn.y
	index++
	if sbn.more != nil {
		index = sbn.more.buildScanBeamTableEntries(index, table)
	}
	return index
}
