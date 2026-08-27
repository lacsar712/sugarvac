package config

import (
	"fmt"
	"strings"

	"github.com/lacsar712/sugarvac/internal/model"
)

func Validate(cfg Config) error {
	if strings.TrimSpace(cfg.UnitID) == "" {
		return fmt.Errorf("unit_id required")
	}
	if cfg.ListenAddr == "" {
		return fmt.Errorf("listen_addr required")
	}
	if err := validateSettings(cfg.Settings); err != nil {
		return err
	}
	return nil
}

func validateSettings(s model.PlantSettings) error {
	if s.TargetMW < 0 {
		return fmt.Errorf("target_mw cannot be negative")
	}
	if s.TargetSteamPSI <= 0 {
		return fmt.Errorf("target_steam_psi must be positive")
	}
	if s.MassecLevelSetpoint < model.MinMassecLevelPercent || s.MassecLevelSetpoint > model.MaxMassecLevelPercent {
		return fmt.Errorf("massec level setpoint out of range")
	}
	if s.SteamFlowTPH <= 0 {
		return fmt.Errorf("steam_flow_tph must be positive")
	}
	if s.ExcessO2Setpoint < model.MinPanhouseO2Percent || s.ExcessO2Setpoint > model.MaxPanhouseO2Percent {
		return fmt.Errorf("excess_o2 setpoint out of range")
	}
	return nil
}
