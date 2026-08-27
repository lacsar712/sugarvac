package model

import "time"

func CloneSnapshot(s PlantSnapshot) PlantSnapshot {
	// PlantSnapshot is otherwise value-typed; only Alarms (a slice) shares
	// backing storage under a shallow copy. A "forwarded screenshot" must be a
	// detached copy so edits on the lab tablet cannot bleed into the live pan
	// curve or bypass the store's revision tracking. Re-allocate the slice so
	// in-place edits to the clone never mutate the original's backing array.
	if s.Alarms != nil {
		alarms := make([]AlarmEvent, len(s.Alarms))
		copy(alarms, s.Alarms)
		s.Alarms = alarms
	}
	return s
}

func DefaultSnapshot(unitID string) PlantSnapshot {
	now := time.Now()
	return PlantSnapshot{
		UnitID: unitID,
		State:  StateColdStandby,
		Settings: PlantSettings{
			Mode:              ModeBaseLoad,
			TargetMW:          150,
			TargetSteamPSI:    NormalSteamPressurePSI,
			MassecLevelSetpoint: 55,
			FeedwaterFlowTPH:  400,
			SteamFlowTPH:       35,
			ExcessO2Setpoint:  3.5,
		},
		Plant: PlantRef{UnitLabel: unitID, PlantCode: "STEAM-PLT"},
		Massec: MassecReading{
			LevelPercent: 50,
			Condition:    MassecNormal,
			FeedwaterTPH: 0,
			SteamFlowTPH: 0,
		},
		Steamjet: SteamjetReading{
			BurnerPhase: BurnerIdle,
		},
		Vacpan: VacpanReading{
			SteamPressurePSI: 0,
			SteamTempF:       70,
		},
		UpdatedAt: now,
	}
}

func (s PlantSnapshot) IsFiring() bool {
	return s.State == StateFiring || s.State == StateLoadFollow || s.State == StateRamp
}

func (s PlantSnapshot) MassecWithinLimits() bool {
	return s.Massec.LevelPercent >= MinMassecLevelPercent && s.Massec.LevelPercent <= MaxMassecLevelPercent
}

func (s PlantSnapshot) PressureWithinLimits() bool {
	if !s.IsFiring() {
		return true
	}
	return s.Vacpan.SteamPressurePSI <= MaxSteamPressurePSI
}
