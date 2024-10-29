package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
)

const FPS = 60

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}
	s.EnableMouse()

	fpsCounter := NewFpsCounter()
	snow := NewSnow()
	snow.Start(s, 10 * time.Millisecond)
	
	// start rendering
	displayers := []DisplayFunc{snow.Display, fpsCounter.Display}
	go func() {
		for range time.NewTicker(time.Second / FPS).C {
			display(s, displayers...)
		}
	}()
	// process events
	for {
		switch e := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			display(s, displayers...)
		case *tcell.EventKey:
			if e.Key() == tcell.KeyEscape || e.Key() == tcell.KeyCtrlC {
				s.Fini()
				os.Exit(0)
			}
		case *tcell.EventMouse:
			if e.Buttons() & tcell.ButtonPrimary == 0 {
				continue
			}
			snow.AddFlake(e.Position())
		}
	}
}

type DisplayFunc func(tcell.Screen)

func display(s tcell.Screen, displayFuncs ...DisplayFunc) {
	s.Clear()
	defer func() {
		s.Show()
	}()

	for _, df := range displayFuncs {
		df(s)
	}
}
