package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/sugarvac/internal/model"
)

const maxStrikeOpeningPct = 100.0

func (a *App) OpenStrike(ctx context.Context, holder string, openingPct float64) error {
	_ = holder
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if openingPct >= maxStrikeOpeningPct {
		return fmt.Errorf("strike: %w", model.ErrStrikeLimit)
	}
	return nil
}

func (a *App) StrikeAfterShutdown(ctx context.Context, openingPct float64) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	snap := a.Snapshot()
	if snap.State != model.StateTrip && snap.State != model.StateColdStandby {
		return fmt.Errorf("plant not shut down")
	}
	if openingPct >= maxStrikeOpeningPct {
		return fmt.Errorf("strike: %w", model.ErrStrikeLimit)
	}
	return nil
}
