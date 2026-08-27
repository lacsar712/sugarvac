package fsm

import (
	"context"
	"fmt"
	"sync"

	"github.com/lacsar712/sugarvac/internal/model"
)

type VacpanFSM struct {
	mu            sync.RWMutex
	state         model.PlantState
	steamPermissive bool
	seedComplete  bool
	hooks          *HookChain
}

func NewVacpanFSM(unitID string) *VacpanFSM {
	_ = unitID
	return &VacpanFSM{state: model.StateColdStandby, hooks: NewHookChain()}
}

func (f *VacpanFSM) Hooks() *HookChain { return f.hooks }

func (f *VacpanFSM) State() model.PlantState {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.state
}

func (f *VacpanFSM) SetSteamPermissive(ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.steamPermissive = ok
}

func (f *VacpanFSM) SetSeedComplete(ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seedComplete = ok
}

func (f *VacpanFSM) SteamPermissive() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.steamPermissive
}

func (f *VacpanFSM) Dispatch(ctx context.Context, event PlantEvent) (model.PlantState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	select {
	case <-ctx.Done():
		return f.state, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if event == EvTrip {
		from := f.state
		if f.hooks != nil {
			if err := f.hooks.RunBefore(ctx, from, model.StateTrip, event); err != nil {
				return f.state, err
			}
		}
		f.state = model.StateTrip
		if f.hooks != nil {
			if err := f.hooks.RunAfter(ctx, from, model.StateTrip, event); err != nil {
				return f.state, err
			}
		}
		return f.state, nil
	}
	next, ok := NextState(f.state, event)
	if !ok {
		if f.hooks != nil {
			_ = f.hooks.RunAfter(ctx, f.state, f.state, event)
		}
		return f.state, fmt.Errorf("%s from %s: %w", event, f.state, ErrIllegalTransition)
	}
	if event == EvIgnite && !f.steamPermissive {
		return f.state, fmt.Errorf("%w", model.ErrSteamPermissive)
	}
	if event == EvSeedComplete && !f.seedComplete {
		return f.state, fmt.Errorf("%w", model.ErrSeedIncomplete)
	}
	from := f.state
	if f.hooks != nil {
		if err := f.hooks.RunBefore(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	f.state = next
	if f.hooks != nil {
		if err := f.hooks.RunAfter(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	return f.state, nil
}

func (f *VacpanFSM) ForceState(state model.PlantState) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = state
}
