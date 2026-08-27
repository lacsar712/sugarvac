package app

import (
	"fmt"

	"github.com/lacsar712/sugarvac/internal/model"
)

func (a *App) CheckMassecLevel(snap model.PlantSnapshot) error {
	if snap.Massec.LevelPercent < model.MinMassecLevelPercent {
		return fmt.Errorf("level low: %.1f", snap.Massec.LevelPercent)
	}
	return nil
}
