package kk

import (
	"image"
	"runtime"

	"golang.org/x/exp/shiny/unit"
	"golang.org/x/exp/shiny/widget"
	"golang.org/x/exp/shiny/widget/node"
	"golang.org/x/exp/shiny/widget/theme"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/geom"
)

/*
 * This file contains all we need to do the layout on the screen.
 * We use widgets from shiny to do this, but that is probably a bad
 * idea because the interfaces used might not be inentionally exposed
 * and probably won't stay stable. It works for now, but it might be
 * a good idea to eventually write some simple package that lays out
 * boxes.
 *
 * Originally I also wanted to use shiny widgets for event routing but
 * that isn't usable yet.
 */

// helper.
func stretch(n node.Node, alongWeight int) node.Node {
	return widget.WithLayoutData(n, widget.FlowLayoutData{
		AlongWeight:  alongWeight,
		ExpandAlong:  true,
		ShrinkAlong:  true,
		ExpandAcross: true,
		ShrinkAcross: true,
	})
}

// helper to extract screen coords for a widget.
func widgetScreenRect(n node.Node) image.Rectangle {
	e := n.Wrappee()
	r := e.Rect
	for e.Parent != nil {
		e = e.Parent
		r = r.Add(e.Rect.Min)
	}
	return r
}

// helper to translate image points to geom points for glutil.Images.
func (s *State) ip2gp(ip image.Point) geom.Point {
	return geom.Point{
		X: geom.Pt(ip.X) / geom.Pt(s.wsz.PixelsPerPt),
		Y: geom.Pt(ip.Y) / geom.Pt(s.wsz.PixelsPerPt),
	}
}

func (s *State) setSize(e size.Event) {
	s.wsz = e
	t := theme.Theme{DPI: float64(unit.PointsPerInch * e.PixelsPerPt)}

	horizontal := e.Orientation == size.OrientationLandscape || e.WidthPx > e.HeightPx

	padPx := float64(e.WidthPx / 50)
	butPx := float64(e.WidthPx / 6)
	bAx := widget.AxisHorizontal
	aAx := widget.AxisVertical
	if horizontal {
		padPx = float64(e.HeightPx / 50)
		butPx = float64(e.HeightPx / 6)
		bAx = widget.AxisVertical
		aAx = widget.AxisHorizontal
	}

	// We abuse shiny widgets to do the layout for us.

	bw := make(map[string]node.Node)
	for i := range s.buttons {
		bw[i] = widget.NewSpace()
	}
	scores := make([]node.Node, len(s.scores))
	for i := range scores {
		scores[i] = stretch(widget.NewSpace(), 1)
	}
	bb := widget.NewPadder(widget.AxisBoth, unit.Pixels(padPx),
		widget.NewFlow(bAx,
			widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), bw["save"]),
			widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), bw["load"]),
			stretch(widget.NewSpace(), 1),
			stretch(widget.NewFlow(bAx, scores...), 1),
			stretch(widget.NewSpace(), 1),
			widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), bw["reset"]),
		),
	)
	// field
	f := widget.NewSpace()

	var all node.Node

	all = widget.NewFlow(aAx,
		stretch(bb, 0),
		stretch(widget.NewPadder(widget.AxisBoth, unit.Pixels(padPx), f), 1),
	)
	// Android tells us to use the whole screen, but the upper
	// part of it is covered by the system, so we need to get out
	// of the way of that.  The 10pt here is completely pulled out
	// of my ass. I have no idea what the actual measure is and
	// how to find what it is or how to disable it, the best I can
	// achieve is "It Works For Me".
	if runtime.GOOS == "android" {
		all = widget.NewPadder(widget.AxisVertical, unit.Points(10), all)
	}
	// do the layout.
	all.Measure(&t, e.WidthPx, e.HeightPx)
	all.Wrappee().Rect = image.Rectangle{Max: image.Pt(e.WidthPx, e.HeightPx)}
	all.Layout(&t)

	r := widgetScreenRect(f)
	dx, dy := r.Dx(), r.Dy()

	// square the field rectangle.
	if dx > dy {
		r.Min.X += dx - dy
	} else {
		r.Min.Y += dy - dx
	}

	s.fr = r

	s.ful = r.Min
	s.tsz.X = r.Dx() / s.f.W()
	s.tsz.Y = r.Dy() / s.f.H()

	for i := range bw {
		s.buttons[i].r = widgetScreenRect(bw[i])
	}
	for i := range scores {
		s.scores[i] = widgetScreenRect(scores[i])
	}

	s.tiles.SetSz(s.tsz, s.buttons["save"].r.Size(), s.scores[0].Size())
}
