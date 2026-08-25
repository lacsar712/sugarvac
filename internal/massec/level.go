package massec

import (
	"math"

	"github.com/lacsar712/sugarvac/internal/clock"
	"github.com/lacsar712/sugarvac/internal/model"
)

type LevelController struct {
	clk clock.ProcessClock
}

func NewLevelController(clk clock.ProcessClock) *LevelController {
	return &LevelController{clk: clk}
}

func (l *LevelController) Compute(snap model.PlantSnapshot, firing bool) (float64, model.MassecCondition) {
	level := snap.Massec.LevelPercent
	if !firing {
		return level, model.MassecNormal
	}
	balance := snap.Massec.FeedwaterTPH - snap.Massec.SteamFlowTPH
	level += balance * 0.01
	level = math.Max(model.MinMassecLevelPercent, math.Min(model.MaxMassecLevelPercent, level))
	cond := l.classify(level, snap)
	return level, cond
}

func (l *LevelController) classify(level float64, snap model.PlantSnapshot) model.MassecCondition {
	setpoint := snap.Settings.MassecLevelSetpoint
	if level > setpoint+15 {
		return model.MassecSwell
	}
	if level < setpoint-15 {
		return model.MassecShrink
	}
	if snap.Vacpan.SteamPressurePSI > snap.Settings.TargetSteamPSI*0.9 && level > setpoint+5 {
		return model.MassecCarry
	}
	return model.MassecNormal
}

func (l *LevelController) RecommendFeedwater(snap model.PlantSnapshot, firing bool) float64 {
	if !firing {
		return 0
	}
	err := snap.Settings.MassecLevelSetpoint - snap.Massec.LevelPercent
	return snap.Settings.FeedwaterFlowTPH + err*3
}

func (l *LevelController) WithinLimits(level float64) bool {
	return level >= model.MinMassecLevelPercent && level <= model.MaxMassecLevelPercent
}

func (l *LevelController) TripLow(level float64) bool  { return level < model.TripMassecLowPercent }
func (l *LevelController) TripHigh(level float64) bool { return level > model.TripMassecHighPercent }

func (l *LevelController) LevelError(snap model.PlantSnapshot) float64 {
	return snap.Massec.LevelPercent - snap.Settings.MassecLevelSetpoint
}

func (l *LevelController) ThreeElementBias(snap model.PlantSnapshot) float64 {
	steam := snap.Massec.SteamFlowTPH
	feed := snap.Massec.FeedwaterTPH
	levelErr := l.LevelError(snap)
	return feed + (steam-feed)*0.5 + levelErr*2
}
