package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/sugarvac/internal/model"
)

func (a *App) CoordinateLoadChange(ctx context.Context, holder string, targetMW float64) error {
	if err := a.coordLock.Acquire(holder); err != nil {
		return err
	}
	defer a.coordLock.Release(holder)

	snap := a.Snapshot()
	if snap.Settings.TargetMW <= 0 {
		return fmt.Errorf("invalid target MW")
	}
	loadPct := targetMW / snap.Settings.TargetMW
	if loadPct > 1.1 {
		return fmt.Errorf("load request exceeds rated capacity")
	}
	if err := a.massec.RequireSettled(snap.Massec); err != nil {
		return err
	}
	settings := snap.Settings
	settings.TargetMW = targetMW
	if err := a.store.UpdateSettings(a.cfg.UnitID, settings); err != nil {
		return err
	}
	return a.RampLoad(ctx, holder, loadPct)
}

func (a *App) CoordinateMassecLevel(ctx context.Context, holder string, setpoint float64) error {
	if err := a.coordLock.Acquire(holder); err != nil {
		return err
	}
	defer a.coordLock.Release(holder)

	if setpoint < model.MinMassecLevelPercent || setpoint > model.MaxMassecLevelPercent {
		return fmt.Errorf("massec setpoint out of range")
	}
	snap := a.Snapshot()
	settings := snap.Settings
	settings.MassecLevelSetpoint = setpoint
	feed := a.massec.CoordinateFeedwater(snap, a.isFiring(snap.State))
	settings.FeedwaterFlowTPH = feed
	return a.store.UpdateSettings(a.cfg.UnitID, settings)
}

func (a *App) CoordinateSteamjetTrim(ctx context.Context, holder string, o2Setpoint float64) error {
	if err := a.coordLock.Acquire(holder); err != nil {
		return err
	}
	defer a.coordLock.Release(holder)

	if o2Setpoint < model.MinPanhouseO2Percent || o2Setpoint > model.MaxPanhouseO2Percent {
		return fmt.Errorf("O2 setpoint out of range")
	}
	snap := a.Snapshot()
	settings := snap.Settings
	settings.ExcessO2Setpoint = o2Setpoint
	if err := a.store.UpdateSettings(a.cfg.UnitID, settings); err != nil {
		return err
	}
	comb := snap.Steamjet
	comb.AirflowTPH = a.steamjet.Airflow().TrimAirflow(comb.AirflowTPH, o2Setpoint, comb.ExcessO2Pct)
	comb.ExcessO2Pct = a.steamjet.Airflow().ExcessO2(comb)
	return a.store.UpdateSteamjet(a.cfg.UnitID, comb)
}

func (a *App) CoordinationHeld() bool { return a.coordLock.Held() }

func (a *App) PlantHealth() map[string]string {
	snap := a.Snapshot()
	a.refreshPermissives(snap)
	out := map[string]string{
		"state":      string(snap.State),
		"massec_ok":    fmt.Sprintf("%v", a.permissives.MassecOK()),
		"pressure_ok": fmt.Sprintf("%v", a.permissives.PressureOK()),
		"steamjet_ok": fmt.Sprintf("%v", a.permissives.SteamjetOK()),
		"steam_ok":    fmt.Sprintf("%v", a.permissives.SteamOK()),
	}
	ready, detail := a.WarmupStatus()
	out["warmup_ready"] = fmt.Sprintf("%v", ready)
	out["warmup_detail"] = detail
	return out
}
