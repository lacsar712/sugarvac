package app

import (
	"context"
	"fmt"
	"time"

	"github.com/lacsar712/sugarvac/internal/clock"
	"github.com/lacsar712/sugarvac/internal/model"
)

func (a *App) advanceClock(d time.Duration) {
	if mc, ok := a.clk.(*clock.ManualClock); ok {
		mc.Advance(d)
		time.Sleep(time.Millisecond)
	} else {
		time.Sleep(d)
	}
}

func (a *App) bindSteamLoop(holder string, ctx context.Context) context.Context {
	a.mu.Lock()
	if cancel, ok := a.steamLoopCancels[holder]; ok {
		cancel()
	}
	child, cancel := context.WithCancel(ctx)
	a.steamLoopCancels[holder] = cancel
	a.mu.Unlock()
	return child
}

func (a *App) cancelSteamLoop(holder string) {
	a.mu.Lock()
	if cancel, ok := a.steamLoopCancels[holder]; ok {
		cancel()
		delete(a.steamLoopCancels, holder)
	}
	a.mu.Unlock()
}

func (a *App) cancelAllSteamLoops() {
	a.mu.Lock()
	for holder, cancel := range a.steamLoopCancels {
		cancel()
		delete(a.steamLoopCancels, holder)
	}
	a.mu.Unlock()
}

func (a *App) CoalFeedTPH() float64 {
	return a.Snapshot().Steamjet.SteamFlowTPH
}

func (a *App) RunSteamRamp(ctx context.Context, holder string, targetTPH float64) error {
	loopCtx := a.bindSteamLoop(holder, ctx)
	defer a.cancelSteamLoop(holder)
	for {
		if err := loopCtx.Err(); err != nil {
			return fmt.Errorf("%w", model.ErrContextDone)
		}
		snap := a.Snapshot()
		current := snap.Steamjet.SteamFlowTPH
		if current >= targetTPH {
			return nil
		}
		comb := snap.Steamjet
		comb.SteamFlowTPH = current + 1.0
		_ = a.store.UpdateSteamjet(a.cfg.UnitID, comb)
		a.telemetry.RecordCoalFeed(comb.SteamFlowTPH)
		a.advanceClock(100 * time.Millisecond)
	}
}

func (a *App) RunCoalFeed(ctx context.Context, holder string, steps int) error {
	loopCtx := a.bindSteamLoop(holder, ctx)
	defer a.cancelSteamLoop(holder)
	for i := 0; steps <= 0 || i < steps; i++ {
		snap := a.Snapshot()
		comb := snap.Steamjet
		comb.SteamFlowTPH += 0.5
		_ = a.store.UpdateSteamjet(a.cfg.UnitID, comb)
		a.telemetry.RecordCoalFeed(comb.SteamFlowTPH)
		a.advanceClock(100 * time.Millisecond)
	}
	return nil
}
