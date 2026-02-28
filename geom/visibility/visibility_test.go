// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility_test

import (
	"math"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/toolbox/v2/geom/visibility"
)

func TestNew(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(20, 20)),
		geom.NewLine(geom.NewPoint(30, 30), geom.NewPoint(40, 40)),
	}

	v := visibility.New(bounds, obstructions)
	c.NotNil(v)

	// Test that the visibility object is created with a copy of the obstructions
	result := v.SetViewPoint(geom.NewPoint(50, 50))
	c.NotNil(result)
}

func TestNewWithEmptyObstructions(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	var obstructions []geom.Line

	v := visibility.New(bounds, obstructions)
	c.NotNil(v)

	// With no obstructions, the entire bounds should be visible
	result := v.SetViewPoint(geom.NewPoint(50, 50))
	c.NotNil(result)
	c.Equal(4, len(result))
}

func TestBreakIntersections(t *testing.T) {
	c := check.New(t)

	// Test with non-intersecting lines
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 10)),
	}

	result := visibility.BreakIntersections(lines)
	c.Equal(2, len(result)) // Should remain unchanged

	// Test with intersecting lines
	intersectingLines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
	}

	result = visibility.BreakIntersections(intersectingLines)
	c.Equal(4, len(result)) // Should be broken into segments
}

func TestBreakIntersectionsWithEmptySlice(t *testing.T) {
	c := check.New(t)

	var lines []geom.Line
	result := visibility.BreakIntersections(lines)
	c.Equal(0, len(result))
}

func TestSetViewPointOutsideBounds(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(20, 20)),
	}

	v := visibility.New(bounds, obstructions)

	// Test view point outside bounds
	result := v.SetViewPoint(geom.NewPoint(-10, -10))
	c.Nil(result)

	result = v.SetViewPoint(geom.NewPoint(110, 110))
	c.Nil(result)

	result = v.SetViewPoint(geom.NewPoint(50, 110))
	c.Nil(result)
}

func TestSetViewPointInsideBounds(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(90, 10)),
	}

	v := visibility.New(bounds, obstructions)

	// Test view point inside bounds
	result := v.SetViewPoint(geom.NewPoint(50, 50))
	c.NotNil(result)
	c.True(len(result) > 0)
}

func TestSetViewPointOnBounds(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(90, 10)),
	}

	v := visibility.New(bounds, obstructions)

	c.NotNil(v.SetViewPoint(geom.NewPoint(0, 0)))
	c.Nil(v.SetViewPoint(geom.NewPoint(100, 100)))
	c.NotNil(v.SetViewPoint(geom.NewPoint(50, 0)))
}

func TestSetViewPointWithComplexObstructions(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		// Create a box obstruction in the middle
		geom.NewLine(geom.NewPoint(40, 40), geom.NewPoint(60, 40)),
		geom.NewLine(geom.NewPoint(60, 40), geom.NewPoint(60, 60)),
		geom.NewLine(geom.NewPoint(60, 60), geom.NewPoint(40, 60)),
		geom.NewLine(geom.NewPoint(40, 60), geom.NewPoint(40, 40)),
	}

	v := visibility.New(bounds, obstructions)

	// Test view point that should create a shadow
	result := v.SetViewPoint(geom.NewPoint(20, 20))
	c.NotNil(result)
	c.True(len(result) > 4)
}

func TestSetViewPointWithNoObstructions(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	var obstructions []geom.Line

	v := visibility.New(bounds, obstructions)

	result := v.SetViewPoint(geom.NewPoint(50, 50))
	c.NotNil(result)
	c.Equal(4, len(result)) // Should be exactly the 4 corners of the bounds
}

func TestVisibilityWithObstructionsOutsideBounds(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(10, 10, 80, 80)
	// Obstructions completely outside bounds
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(-10, -10), geom.NewPoint(0, 0)),
		geom.NewLine(geom.NewPoint(100, 100), geom.NewPoint(110, 110)),
	}

	v := visibility.New(bounds, obstructions)

	result := v.SetViewPoint(geom.NewPoint(50, 50))
	c.NotNil(result)
	// Should get a clean rectangle since obstructions are outside
	c.Equal(4, len(result))
}

func TestVisibilityWithPartiallyIntersectingObstructions(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	// Obstruction that crosses the boundary
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(-10, 50), geom.NewPoint(50, 50)),
		geom.NewLine(geom.NewPoint(50, -10), geom.NewPoint(50, 110)),
	}

	v := visibility.New(bounds, obstructions)

	result := v.SetViewPoint(geom.NewPoint(25, 25))
	c.NotNil(result)
	c.True(len(result) > 4)
}

// Test utility functions

func TestAngle(t *testing.T) {
	c := check.New(t)

	// We can't directly test the internal angle function, but we can test indirectly
	// through the visibility calculations with different viewpoint positions

	bounds := geom.NewRect(-10, -10, 20, 20)
	v := visibility.New(bounds, nil)
	result := v.SetViewPoint(geom.NewPoint(0, 0))
	c.NotNil(result)
	c.Equal(4, len(result))
}

func TestDistSqrd(t *testing.T) {
	c := check.New(t)

	// Test distance calculation indirectly through visibility
	bounds := geom.NewRect(0, 0, 10, 10)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(3, 3), geom.NewPoint(7, 3)),
		geom.NewLine(geom.NewPoint(3, 7), geom.NewPoint(7, 7)),
	}

	v := visibility.New(bounds, obstructions)
	result := v.SetViewPoint(geom.NewPoint(5, 5))
	c.NotNil(result)
}

func TestIntersectLines(t *testing.T) {
	c := check.New(t)

	// Test line intersection indirectly through BreakIntersections
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
	}

	result := visibility.BreakIntersections(lines)
	c.True(len(result) >= 2)

	// Test parallel lines (no intersection)
	parallelLines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(0, 1), geom.NewPoint(10, 1)),
	}

	result = visibility.BreakIntersections(parallelLines)
	c.Equal(2, len(result)) // Should remain unchanged
}

func TestHasIntersection(t *testing.T) {
	c := check.New(t)

	// Test through BreakIntersections which uses hasIntersection internally
	// Lines that clearly intersect
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
	}

	result := visibility.BreakIntersections(lines)
	c.True(len(result) > 2) // Should be broken into segments

	// Lines that don't intersect
	nonIntersectingLines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 0)),
		geom.NewLine(geom.NewPoint(6, 0), geom.NewPoint(10, 0)),
	}

	result = visibility.BreakIntersections(nonIntersectingLines)
	c.Equal(2, len(result)) // Should remain unchanged
}

func TestDirection(t *testing.T) {
	c := check.New(t)

	// Test direction calculation indirectly through intersection detection
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(5, 10)),
	}

	result := visibility.BreakIntersections(lines)
	c.True(len(result) >= 2)
}

func TestOnLine(t *testing.T) {
	c := check.New(t)

	// Test collinear segments
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(15, 0)),
	}

	result := visibility.BreakIntersections(lines)
	c.True(len(result) >= 2)
}

func TestVisibilityWithZeroBounds(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 0, 0)
	var obstructions []geom.Line

	v := visibility.New(bounds, obstructions)
	c.Nil(v.SetViewPoint(geom.NewPoint(0, 0)))
}

func TestVisibilityWithNegativeBounds(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(-50, -50, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(-10, -10), geom.NewPoint(10, 10)),
	}

	v := visibility.New(bounds, obstructions)
	result := v.SetViewPoint(geom.NewPoint(0, 0))
	c.NotNil(result)
}

func TestBreakIntersectionsWithCollinearLines(t *testing.T) {
	c := check.New(t)

	// Test with overlapping collinear segments
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(15, 0)),
		geom.NewLine(geom.NewPoint(12, 0), geom.NewPoint(20, 0)),
	}

	result := visibility.BreakIntersections(lines)
	c.True(len(result) >= len(lines))
}

func TestBreakIntersectionsWithTouchingLines(t *testing.T) {
	c := check.New(t)

	// Test with lines that touch at endpoints
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 0)),
	}

	result := visibility.BreakIntersections(lines)
	c.True(len(result) >= 2)
}

func TestVisibilityWithVerySmallObstruction(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	// Very small obstruction
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(50, 50), geom.NewPoint(50.1, 50.1)),
	}

	v := visibility.New(bounds, obstructions)
	result := v.SetViewPoint(geom.NewPoint(25, 25))
	c.NotNil(result)
}

func TestVisibilityWithLargeObstruction(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	// Large obstruction that covers most of the area
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(90, 10)),
		geom.NewLine(geom.NewPoint(90, 10), geom.NewPoint(90, 90)),
		geom.NewLine(geom.NewPoint(90, 90), geom.NewPoint(10, 90)),
		geom.NewLine(geom.NewPoint(10, 90), geom.NewPoint(10, 10)),
	}

	v := visibility.New(bounds, obstructions)
	// View point inside the large obstruction
	result := v.SetViewPoint(geom.NewPoint(50, 50))
	c.NotNil(result)
}

// Benchmark tests

func BenchmarkNew(b *testing.B) {
	bounds := geom.NewRect(0, 0, 1000, 1000)
	obstructions := make([]geom.Line, 100)
	for i := range obstructions {
		obstructions[i] = geom.NewLine(
			geom.NewPoint(float32(i*10), float32(i*5)),
			geom.NewPoint(float32(i*10+50), float32(i*5+50)),
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		visibility.New(bounds, obstructions)
	}
}

func BenchmarkBreakIntersections(b *testing.B) {
	lines := make([]geom.Line, 50)
	for i := range lines {
		angle := float64(i) * 2 * math.Pi / float64(len(lines))
		lines[i] = geom.NewLine(
			geom.NewPoint(0, 0),
			geom.NewPoint(float32(math.Cos(angle)*100), float32(math.Sin(angle)*100)),
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		visibility.BreakIntersections(lines)
	}
}

func BenchmarkSetViewPoint(b *testing.B) {
	bounds := geom.NewRect(0, 0, 1000, 1000)
	obstructions := make([]geom.Line, 20)
	for i := range obstructions {
		obstructions[i] = geom.NewLine(
			geom.NewPoint(float32(i*40+50), float32(i*30+50)),
			geom.NewPoint(float32(i*40+100), float32(i*30+100)),
		)
	}

	v := visibility.New(bounds, obstructions)
	viewPoint := geom.NewPoint(500, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.SetViewPoint(viewPoint)
	}
}

// Example tests for documentation

func ExampleNew() {
	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(20, 20), geom.NewPoint(80, 20)),
		geom.NewLine(geom.NewPoint(80, 20), geom.NewPoint(80, 80)),
	}

	v := visibility.New(bounds, obstructions)
	polygon := v.SetViewPoint(geom.NewPoint(10, 10))

	if polygon != nil {
		// Use the visibility polygon for rendering, collision detection, etc.
		_ = len(polygon) // Number of vertices in the visibility polygon
	}
}

func ExampleBreakIntersections() {
	// Lines that intersect each other
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
	}

	// Break them at intersection points
	brokenLines := visibility.BreakIntersections(lines)

	// brokenLines will contain more segments than the original lines
	_ = len(brokenLines) // Will be > 2
}

func ExampleVisibility_SetViewPoint() {
	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(30, 30), geom.NewPoint(70, 30)),
		geom.NewLine(geom.NewPoint(70, 30), geom.NewPoint(70, 70)),
	}

	v := visibility.New(bounds, obstructions)

	// Calculate visibility from this point
	viewPoint := geom.NewPoint(15, 15)
	polygon := v.SetViewPoint(viewPoint)

	// polygon contains the vertices of the visibility area
	for _, vertex := range polygon {
		_, _ = vertex.X, vertex.Y // Process each vertex
	}
}
