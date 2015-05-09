package main

import (
	. "github.com/kdar/gtest"
	"testing"
)

func TestPointInPolygon(t *testing.T) {
	boxX := []float64{0.0, 50.0, 50.0, 0.0}
	boxY := []float64{0.0, 0.0, 50.0, 50.0}
	for x, y := 1.0, 1.0; x < 50.0; x, y = x+0.1, y+0.1 {
		So(pointInPolygon(x, y, boxX, boxY), ShouldBeTrue).ElseFatal(t)
	}
}
