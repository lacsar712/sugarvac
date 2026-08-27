package vacpan

import (
	"math"

	"github.com/lacsar712/sugarvac/internal/clock"
	"github.com/lacsar712/sugarvac/internal/model"
)

type FeedwaterRegulator struct {
	clk    clock.ProcessClock
	ramp   *clock.FeedwaterRampTracker
	lastTPH float64
}

func NewFeedwaterRegulator(clk clock.ProcessClock) *FeedwaterRegulator {
	return &FeedwaterRegulator{clk: clk, ramp: clock.NewFeedwaterRampTracker(clk)}
}

func (f *FeedwaterRegulator) ComputeFlow(snap model.PlantSnapshot, firing bool) float64 {
	if !firing {
		return 0
	}
	target := snap.Settings.FeedwaterFlowTPH
	steamDemand := snap.Vacpan.MainSteamFlowTPH
	levelErr := snap.Settings.MassecLevelSetpoint - snap.Massec.LevelPercent
	correction := levelErr * 2.0
	desired := target + correction + steamDemand*0.1
	return math.Max(0, desired)
}

func (f *FeedwaterRegulator) ApplyRamp(targetTPH float64) (float64, error) {
	if f.ramp.Satisfied() {
		f.lastTPH = targetTPH
		f.ramp.Reset()
		return targetTPH, nil
	}
	step := (targetTPH - f.lastTPH) * 0.1
	f.lastTPH += step
	return f.lastTPH, f.ramp.Require()
}

func (f *FeedwaterRegulator) BalanceError(snap model.PlantSnapshot) float64 {
	return snap.Massec.FeedwaterTPH - snap.Massec.SteamFlowTPH
}

func (f *FeedwaterRegulator) RecommendSetpoint(snap model.PlantSnapshot) float64 {
	base := snap.Settings.FeedwaterFlowTPH
	if snap.Vacpan.MainSteamFlowTPH > 0 {
		return snap.Vacpan.MainSteamFlowTPH * 1.02
	}
	return base
}

func (f *FeedwaterRegulator) Reset() {
	f.lastTPH = 0
	f.ramp.Reset()
}
