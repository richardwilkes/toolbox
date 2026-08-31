// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility

import (
	"cmp"
	"math"
	"slices"

	"github.com/richardwilkes/toolbox/v2/collection/quadtree"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/toolbox/v2/xmath"
)

const (
	// pointEpsilonRatio is the fraction of a scene's finest meaningful dimension used as the tolerance when comparing
	// points for equality. Deriving the tolerance from the size of the data rather than fixing it at some number of
	// world units keeps the results consistent as a scene is scaled up or down.
	pointEpsilonRatio = 0.0001
	// extentNoiseRatio is the fraction of a scene's overall extent used as a secondary floor for the point comparison
	// tolerance. All computation happens in coordinates relative to the scene's center, so the float32 rounding noise
	// a coordinate carries is proportional to the scene's extent rather than to its distance from the origin. The
	// ratio is about eight float32 ULPs.
	extentNoiseRatio = 1e-6
	// minPointEpsilon is the floor for the point comparison tolerance, used when a scene has no measurable extent.
	minPointEpsilon = 1e-6
	// minAngleEpsilon is the floor, in degrees, for the tolerance used when deciding whether two endpoints lie at the
	// same angle from the view point. It covers the quantization of float32 angles; the working tolerance comes from
	// angleEpsilon, which derives it from the point comparison tolerance so that only endpoints that could be the
	// same point within that tolerance are merged.
	minAngleEpsilon = 0.0001
	// parallelEpsilonSqrd is the square of the largest sine of the angle between two lines that still counts as
	// parallel. Squared because intersectLines can obtain the squared sine without a square root.
	parallelEpsilonSqrd = 1e-12
)

// Visibility holds state for computing a visibility polygon. A Visibility is immutable once created, so it is safe for
// concurrent use by multiple goroutines.
type Visibility struct {
	lines    []geom.Line // In scene-local coordinates.
	bounds   geom.Rect   // In scene-local coordinates.
	world    geom.Rect   // The bounds as the caller supplied them.
	offset   geom.Point  // World coordinates are scene-local coordinates plus this offset.
	epsilon  float32
	angleEps float32 // In degrees.
}

// New creates a Visibility object.
//
// bounds should not be empty. Nothing rejects one that is, but no point lies inside an empty rectangle, so every call
// to PolygonFrom on such a Visibility returns nil. A bounds with a NaN or infinite coordinate is treated as empty,
// since the four implicit edges built from its corners could never produce the finite vertices PolygonFrom promises.
//
// The obstructions must not intersect each other, which is not verified. If they do, call BreakIntersections() and
// pass the result instead.
func New(bounds geom.Rect, obstructions []geom.Line) *Visibility {
	if !finite(bounds.Point) || !finite(geom.NewPoint(bounds.Right(), bounds.Bottom())) {
		bounds = geom.Rect{}
	}
	// Everything is computed relative to the center of the bounds. Working in scene-local coordinates keeps the
	// float32 rounding noise proportional to the scene's extent rather than to its distance from the origin, so a
	// small scene far from the origin resolves exactly as well as the same scene placed at the origin.
	offset := bounds.Center()
	if !finite(offset) {
		offset = geom.Point{}
	}
	local := bounds
	local.Point = local.Point.Sub(offset)
	// The comparison tolerance tracks the smaller dimension so that an elongated bounds does not have its short
	// dimension swallowed by a tolerance derived from the long one. The noise floor stays proportional to the larger
	// dimension, which is what bounds the rounding error of the scene-local coordinates.
	epsilon := pointEpsilon(min(bounds.Width, bounds.Height), max(bounds.Width, bounds.Height))
	v := &Visibility{
		lines:    make([]geom.Line, len(obstructions)),
		bounds:   local,
		world:    bounds,
		offset:   offset,
		epsilon:  epsilon,
		angleEps: angleEpsilon(epsilon, xmath.Hypot(bounds.Width, bounds.Height)),
	}
	for i, line := range obstructions {
		v.lines[i] = geom.NewLine(line.Start.Sub(offset), line.End.Sub(offset))
	}
	return v
}

// BreakIntersections breaks the lines at their intersections, returning a new slice of lines that do not intersect.
// Collinear lines that overlap each other are split at the ends of the shared portions, and each shared portion is
// returned only once. The cut points of a collinear overlap are merged when they fall within the point comparison
// tolerance of each other, so the pieces of an overlap land on identical coordinates for every line that contains
// them, even when an overlap boundary sits within the tolerance of one line's endpoint but not another's.
//
// Lines with a NaN or infinite coordinate are dropped, since they cannot be reasoned about geometrically. Zero-length
// lines, and lines or fragments whose endpoints compare equal within the point comparison tolerance, are dropped as
// well; the tolerance is derived from the dimensions of the input's bounding box, mirroring how New derives its own
// from the bounds, and, like all point comparisons here, is applied to each axis independently. The splitting and
// merging happen in coordinates relative to the center of the input, and the intersections themselves are computed in
// float64, so input far from the origin does not lose precision to its coordinate magnitude.
func BreakIntersections(lines []geom.Line) []geom.Line {
	finiteLines := make([]geom.Line, 0, len(lines))
	minX, minY := float32(math.Inf(1)), float32(math.Inf(1))
	maxX, maxY := float32(math.Inf(-1)), float32(math.Inf(-1))
	for _, line := range lines {
		// A single NaN or infinite coordinate would poison the tolerance derivation for the entire batch, collapsing
		// or inflating the epsilon for every other line, so such lines are dropped here just as PolygonFrom drops
		// them.
		if !finite(line.Start) || !finite(line.End) {
			continue
		}
		finiteLines = append(finiteLines, line)
		minX = min(minX, line.Start.X, line.End.X)
		minY = min(minY, line.Start.Y, line.End.Y)
		maxX = max(maxX, line.Start.X, line.End.X)
		maxY = max(maxY, line.Start.Y, line.End.Y)
	}
	if len(finiteLines) == 0 {
		return finiteLines
	}
	// Mirror New's derivation from its bounds: the comparison tolerance tracks the finer dimension of the input's
	// bounding box so that an elongated batch is not preprocessed at a coarser tolerance than the sweep itself works
	// to, which would silently delete short geometry PolygonFrom could have resolved.
	spanX := maxX - minX
	spanY := maxY - minY
	eps := pointEpsilon(min(spanX, spanY), max(spanX, spanY))
	offset := geom.NewPoint((minX+maxX)/2, (minY+maxY)/2)
	if !finite(offset) {
		offset = geom.Point{}
	}
	work := make([]geom.Line, 0, len(finiteLines))
	worldLines := make([]geom.Line, 0, len(finiteLines))
	byIndex := make(map[geom.Line]int, len(finiteLines))
	var qt quadtree.QuadTree[geom.Line]
	for _, line := range finiteLines {
		local := geom.NewLine(line.Start.Sub(offset), line.End.Sub(offset))
		// A line shorter than the comparison tolerance has nothing to split and nothing to contribute as an
		// obstruction; left in, it would still cut other lines it happens to lie on.
		if local.Start.EqualWithin(local.End, eps) {
			continue
		}
		// Identical duplicate lines reduce to a single copy here, which both returns their shared portion exactly
		// once and makes every work entry unique, so a line meeting its own value among the quadtree candidates below
		// can only be meeting itself.
		if _, exists := byIndex[local]; exists {
			continue
		}
		byIndex[local] = len(work)
		work = append(work, local)
		worldLines = append(worldLines, line)
		qt.Insert(local)
	}
	// First pass: visit every intersecting pair once, recording single-point crossings as cut points and joining
	// collinear overlapping lines into groups with a union-find. Groups are handled jointly afterwards, since the
	// pieces of an overlap have to be carved identically for every line that contains them.
	parent := make([]int, len(work))
	for i := range parent {
		parent[i] = i
	}
	find := func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	cuts := make([]geom.Point, 0, len(work)*2)
	cutRanges := make([][2]int, len(work))
	grouped := false
	for i, line := range work {
		from := len(cuts)
		for _, candidate := range qt.FindIntersects(line.Bounds()) {
			if candidate == line {
				continue
			}
			j := byIndex[candidate]
			// Intersect in the caller's original coordinates: re-centering rounds each axis independently, which
			// nudges exactly collinear world input (a diagonal through points whose x and y coordinates are equal,
			// say) off exact collinearity, and the overlap of such a pair would then go undetected, leaving
			// overlapping lines in the output. The float64 arithmetic inside segmentIntersection keeps its precision
			// without the re-centering, and the cut points it returns are converted to scene-local coordinates here,
			// where the conversion's rounding sits far below the comparison tolerance.
			pts, count := segmentIntersection(worldLines[i].Start, worldLines[i].End, worldLines[j].Start,
				worldLines[j].End)
			switch count {
			case 1:
				pt := pts[0].Sub(offset)
				if !pt.EqualWithin(line.Start, eps) && !pt.EqualWithin(line.End, eps) {
					cuts = append(cuts, pt)
				}
			case 2:
				ri, rj := find(i), find(j)
				if ri != rj {
					if ri < rj {
						parent[rj] = ri
					} else {
						parent[ri] = rj
					}
					grouped = true
				}
			}
		}
		cutRanges[i] = [2]int{from, len(cuts)}
	}
	// Second pass: emit the split lines. A line that overlaps nothing is split at its crossings; each group of
	// collinear overlapping lines is emitted as the spans between consecutive merged cut points, each span exactly
	// once.
	revised := make([]geom.Line, 0, len(work)*2)
	if !grouped {
		for i, line := range work {
			revised = collectLines(line, cuts[cutRanges[i][0]:cutRanges[i][1]], revised, nil, eps)
		}
	} else {
		groups := make(map[int][]int)
		for i := range work {
			root := find(i)
			groups[root] = append(groups[root], i)
		}
		for i, line := range work {
			members := groups[find(i)]
			switch {
			case len(members) == 1:
				revised = collectLines(line, cuts[cutRanges[i][0]:cutRanges[i][1]], revised, nil, eps)
			case members[0] == i:
				revised = appendGroupSpans(work, members, cuts, cutRanges, revised, eps)
			}
		}
	}
	for i := range revised {
		revised[i] = geom.NewLine(revised[i].Start.Add(offset), revised[i].End.Add(offset))
	}
	return slices.Clip(revised)
}

// appendGroupSpans emits the pieces of one group of collinear overlapping lines. Every cut point in the group -- the
// members' endpoints and any crossings with lines outside the group -- is projected onto the group's shared line, cut
// points within the comparison tolerance of each other are merged, and each span between consecutive merged points
// that at least one member covers is appended exactly once. Merging before emission is what makes the boundary of an
// overlap land on the same coordinates for every member containing it, no matter whose endpoint it came from.
func appendGroupSpans(work []geom.Line, members []int, cuts []geom.Point, cutRanges [][2]int, revised []geom.Line,
	eps float32,
) []geom.Line {
	// Project along the direction of the longest member, which gives the most numerically stable parameterization of
	// the shared line.
	var dir geom.Point
	var bestSqrd float32
	for _, m := range members {
		d := work[m].End.Sub(work[m].Start)
		if lenSqrd := d.Dot(d); lenSqrd > bestSqrd {
			bestSqrd = lenSqrd
			dir = d
		}
	}
	dir = dir.Div(xmath.Sqrt(bestSqrd))
	origin := work[members[0]].Start
	type groupPoint struct {
		pt geom.Point
		t  float32
	}
	// The members' endpoints come first, two per member, so that the coverage step below can find member k's
	// endpoints at indices 2k and 2k+1. The crossings only add split points and never extend coverage, so their
	// positions in the slice do not matter.
	points := make([]groupPoint, 0, len(members)*3)
	for _, m := range members {
		points = append(points,
			groupPoint{pt: work[m].Start, t: work[m].Start.Sub(origin).Dot(dir)},
			groupPoint{pt: work[m].End, t: work[m].End.Sub(origin).Dot(dir)})
	}
	for _, m := range members {
		r := cutRanges[m]
		for _, pt := range cuts[r[0]:r[1]] {
			points = append(points, groupPoint{pt: pt, t: pt.Sub(origin).Dot(dir)})
		}
	}
	order := make([]int, len(points))
	for i := range order {
		order[i] = i
	}
	slices.SortFunc(order, func(a, b int) int { return cmp.Compare(points[a].t, points[b].t) })
	// Merge cut points that fall within the comparison tolerance of the first point of their cluster. Each cluster is
	// represented by that first point, so representative coordinates are always taken verbatim from an endpoint or a
	// crossing rather than re-derived, and consecutive representatives are always more than the tolerance apart,
	// which is what keeps fragments shorter than the tolerance out of the result.
	clusterOf := make([]int, len(points))
	var reps []geom.Point
	var repT []float32
	for _, idx := range order {
		if len(reps) == 0 || points[idx].t > repT[len(repT)-1]+eps {
			reps = append(reps, points[idx].pt)
			repT = append(repT, points[idx].t)
		}
		clusterOf[idx] = len(reps) - 1
	}
	if len(reps) < 2 {
		return revised
	}
	covered := make([]bool, len(reps)-1)
	for k := range members {
		a := clusterOf[2*k]
		b := clusterOf[2*k+1]
		if a > b {
			a, b = b, a
		}
		for span := a; span < b; span++ {
			covered[span] = true
		}
	}
	for span, ok := range covered {
		if ok {
			revised = append(revised, geom.NewLine(reps[span], reps[span+1]))
		}
	}
	return revised
}

// collectLines splits line at each of the given intersection points and appends the resulting segments to lines,
// returning the extended slice. Zero-length segments are dropped, as are segments that fall outside viewPort when one
// is supplied. The intersections slice is sorted in place.
func collectLines(line geom.Line, intersections []geom.Point, lines []geom.Line, viewPort *geom.Rect, eps float32,
) []geom.Line {
	// The intersection points all lie on line, so ordering them by distance from its start walks them from one end to
	// the other.
	slices.SortFunc(intersections, func(a, b geom.Point) int {
		return cmp.Compare(distSqrd(line.Start, a), distSqrd(line.Start, b))
	})
	start := line.Start
	for _, end := range intersections {
		// Three or more lines meeting at a single point yield that point once per line, so skip any repeat of the
		// current start rather than emitting a zero-length segment.
		if end.EqualWithin(start, eps) {
			continue
		}
		if viewPort == nil || withinViewPort(*viewPort, start, end) {
			lines = append(lines, geom.NewLine(start, end))
		}
		start = end
	}
	if !line.End.EqualWithin(start, eps) && (viewPort == nil || withinViewPort(*viewPort, start, line.End)) {
		lines = append(lines, geom.NewLine(start, line.End))
	}
	return lines
}

// withinViewPort reports whether the segment lies inside the view port. Callers split each segment at the crossings
// of the port's edges, snapping an endpoint that sits within the point comparison tolerance of a crossing onto it, so
// a kept fragment lies flush with the port to within rounding. The midpoint test classifies each fragment dependably,
// and appendVertex clamps any vertex derived from one back inside the bounds. Rect.IntersectsLine cannot be used
// here: it also reports true for a segment lying outside that merely touches an edge or a corner, and the sweep would
// then cast rays out to that segment's far endpoint and place polygon vertices beyond the bounds.
func withinViewPort(viewPort geom.Rect, start, end geom.Point) bool {
	midX := (start.X + end.X) / 2
	midY := (start.Y + end.Y) / 2
	return midX >= viewPort.X && midX <= viewPort.Right() && midY >= viewPort.Y && midY <= viewPort.Bottom()
}

// PolygonFrom returns the polygon of the unobstructed area visible from viewPt, or nil if viewPt is outside the bounds
// or no area is visible. Consecutive vertices are always distinct, so the polygon has no zero-length edges.
//
// An obstruction that passes within the point comparison tolerance of viewPt is ignored: a zero-thickness wall
// through or touching the eye occludes only a measure-zero set of rays, so it blocks nothing measurable.
func (v *Visibility) PolygonFrom(viewPt geom.Point) []geom.Point {
	viewPt = viewPt.Sub(v.offset)
	if !viewPt.In(v.bounds) {
		return nil
	}

	// Generate a revised line list by clipping the lines against the viewport and throwing out any that aren't within
	// the viewport.
	lines := make([]geom.Line, 0, len(v.lines)*2+4)
	viewport := [4]geom.Point{
		v.bounds.Point,
		v.bounds.TopRight(),
		v.bounds.BottomRight(),
		v.bounds.BottomLeft(),
	}
	var intersections []geom.Point
	for _, line := range v.lines {
		if !finite(line.Start) || !finite(line.End) {
			continue
		}
		// A zero-length obstruction has no extent to block anything, and its degenerate direction would divide by
		// zero in the collinear branch of segmentIntersection.
		if line.Start == line.End {
			continue
		}
		// A wall through, or within the comparison tolerance of, the eye occludes only rays along its own line, a
		// measure-zero set, so it is dropped. An endpoint at the view point also has no meaningful angle, and a ray
		// cast towards it no meaningful direction, so keeping such a wall would corrupt the ordering of every line it
		// is compared against during its batch of the sweep.
		if geom.PointSegmentDistanceSquared(line.Start, line.End, viewPt) <= v.epsilon*v.epsilon {
			continue
		}
		if (line.Start.X < v.bounds.X && line.End.X < v.bounds.X) ||
			(line.Start.Y < v.bounds.Y && line.End.Y < v.bounds.Y) ||
			(line.Start.X > v.bounds.Right() && line.End.X > v.bounds.Right()) ||
			(line.Start.Y > v.bounds.Bottom() && line.End.Y > v.bounds.Bottom()) {
			continue
		}
		intersections = intersections[:0]
		for j := range viewport {
			k := (j + 1) % len(viewport)
			pts, count := segmentIntersection(line.Start, line.End, viewport[j], viewport[k])
			for _, pt := range pts[:count] {
				// A crossing within the comparison tolerance of an endpoint is not a meaningful cut, but simply
				// discarding it would leave an endpoint that lies within the tolerance of a viewport edge poking up
				// to that far beyond it. Along the ray towards such an endpoint the edge then sits measurably in
				// front of the wall, so the sweep would keep the edge as the occluder past the event instead of
				// handing over to the wall, leaking whatever the wall should have hidden. Snapping the endpoint onto
				// the crossing keeps the fragment flush with the edge, so the two distances tie and the front is
				// chosen by the ordering's shared-point tie-break.
				switch {
				case pt.EqualWithin(line.Start, v.epsilon):
					line.Start = pt
				case pt.EqualWithin(line.End, v.epsilon):
					line.End = pt
				default:
					intersections = append(intersections, pt)
				}
			}
		}
		lines = collectLines(line, intersections, lines, &v.bounds, v.epsilon)
	}

	lines = append(lines,
		geom.NewLine(v.bounds.Point, v.bounds.TopRight()),
		geom.NewLine(v.bounds.TopRight(), v.bounds.BottomRight()),
		geom.NewLine(v.bounds.BottomRight(), v.bounds.BottomLeft()),
		geom.NewLine(v.bounds.BottomLeft(), v.bounds.Point),
	)

	return v.toWorld(v.computePolygon(viewPt, lines))
}

// toWorld translates a polygon from scene-local coordinates back to world coordinates, clamping each vertex to the
// bounds. Adding the offset back rounds to the resolution the world coordinates support, which for a scene far from
// the origin is coarser than the scene-local tolerance, so consecutive vertices that collapse onto the same world
// coordinates are reduced to one and the documented at-least-three-distinct-vertices contract is re-checked.
func (v *Visibility) toWorld(polygon []geom.Point) []geom.Point {
	out := polygon[:0]
	for _, pt := range polygon {
		p := pt.Add(v.offset)
		p.X = min(max(p.X, v.world.X), v.world.Right())
		p.Y = min(max(p.Y, v.world.Y), v.world.Bottom())
		if len(out) != 0 && out[len(out)-1] == p {
			continue
		}
		out = append(out, p)
	}
	for len(out) > 1 && out[0] == out[len(out)-1] {
		out = out[:len(out)-1]
	}
	if len(out) < 3 {
		return nil
	}
	return out
}

func (v *Visibility) computePolygon(viewPt geom.Point, lines []geom.Line) []geom.Point {
	// Sweep through the points to generate the visibility polygon
	sorted := sortLines(viewPt, lines)
	heap := newLineHeap(lines, viewPt, v.epsilon)
	// The sweep starts along the +X direction from the view point. The offset only fixes a direction, but it has to
	// be large enough that adding it to the view point's coordinate actually changes it: a fixed +1 vanishes below
	// the float32 resolution of a large enough coordinate, collapsing the ray to a point and with it the ordering of
	// the entire seeded heap, so it is derived from the scene's extent instead.
	start := geom.NewPoint(viewPt.X+max(1, v.bounds.Width, v.bounds.Height), viewPt.Y)
	for i := range lines {
		a1 := angle(lines[i].Start, viewPt)
		a2 := angle(lines[i].End, viewPt)
		// Seed the heap with the lines that straddle the ±180° discontinuity the sweep starts from. angle returns an
		// atan2 result, so both angles already lie within [-180,180]; all that is left to test is that they fall on
		// opposite sides of zero and are more than 180° apart.
		if (a1 <= 0 && a2 >= 0 && a2-a1 > 180) || (a2 <= 0 && a1 >= 0 && a1-a2 > 180) {
			heap.insert(i, start)
		}
	}
	// The sweep emits at most two points per batch of endpoints that share an angle, but in practice lands a little
	// above one per endpoint, so size the polygon for that rather than for the worst case.
	polygon := make([]geom.Point, 0, len(sorted)+len(sorted)/2)
	i := 0
	for i < len(sorted) {
		orig := i
		for i++; i < len(sorted) && sorted[i].angle < sorted[orig].angle+v.angleEps; i++ {
		}
		extend := false
		shorten := false
		vertex := sorted[orig].pt(lines)
		oldLine := heap.nearest()
		// The batch's removals happen before its insertions. A line ending at a point that another line passes
		// through or starts from swaps its order against that line exactly at this batch's angle, and the heap is
		// only ever reordered by the operations themselves, so the ending line has to leave before a newcomer is
		// sifted in: a sift compares the newcomer against only its ancestors and relies on every other pair already
		// being consistently ordered, which such a stale pair no longer is, and a newcomer can otherwise end up
		// parked in front of a line that actually occludes it.
		for j := orig; j < i; j++ {
			if heap.contains(sorted[j].lineIndex) {
				if sorted[j].lineIndex == oldLine {
					extend = true
					vertex = sorted[j].pt(lines)
				}
				heap.remove(sorted[j].lineIndex, vertex)
				sorted[j].handled = true
			}
		}
		for j := orig; j < i; j++ {
			if sorted[j].handled {
				continue
			}
			if heap.contains(sorted[j].lineIndex) {
				// The line's other endpoint fell within this same batch: it entered moments ago at that event and
				// leaves again here without ever spanning a full batch.
				heap.remove(sorted[j].lineIndex, vertex)
				continue
			}
			heap.insert(sorted[j].lineIndex, vertex)
			if heap.nearest() != oldLine {
				shorten = true
			}
		}
		// The heap can be emptied by the removals above, so every nearest-occluder lookup has to tolerate the empty
		// case rather than indexing blindly.
		if extend {
			polygon = v.appendVertex(polygon, vertex)
			if nearest := heap.nearest(); nearest != -1 {
				line := lines[nearest]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects &&
					!cur.EqualWithin(vertex, v.epsilon) {
					polygon = v.appendVertex(polygon, cur)
				}
			}
		} else if shorten {
			if oldLine != -1 {
				line := lines[oldLine]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects {
					polygon = v.appendVertex(polygon, cur)
				}
			}
			if nearest := heap.nearest(); nearest != -1 {
				line := lines[nearest]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects {
					polygon = v.appendVertex(polygon, cur)
				}
			}
		}
	}
	// The sweep wraps all the way around, so it can finish on the vertex it started from. Drop that repeat as well,
	// since the closing edge is implicit.
	if len(polygon) > 1 && polygon[0].EqualWithin(polygon[len(polygon)-1], v.epsilon) {
		polygon = polygon[:len(polygon)-1]
	}
	// A polygon needs at least three vertices to enclose anything. A scene whose extent is at the scale of the
	// comparison tolerance, such as a sliver bounds far thinner than it is long, can collapse the sweep's output to
	// one or two points, and the documented contract is a real polygon or nil, never a degenerate point or edge.
	if len(polygon) < 3 {
		return nil
	}
	return polygon
}

// appendVertex clamps pt to the bounds and appends it to the polygon, unless it repeats the vertex already at the end.
//
// The clamp holds the documented postcondition that the visible area lies within the bounds. The sweep intersects the
// ray with the infinite line through the nearest occluder, so when that occluder is a bounds edge and the ray runs
// almost parallel to it -- which happens near a corner of a very elongated rectangle -- the crossing can land well
// beyond the edge's end.
//
// Emitting a repeated vertex would leave a zero-length edge for every consumer to filter out, and would break polygon
// algorithms that assume each edge has a direction.
func (v *Visibility) appendVertex(polygon []geom.Point, pt geom.Point) []geom.Point {
	pt.X = min(max(pt.X, v.bounds.X), v.bounds.Right())
	pt.Y = min(max(pt.Y, v.bounds.Y), v.bounds.Bottom())
	if len(polygon) != 0 && polygon[len(polygon)-1].EqualWithin(pt, v.epsilon) {
		return polygon
	}
	return append(polygon, pt)
}

// lineHeap is a min-heap of indices into lines, ordered by how close each line is to viewPt along a ray towards a
// destination supplied per operation. order holds line indices keyed by heap position, and slot holds heap positions
// keyed by line index; the two are inverses of each other, which is what lets a line be located in the heap without
// searching.
type lineHeap struct {
	lines   []geom.Line
	order   []int
	slot    []int
	viewPt  geom.Point
	epsilon float32
}

// newLineHeap returns an empty heap over the given lines. Both index slices come from one allocation, since neither
// can ever hold more than one entry per line.
func newLineHeap(lines []geom.Line, viewPt geom.Point, epsilon float32) lineHeap {
	backing := make([]int, len(lines)*2)
	h := lineHeap{
		lines:   lines,
		order:   backing[:0:len(lines)],
		slot:    backing[len(lines):],
		viewPt:  viewPt,
		epsilon: epsilon,
	}
	for i := range h.slot {
		h.slot[i] = -1
	}
	return h
}

// contains reports whether the given line is currently in the heap.
func (h *lineHeap) contains(lineIndex int) bool {
	return h.slot[lineIndex] != -1
}

// nearest returns the index of the line closest to the view point, or -1 if the heap is empty.
func (h *lineHeap) nearest() int {
	if len(h.order) == 0 {
		return -1
	}
	return h.order[0]
}

// insert adds the given line to the heap, ordering it against the ray from the view point towards destination.
func (h *lineHeap) insert(lineIndex int, destination geom.Point) {
	cur := len(h.order)
	h.order = append(h.order, lineIndex)
	h.slot[lineIndex] = cur
	h.siftUp(cur, destination)
}

// remove takes the given line back out of the heap, reordering what is left against the ray from the view point
// towards destination.
func (h *lineHeap) remove(lineIndex int, destination geom.Point) {
	cur := h.slot[lineIndex]
	h.slot[lineIndex] = -1
	last := len(h.order) - 1
	if cur != last {
		h.order[cur] = h.order[last]
		h.slot[h.order[cur]] = cur
	}
	h.order = h.order[:last]
	if cur != last && !h.siftUp(cur, destination) {
		h.siftDown(cur, destination)
	}
}

// swap exchanges the entries at two heap positions, keeping slot in step with order.
func (h *lineHeap) swap(a, b int) {
	h.slot[h.order[a]], h.slot[h.order[b]] = b, a
	h.order[a], h.order[b] = h.order[b], h.order[a]
}

// siftUp moves the entry at the given heap position towards the front while it is closer than its parent, reporting
// whether it moved at all.
func (h *lineHeap) siftUp(cur int, destination geom.Point) bool {
	moved := false
	for cur > 0 {
		parent := (cur - 1) / 2
		if !h.lessThan(h.order[cur], h.order[parent], destination) {
			break
		}
		h.swap(cur, parent)
		cur = parent
		moved = true
	}
	return moved
}

// siftDown moves the entry at the given heap position towards the back while either child is closer than it is.
func (h *lineHeap) siftDown(cur int, destination geom.Point) {
	for {
		left := 2*cur + 1
		if left >= len(h.order) {
			return
		}
		closest := left
		if right := left + 1; right < len(h.order) && h.lessThan(h.order[right], h.order[left], destination) {
			closest = right
		}
		if !h.lessThan(h.order[closest], h.order[cur], destination) {
			return
		}
		h.swap(cur, closest)
		cur = closest
	}
}

// lessThan reports whether the line at index1 is closer to the view point along the ray towards destination than the
// line at index2. A line that the ray misses entirely sorts as infinitely far away, so that it can always be
// displaced from the front of the heap by one that is actually hit, and two missed lines are equivalent.
//
// The comparison is a genuine strict weak ordering, which the sift-up and sift-down code depends on: distances are
// compared by the tolerance-sized bucket they fall in rather than gated by pairwise epsilon equality, which is
// non-transitive and admits ordering cycles among lines whose crossings sit just under a tolerance apart. Bucket ties
// -- lines sharing a point on the current ray, such as a clipped wall endpoint lying on a bounds edge or several
// walls meeting at a cut -- are broken by frontKey, a totally ordered scalar, with exactly collinear lines untangled
// by length.
func (h *lineHeap) lessThan(index1, index2 int, destination geom.Point) bool {
	pt1, intersects1 := intersectLines(h.lines[index1].Start, h.lines[index1].End, h.viewPt, destination)
	pt2, intersects2 := intersectLines(h.lines[index2].Start, h.lines[index2].End, h.viewPt, destination)
	if !intersects1 {
		return false
	}
	if !intersects2 {
		return true
	}
	b1 := xmath.Floor(xmath.Sqrt(distSqrd(pt1, h.viewPt)) / h.epsilon)
	b2 := xmath.Floor(xmath.Sqrt(distSqrd(pt2, h.viewPt)) / h.epsilon)
	if b1 != b2 {
		return b1 < b2
	}
	k1 := h.frontKey(index1, pt1)
	k2 := h.frontKey(index2, pt2)
	if k1 != k2 {
		return k1 < k2
	}
	// Exactly collinear lines have equal keys and are interchangeable as occluders. Putting the shorter one in front
	// keeps its removal event visible to the sweep, so the endpoints of a wall lying along a bounds edge still become
	// vertices; two lines equal in both key and length are genuinely equivalent.
	return distSqrd(h.lines[index1].Start, h.lines[index1].End) < distSqrd(h.lines[index2].Start, h.lines[index2].End)
}

// frontKey returns the rate at which the crossing distance of the line at the given index changes as the sweep's ray
// rotates forward past its crossing point pt, up to a positive factor shared by every line crossing the ray at that
// point: the cotangent of the angle between the line and the direction from pt back to the view point. Two lines
// crossing the current ray at the same point are ordered by which one bends in front of the other just beyond that
// point, which is the one whose distance shrinks faster or grows more slowly, so the smaller key is the nearer line.
//
// The cotangent is periodic in the line's direction, so the key does not depend on which way the line's endpoints
// happen to be ordered. The angle-based tie-break this replaces did depend on it -- it picked a direction from
// whichever endpoint the crossing was not near, which is ambiguous for a line crossed mid-span -- and its two-class
// comparison key was split exactly at the angle that an eye lying close to the line through a wall or bounds edge
// produces, so the ordering there flipped on rounding noise, leaking scene-scale regions past the wall or carving
// them away.
func (h *lineHeap) frontKey(index int, pt geom.Point) float64 {
	ux := float64(h.lines[index].End.X) - float64(h.lines[index].Start.X)
	uy := float64(h.lines[index].End.Y) - float64(h.lines[index].Start.Y)
	ex := float64(h.viewPt.X) - float64(pt.X)
	ey := float64(h.viewPt.Y) - float64(pt.Y)
	cross := ex*uy - ey*ux
	if cross == 0 {
		// The line runs along the ray itself, which intersectLines rejects as parallel before the ordering ever
		// compares it, so this is unreachable in practice; +Inf orders such a line as maximally behind if the guard is
		// ever reached.
		return math.Inf(1)
	}
	return (ux*ex + uy*ey) / cross
}

func sortLines(position geom.Point, lines []geom.Line) []endPoint {
	points := make([]endPoint, len(lines)*2)
	pos := 0
	for i, line := range lines {
		points[pos].lineIndex = i
		points[pos].angle = angle(line.Start, position)
		points[pos].start = true
		pos++
		points[pos].lineIndex = i
		points[pos].angle = angle(line.End, position)
		points[pos].start = false
		pos++
	}
	slices.SortFunc(points, func(a, b endPoint) int {
		if result := cmp.Compare(a.angle, b.angle); result != 0 {
			return result
		}
		if result := cmp.Compare(distSqrd(a.pt(lines), position), distSqrd(b.pt(lines), position)); result != 0 {
			return result
		}
		if a.start == b.start {
			return 0
		}
		if a.start {
			return 1
		}
		return -1
	})
	return points
}

func angle(a, b geom.Point) float32 {
	return xmath.Atan2(b.Y-a.Y, b.X-a.X) * 180 / math.Pi
}

// pointEpsilon returns the tolerance to use when comparing points for equality in a scene whose finest meaningful
// dimension is feature and whose overall extent is extent. The primary term is proportional to the feature size,
// which keeps behavior the same at every scale. The extent term keeps the tolerance above the float32 rounding noise
// carried by scene-local coordinates, which is proportional to the scene's extent, and the fixed floor keeps the
// tolerance from collapsing to zero for a scene with no measurable extent.
func pointEpsilon(feature, extent float32) float32 {
	return max(feature*pointEpsilonRatio, extent*extentNoiseRatio, minPointEpsilon)
}

// angleEpsilon returns the tolerance, in degrees, for treating two endpoints as lying at the same angle from the view
// point. It is the angle the point comparison tolerance subtends across the scene's diagonal, so two endpoints merge
// into one batch of the sweep only when they could be the same point to within that tolerance. A fixed angular
// tolerance instead merges a wall endpoint with a bounds corner it merely aligns with, silently deleting or inflating
// the wall's shadow. The floor covers the quantization of float32 angles.
func angleEpsilon(epsilon, diagonal float32) float32 {
	if diagonal <= 0 || xmath.IsInf(diagonal, 0) {
		return minAngleEpsilon
	}
	return max(epsilon/diagonal*(180/math.Pi), minAngleEpsilon)
}

// finite reports whether both of a point's coordinates can take part in the sweep. A NaN or infinite coordinate has
// neither a meaningful angle nor a meaningful distance, so a line carrying one is dropped before it reaches anything
// that would have to reason about it. The viewport clipping happens to reject such a line as well, since every
// comparison against a NaN is false, but relying on that would leave the behavior resting on an accident.
func finite(pt geom.Point) bool {
	return !xmath.IsNaN(pt.X) && !xmath.IsNaN(pt.Y) && !xmath.IsInf(pt.X, 0) && !xmath.IsInf(pt.Y, 0)
}

func distSqrd(a, b geom.Point) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

// intersectLines returns the point where the infinite lines through the two segments cross. It reports false when the
// lines are parallel, or so nearly parallel that the crossing point would be far enough away to swamp every distance
// it was compared against. ub is the cross product of the two direction vectors, so comparing its square against the
// product of their squared lengths tests the squared sine of the angle between them, which is independent of scale.
//
// The arithmetic is done in float64. The heap comparator buckets the crossing distances it gets from here by the
// comparison tolerance, and when two lines share a point on the current ray -- a clipped wall endpoint lying on a
// bounds edge, or several walls meeting at a cut -- their equal distances must land in the same bucket for the
// angle-based tie-break to engage at all. Float32 rounding here is a meaningful fraction of the bucket width, so it
// randomly separated such ties across a bucket boundary, leaving whichever line the noise favored in front.
func intersectLines(s1, e1, s2, e2 geom.Point) (geom.Point, bool) {
	dbx := float64(e2.X) - float64(s2.X)
	dby := float64(e2.Y) - float64(s2.Y)
	dax := float64(e1.X) - float64(s1.X)
	day := float64(e1.Y) - float64(s1.Y)
	ub := dby*dax - dbx*day
	if ub*ub <= (parallelEpsilonSqrd*(dax*dax+day*day))*(dbx*dbx+dby*dby) {
		return geom.Point{}, false
	}
	ua := (dbx*(float64(s1.Y)-float64(s2.Y)) - dby*(float64(s1.X)-float64(s2.X))) / ub
	return geom.Point{X: float32(float64(s1.X) + ua*dax), Y: float32(float64(s1.Y) + ua*day)}, true
}

// segmentIntersection returns the intersection of two segments, if any, without allocating. A count of 0 means no
// intersection, 1 a single shared point (in pts[0]), and 2 a collinear overlapping span, with pts holding the ends of
// the shared portion interpolated along the first segment. It mirrors geom.LineIntersection for segments of nonzero
// length, which callers guarantee.
//
// The arithmetic is done in float64, for two reasons. First, the products of float32 differences carry at most 48
// significant bits, so they are exact in float64 and the exact-zero collinearity tests below become reliable; in
// float32, fused multiply-subtract contraction of a*b - c*d (which the compiler emits on arm64) returns the rounding
// error of the second product instead of zero for exactly collinear segments, sending every such overlap down the
// "not parallel" branch and leaving it undetected. Second, an interpolated crossing lands within a float32 rounding
// of the true crossing even when one segment is far longer than the other, where float32 interpolation drifts by the
// long segment's length times the float32 rounding unit -- for a caller-supplied obstruction clipped against the
// viewport, far enough past the viewport edge to defeat the sweep's tolerance. The explicit float64 conversions
// around each product keep the compiler from fusing even the float64 subtractions, which matters for differences that
// do not happen to be exact.
func segmentIntersection(a1, a2, b1, b2 geom.Point) (pts [2]geom.Point, count int) {
	abdx := float64(a1.X) - float64(b1.X)
	abdy := float64(a1.Y) - float64(b1.Y)
	bdx := float64(b2.X) - float64(b1.X)
	bdy := float64(b2.Y) - float64(b1.Y)
	uat := float64(bdx*abdy) - float64(bdy*abdx)
	adx := float64(a2.X) - float64(a1.X)
	ady := float64(a2.Y) - float64(a1.Y)
	ubt := float64(adx*abdy) - float64(ady*abdx)
	ub := float64(bdy*adx) - float64(bdx*ady)
	if ub != 0 {
		// Not parallel, so there is at most a single crossing.
		a := uat / ub
		if a >= 0 && a <= 1 {
			if b := ubt / ub; b >= 0 && b <= 1 {
				pts[0] = geom.Point{X: float32(float64(a1.X) + a*adx), Y: float32(float64(a1.Y) + a*ady)}
				count = 1
			}
		}
		return pts, count
	}
	// Parallel. Collinearity requires both cross-product numerators to be zero; in exact arithmetic either being zero
	// implies the other, but requiring both guards against a phantom overlap when rounding zeroes only one of them
	// for parallel-but-offset segments.
	if uat != 0 || ubt != 0 {
		return pts, 0
	}
	var ub1, ub2 float64
	if math.Abs(adx) > math.Abs(ady) {
		ub1 = (float64(b1.X) - float64(a1.X)) / adx
		ub2 = (float64(b2.X) - float64(a1.X)) / adx
	} else {
		ub1 = (float64(b1.Y) - float64(a1.Y)) / ady
		ub2 = (float64(b2.Y) - float64(a1.Y)) / ady
	}
	left := max(0, min(ub1, ub2))
	right := min(1, max(ub1, ub2))
	if left > right {
		return pts, 0
	}
	pts[0] = geom.Point{
		X: float32(float64(a2.X)*left + float64(a1.X)*(1-left)),
		Y: float32(float64(a2.Y)*left + float64(a1.Y)*(1-left)),
	}
	if left == right {
		return pts, 1
	}
	pts[1] = geom.Point{
		X: float32(float64(a2.X)*right + float64(a1.X)*(1-right)),
		Y: float32(float64(a2.Y)*right + float64(a1.Y)*(1-right)),
	}
	return pts, 2
}

type endPoint struct {
	lineIndex int
	angle     float32
	start     bool
	// handled marks an event the sweep's removal pass has already consumed, so the insertion pass skips it.
	handled bool
}

func (ep *endPoint) pt(lines []geom.Line) geom.Point {
	if ep.start {
		return lines[ep.lineIndex].Start
	}
	return lines[ep.lineIndex].End
}
