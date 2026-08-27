package app

import (
	"context"
	"time"
)

func (a *App) RunWarmupBurnScheduler(ctx context.Context, ignitionAt time.Time) error {
	for !a.warmupWindow.Ready(ignitionAt) {
		if err := ctx.Err(); err != nil {
			return err
		}
		a.advanceClock(100 * time.Millisecond)
	}
	snap, err := a.store.Require(a.cfg.UnitID)
	if err != nil {
		return err
	}
	return a.scheduler.InstallBurnPlan(snap.Settings, "warmup-burn")
}

// CancelWarmupBurn voids the scheduled warmup burn plan at the console. It must
// retract the plan items from the scheduler timeline, otherwise the plan page
// keeps showing the voided burn tail even after the process card is stamped.
func (a *App) CancelWarmupBurn() {
	a.scheduler.CancelBurnPlan("warmup-burn")
}

func (a *App) SchedulerItemCount() int {
	return a.scheduler.ItemCount()
}
