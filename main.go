package main

import (
	"kk/kk"
	"log"

	"golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/paint"
)

func main() {
	smain()
}

func smain() {
	gldriver.Main(func(s screen.Screen) {
		st := kk.New()
		w, err := s.NewWindow(&screen.NewWindowOptions{800, 800})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()
		for {
			repaint, quit, publish := st.Handle(w.NextEvent())
			if quit {
				return
			}
			if repaint {
				w.Send(paint.Event{})
			}
			if publish {
				w.Publish()
			}
		}
	})
}

/*
func mmain() {
	app.Main(func(a app.App) {
		s := st()

		for e := range a.Events() {
			repaint, quit, publish := s.handle(a.Filter(e))
			if quit {
				return
			}
			if repaint {
				a.Send(paint.Event{})
			}
			if publish {
				a.Publish()
			}
		}
	})
}
*/
