package kk

import (
	"image"

	"golang.org/x/exp/shiny/unit"
	"golang.org/x/exp/shiny/widget"
	"golang.org/x/exp/shiny/widget/node"
	"golang.org/x/exp/shiny/widget/theme"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/geom"
)

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
func widgetScreenRect(e *node.Embed) image.Rectangle {
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

	s.buttons = []butt{
		{b: Butt{Label: "Save"}},
		{b: Butt{Label: "Load"}},
		{b: Butt{Label: "Reset"}},
	}

	bw := make([]*widget.Uniform, len(s.buttons))
	for i := range bw {
		bw[i] = widget.NewUniform(theme.Light, nil)
	}

	bb := widget.NewUniform(theme.Neutral,
		widget.NewPadder(widget.AxisBoth, unit.Pixels(padPx),
			widget.NewFlow(bAx,
				widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), bw[0]),
				widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), bw[1]),
				stretch(widget.NewSpace(), 1),
				widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), bw[2]),
			),
		),
	)
	// field
	f := widget.NewUniform(theme.Light, nil)

	all := widget.NewFlow(aAx,
		stretch(bb, 0),
		stretch(widget.NewPadder(widget.AxisBoth, unit.Pixels(padPx), f), 1),
	)
	// do the layout.
	all.Measure(&t, e.WidthPx, e.HeightPx)
	all.Rect = image.Rectangle{Max: image.Pt(e.WidthPx, e.HeightPx)}
	all.Layout(&t)

	r := widgetScreenRect(&f.Embed)
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

	for i := range s.buttons {
		s.buttons[i].r = widgetScreenRect(&bw[i].Embed)
	}

	s.tiles.SetSz(s.tsz, s.buttons[0].r.Size())
}
