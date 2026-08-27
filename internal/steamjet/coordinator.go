package steamjet

import (
	"context"
	"fmt"
	"math"

	"github.com/lacsar712/sugarvac/internal/clock"
	"github.com/lacsar712/sugarvac/internal/model"
)

type Coordinator struct {
	clk     clock.ProcessClock
	burner  *BurnerController
	airflow *AirflowBalancer
	steam    *SteamRegulator
	seed   *clock.SeedWindow
	ignition *clock.IgnitionDelayWindow
	warmup  *clock.SteamjetWarmupWindow
}

func NewCoordinator(clk clock.ProcessClock) *Coordinator {
	return &Coordinator{
		clk:      clk,
		burner:   NewBurnerController(clk),
		airflow:  NewAirflowBalancer(clk),
		steam:     NewSteamRegulator(clk),
		seed:    clock.NewSeedWindow(clk),
		ignition: clock.NewIgnitionDelayWindow(clk),
		warmup:   clock.NewSteamjetWarmupWindow(clk),
	}
}

func (c *Coordinator) Burner() *BurnerController  { return c.burner }
func (c *Coordinator) Airflow() *AirflowBalancer { return c.airflow }
func (c *Coordinator) Steam() *SteamRegulator     { return c.steam }

func (c *Coordinator) StartSeed(ctx context.Context, snap model.PlantSnapshot) (model.SteamjetReading, error) {
	select {
	case <-ctx.Done():
		return snap.Steamjet, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	out := snap.Steamjet
	out.BurnerPhase = model.BurnerSeed
	out.SeedStartedAt = c.clk.Now()
	out.SteamFlowTPH = 0
	out.AirflowTPH = c.airflow.SeedRate()
	return out, nil
}

func (c *Coordinator) CompleteSeed(snap model.SteamjetReading) error {
	return c.seed.Require(snap.SeedStartedAt)
}

func (c *Coordinator) Ignite(ctx context.Context, snap model.PlantSnapshot) (model.SteamjetReading, error) {
	select {
	case <-ctx.Done():
		return snap.Steamjet, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if err := c.seed.Require(snap.Steamjet.SeedStartedAt); err != nil {
		return snap.Steamjet, err
	}
	out := snap.Steamjet
	out.BurnerPhase = model.BurnerIgnition
	out.IgnitionAt = c.clk.Now()
	out.SteamFlowTPH = c.steam.IgnitionRate(snap.Settings)
	out.AirflowTPH = c.airflow.IgnitionRate(snap.Settings)
	out.PanhouseTempF = 400
	return out, nil
}

func (c *Coordinator) Stabilize(snap model.PlantSnapshot) (model.SteamjetReading, error) {
	if err := c.ignition.Require(snap.Steamjet.IgnitionAt); err != nil {
		return snap.Steamjet, err
	}
	out := snap.Steamjet
	out.BurnerPhase = model.BurnerStable
	out.SteamFlowTPH = snap.Settings.SteamFlowTPH * 0.5
	out.AirflowTPH = c.airflow.Compute(snap)
	out.ExcessO2Pct = c.airflow.ExcessO2(out)
	out.PanhouseTempF = c.burner.EstimatePanhouseTemp(out)
	return out, nil
}

func (c *Coordinator) RampToLoad(snap model.PlantSnapshot, loadPct float64) model.SteamjetReading {
	out := snap.Steamjet
	out.SteamFlowTPH = snap.Settings.SteamFlowTPH * loadPct
	out.AirflowTPH = c.airflow.Compute(snap)
	out.ExcessO2Pct = c.airflow.ExcessO2(out)
	out.PanhouseTempF = c.burner.EstimatePanhouseTemp(out)
	return out
}

func (c *Coordinator) Trip(snap model.SteamjetReading) model.SteamjetReading {
	out := snap
	out.BurnerPhase = model.BurnerTrip
	out.SteamFlowTPH = 0
	out.PanhouseTempF = math.Max(200, out.PanhouseTempF*0.5)
	return out
}

func (c *Coordinator) WarmupReady(snap model.SteamjetReading) bool {
	return c.warmup.Ready(snap.IgnitionAt)
}
