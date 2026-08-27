package clock

import (
	"testing"
	"time"
)

func TestCase(t *testing.T) {
	clk := NewManual(time.Unix(0, 0))
	p := NewPreheatWindow(clk)
	anchor := clk.Now()
	if p.Ready(anchor) {
		t.Fatal("preheat window should not be satisfied before process clock advances")
	}
	time.Sleep(50 * time.Millisecond)
	if p.Ready(anchor) {
		t.Fatal("preheat window closed on wall clock while process clock frozen")
	}
	clk.Advance(6 * time.Minute)
	if !p.Ready(anchor) {
		t.Fatal("preheat window should open after process clock advance")
	}
}
