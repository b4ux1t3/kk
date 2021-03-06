// +build !android

package main

import (
	"log"

	"github.com/art4711/kk/kk"

	"golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
)

func keyFilter(ei interface{}) interface{} {
	if e, ok := ei.(key.Event); ok && e.Direction == key.DirPress {
		switch e.Code {
		case key.CodeLeftArrow:
			return kk.EvL{}
		case key.CodeRightArrow:
			return kk.EvR{}
		case key.CodeUpArrow:
			return kk.EvU{}
		case key.CodeDownArrow:
			return kk.EvD{}
		case key.CodeEscape:
			return kk.EvQ{}
		}
	}
	return ei
}

func main() {
	gldriver.Main(func(s screen.Screen) {
		st := kk.New()
		w, err := s.NewWindow(&screen.NewWindowOptions{1080, 1776})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()
		for st.Handle(st.EvFilter(keyFilter(w.NextEvent())), func() { w.Publish() }) {
		}
	})
}
