package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Pos struct {
	X, Y int
}

type Snow struct {
	flakes map[Pos]rune
	mu     *sync.Mutex
}

func NewSnow() *Snow {
	return &Snow{
		flakes: make(map[Pos]rune),
		mu:     &sync.Mutex{},
	}
}

func (s *Snow) AddFlake(x, y int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.flakes[Pos{x, y}] = '*'
}

func (s *Snow) Start(sc tcell.Screen, period time.Duration) {
	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for ;;<-ticker.C {
			w, _ := sc.Size()
			x := rand.Intn(w)
			s.AddFlake(x, 0)
		}
	}()
}

func (s *Snow) Display(sc tcell.Screen) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.update(sc.Size())
	for p, r := range s.flakes {
		sc.SetContent(p.X, p.Y, r, nil, tcell.StyleDefault)
	}
}

func (s *Snow) update(w, h int) {
	next := make(map[Pos]rune)
	for p := range s.flakes {
		// skip any flakes which are out of range
		if p.X < 0 || p.X >= w || p.Y < 0 || p.Y >= h {
			continue
		}

		p, r := s.fall(w, h, p)
		next[p] = r
	}
	s.flakes = next
}

func (s *Snow) fall(w, h int, p Pos) (Pos, rune) {
	below := Pos{p.X, p.Y + 1}
	left := Pos{p.X - 1, p.Y + 1}
	right := Pos{p.X + 1, p.Y + 1}

	_, snowBelow := s.flakes[below] 
	canFall := !snowBelow && below.Y < h
	_, snowLeft := s.flakes[left]
	canFallLeft := !snowLeft && left.Y < h && left.X >= 0
	_, snowRight := s.flakes[right]
	canFallRight := !snowRight && right.Y < h && right.X < w

	// update position
	if canFall {
		p = below
	} else if canFallLeft && !canFallRight {
		p = left
	} else if !canFallLeft && canFallRight {
		p = right
	} else if canFallLeft && canFallRight {
		p = [2]Pos{left, right}[rand.Intn(2)]
	}

	// update rune
	r := '*'
	if !canFall {
		r = '█'
	}
	return p, r
}
