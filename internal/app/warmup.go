package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/sugarvac/internal/model"
)

func (a *App) WarmupStatus() (ready bool, detail string) {
	snap := a.Snapshot()
	if snap.Steamjet.SeedStartedAt.IsZero() {
		return false, "seed not started"
	}
	if !a.seedWindow.Ready(snap.Steamjet.SeedStartedAt) {
		return false, "seed window open"
	}
	if !snap.Steamjet.IgnitionAt.IsZero() && !a.warmupWindow.Ready(snap.Steamjet.IgnitionAt) {
		return false, "steamjet warmup window open"
	}
	if !snap.Massec.LastSwellAt.IsZero() {
		if err := a.massec.RequireSettled(snap.Massec); err != nil {
			return false, "massec swell settling"
		}
	}
	return true, "ready"
}

func (a *App) WaitWarmup(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", model.ErrContextDone)
		default:
		}
		ready, _ := a.WarmupStatus()
		if ready {
			return nil
		}
	}
}

func (a *App) SeedRemaining() string {
	snap := a.Snapshot()
	if snap.Steamjet.SeedStartedAt.IsZero() {
		return "not started"
	}
	if a.seedWindow.Ready(snap.Steamjet.SeedStartedAt) {
		return "complete"
	}
	return "in progress"
}

func (a *App) SteamjetWarmupRemaining() string {
	snap := a.Snapshot()
	if snap.Steamjet.IgnitionAt.IsZero() {
		return "not ignited"
	}
	if a.warmupWindow.Ready(snap.Steamjet.IgnitionAt) {
		return "complete"
	}
	return "in progress"
}
