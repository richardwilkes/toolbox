package poly

import (
	"sort"
)

type queue struct {
	edges  []*edge
	sorted bool
}

func (q *queue) more() bool {
	return len(q.edges) > 0
}

func (q *queue) enqueue(e *edge) {
	if !q.sorted || len(q.edges) == 0 {
		q.edges = append(q.edges, e)
		return
	}
	i := len(q.edges) - 1
	q.edges = append(q.edges, nil)
	for i >= 0 && e.less(q.edges[i]) {
		q.edges[i+1] = q.edges[i]
		i--
	}
	q.edges[i+1] = e
}

func (q *queue) dequeue() *edge {
	if !q.sorted {
		sort.Slice(q.edges, func(i, j int) bool {
			return q.edges[i].less(q.edges[j])
		})
		q.sorted = true
	}
	x := q.edges[len(q.edges)-1]
	q.edges = q.edges[:len(q.edges)-1]
	return x
}
