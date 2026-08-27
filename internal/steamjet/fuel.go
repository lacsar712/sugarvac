package steamjet

import (
	"math"

	"github.com/lacsar712/sugarvac/internal/clock"
	"github.com/lacsar712/sugarvac/internal/model"
)

type SteamRegulator struct {
	clk clock.ProcessClock
}

func NewSteamRegulator(clk clock.ProcessClock) *SteamRegulator {
	return &SteamRegulator{clk: clk}
}

func (f *SteamRegulator) IgnitionRate(settings model.PlantSettings) float64 {
	return settings.SteamFlowTPH * 0.08
}

func (f *SteamRegulator) ComputeForLoad(settings model.PlantSettings, loadPct float64) float64 {
	loadPct = math.Max(0, math.Min(1, loadPct))
	return settings.SteamFlowTPH * loadPct
}

func (f *SteamRegulator) Ramp(current, target, maxStep float64) float64 {
	delta := target - current
	if math.Abs(delta) <= maxStep {
		return target
	}
	if delta > 0 {
		return current + maxStep
	}
	return current - maxStep
}

func (f *SteamRegulator) BtuPerHour(flowTPH float64) float64 {
	return flowTPH * 19_500_000
}

func (f *SteamRegulator) HeatInputMW(flowTPH float64) float64 {
	return flowTPH * 11.6
}

func (f *SteamRegulator) ValidatePermissive(settings model.PlantSettings, massecOK, seedOK bool) error {
	if !seedOK {
		return model.ErrSeedIncomplete
	}
	if !massecOK {
		return model.ErrMassecLevelTrip
	}
	if settings.SteamFlowTPH <= 0 {
		return model.ErrSteamPermissive
	}
	return nil
}

func (f *SteamRegulator) MinFlow(settings model.PlantSettings) float64 {
	return settings.SteamFlowTPH * 0.2
}

func (f *SteamRegulator) MaxFlow(settings model.PlantSettings) float64 {
	return settings.SteamFlowTPH * 1.1
}
