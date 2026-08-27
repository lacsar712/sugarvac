package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/lacsar712/sugarvac/internal/vacpan"
	"github.com/lacsar712/sugarvac/internal/clock"
	"github.com/lacsar712/sugarvac/internal/steamjet"
	"github.com/lacsar712/sugarvac/internal/config"
	"github.com/lacsar712/sugarvac/internal/massec"
	"github.com/lacsar712/sugarvac/internal/fsm"
	"github.com/lacsar712/sugarvac/internal/interlock"
	"github.com/lacsar712/sugarvac/internal/model"
	"github.com/lacsar712/sugarvac/internal/store"
)

type App struct {
	cfg           config.Config
	clk           clock.ProcessClock
	store         *store.PlantStore
	journal       *store.Journal
	fsm           *fsm.VacpanFSM
	vacpan        *vacpan.Controller
	steamjet    *steamjet.Coordinator
	massec          *massec.Coordinator
	interlock     *interlock.Interlock
	permissives   *interlock.PermissiveSet
	coordLock     *interlock.CoordinationLock
	scheduler     *clock.Scheduler
	seedWindow   *clock.SeedWindow
	warmupWindow  *clock.SteamjetWarmupWindow
	telemetry     *Telemetry
	tickCancels    map[string]context.CancelFunc
	steamLoopCancels map[string]context.CancelFunc
	mu             sync.RWMutex
}

func New(cfg config.Config, clk clock.ProcessClock) *App {
	return &App{
		cfg:          cfg,
		clk:          clk,
		store:        store.NewPlantStore(),
		journal:      store.NewJournal(cfg.JournalPath, cfg.JournalCapacity),
		fsm:          fsm.NewVacpanFSM(cfg.UnitID),
		vacpan:       vacpan.NewController(clk),
		steamjet:   steamjet.NewCoordinator(clk),
		massec:         massec.NewCoordinator(clk),
		interlock:    interlock.NewInterlock(cfg.LeaseTTL),
		permissives:  interlock.NewPermissiveSet(),
		coordLock:    interlock.NewCoordinationLock(),
		scheduler:    clock.NewScheduler(clk),
		seedWindow:  clock.NewSeedWindow(clk),
		warmupWindow: clock.NewSteamjetWarmupWindow(clk),
		telemetry:    NewTelemetry(cfg.UnitID),
		tickCancels:     make(map[string]context.CancelFunc),
		steamLoopCancels: make(map[string]context.CancelFunc),
	}
}

func (a *App) Snapshot() model.PlantSnapshot {
	snap, err := a.store.Require(a.cfg.UnitID)
	if err != nil {
		return model.DefaultSnapshot(a.cfg.UnitID)
	}
	return snap
}

func (a *App) Config() config.Config              { return a.cfg }
func (a *App) Clock() clock.ProcessClock          { return a.clk }
func (a *App) FSM() *fsm.VacpanFSM                { return a.fsm }
func (a *App) UnitID() string                     { return a.cfg.UnitID }
func (a *App) Store() *store.PlantStore           { return a.store }
func (a *App) Interlock() *interlock.Interlock    { return a.interlock }
func (a *App) Telemetry() TelemetrySnapshot       { return a.telemetry.Snapshot() }
func (a *App) Journal() *store.Journal            { return a.journal }

func (a *App) journalEvent(ev, payload string) {
	_, _ = a.journal.Append(a.cfg.UnitID, ev, payload)
}

func (a *App) syncState(state model.PlantState) {
	_ = a.store.UpdateState(a.cfg.UnitID, state)
}

func (a *App) isFiring(state model.PlantState) bool {
	return state == model.StateFiring || state == model.StateLoadFollow || state == model.StateRamp
}

func (a *App) refreshPermissives(snap model.PlantSnapshot) {
	a.permissives.SetMassec(a.massec.Level().WithinLimits(snap.Massec.LevelPercent))
	a.permissives.SetPressure(a.vacpan.Pressure().WithinTripLimits(snap.Vacpan.SteamPressurePSI, a.isFiring(snap.State)))
	a.permissives.SetSteamjet(a.steamjet.Burner().BoilStable(snap.Steamjet))
	a.permissives.SetSteam(snap.Steamjet.SteamFlowTPH > 0 || snap.State == model.StateSeed)
	a.permissives.SetIgnition(snap.Steamjet.BurnerPhase == model.BurnerStable || snap.Steamjet.BurnerPhase == model.BurnerIgnition)
	a.fsm.SetSteamPermissive(a.permissives.SteamOK())
	a.fsm.SetSeedComplete(a.seedWindow.Ready(snap.Steamjet.SeedStartedAt))
}

func (a *App) tickLabel() string {
	return fmt.Sprintf("%s-tick", a.cfg.UnitID)
}
