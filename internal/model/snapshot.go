package model

import "time"

func CloneSnapshot(s PlantSnapshot) PlantSnapshot {
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
