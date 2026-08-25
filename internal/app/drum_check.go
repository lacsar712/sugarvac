package app

import (
	"fmt"

	"github.com/lacsar712/sugarvac/internal/model"
)

func (a *App) CheckMassecLevel(snap model.PlantSnapshot) error {
	if snap.Massec.LevelPercent < model.MinMassecLevelPercent {
		return fmt.Errorf("%w", model.ErrMassecLevelLow)
	}
	return nil
}
