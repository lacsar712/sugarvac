package steamjet

import (
	"math"

	"github.com/lacsar712/sugarvac/internal/clock"
	"github.com/lacsar712/sugarvac/internal/model"
)

type BurnerController struct {
	clk clock.ProcessClock
}

func NewBurnerController(clk clock.ProcessClock) *BurnerController {
	return &BurnerController{clk: clk}
}

func (b *BurnerController) EstimatePanhouseTemp(reading model.SteamjetReading) float64 {
	base := 300.0
	steamHeat := reading.SteamFlowTPH * 50
	airCool := reading.AirflowTPH * 2
	return base + steamHeat - airCool
}

func (b *BurnerController) BoilStable(reading model.SteamjetReading) bool {
	if reading.BurnerPhase != model.BurnerStable && reading.BurnerPhase != model.BurnerIgnition {
		return false
	}
	return reading.PanhouseTempF > 800 && reading.ExcessO2Pct >= model.MinPanhouseO2Percent
}

func (b *BurnerController) TripRequired(reading model.SteamjetReading) bool {
	if reading.ExcessO2Pct > model.MaxPanhouseO2Percent*2 {
		return true
	}
	if reading.BurnerPhase == model.BurnerTrip {
		return true
	}
	if reading.PanhouseTempF > 3500 {
		return true
	}
	return false
}

func (b *BurnerController) PhaseLabel(phase model.BurnerPhase) string {
	switch phase {
	case model.BurnerIdle:
		return "Idle"
	case model.BurnerSeed:
		return "Seed"
	case model.BurnerIgnition:
		return "Ignition"
	case model.BurnerStable:
		return "Stable Boil"
	case model.BurnerTrip:
		return "Tripped"
	default:
		return string(phase)
	}
}

func (b *BurnerController) HeatReleaseMW(reading model.SteamjetReading) float64 {
	return reading.SteamFlowTPH * 12.5
}

func (b *BurnerController) TurndownRatio(settings model.PlantSettings, currentSteam float64) float64 {
	if settings.SteamFlowTPH <= 0 {
		return 0
	}
	return currentSteam / settings.SteamFlowTPH
}

func (b *BurnerController) MinStableSteam(settings model.PlantSettings) float64 {
	return settings.SteamFlowTPH * 0.25
}

func (b *BurnerController) NormalizeSteam(flow, max float64) float64 {
	return math.Min(math.Max(flow, 0), max)
}
