package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type FpsCounter struct {
	log map[int64]int
	mu  *sync.Mutex
}

func NewFpsCounter() *FpsCounter {
	fps := &FpsCounter{
		log: make(map[int64]int),
		mu:  &sync.Mutex{},
	}
	// prevent the log blowing up memory usage
	go fps.startCleanup(10 * time.Second)

	return fps
}

func (f *FpsCounter) Display(s tcell.Screen) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.log[time.Now().Unix()]++

	var fps float64
	for _, count := range f.log {
		fps += float64(count)
	}
	fps /= float64(len(f.log))

	emitStr(s, 0, 0, fmt.Sprintf("[%.2f FPS]", fps))
}

func (f *FpsCounter) startCleanup(ttl time.Duration) {
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()

	for range ticker.C {
		f.cleanup(ttl)
	}
}

func (f *FpsCounter) cleanup(ttl time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()

	expiry := time.Now().Add(-1 * ttl).Unix()

	// create a fresh map as maps don't shrink after deletes
	// https://github.com/golang/go/issues/20135
	fresh := make(map[int64]int)
	for t, count := range f.log {
		if t < expiry {
			continue
		}
		fresh[t] = count
	}
	f.log = fresh
}

func emitStr(s tcell.Screen, x, y int, str string) {
	for i, c := range str {
		s.SetContent(x+i, y, c, nil, tcell.StyleDefault)
	}
}
