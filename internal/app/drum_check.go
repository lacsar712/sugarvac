package app

import (
	"github.com/lacsar712/sugarvac/internal/model"
)

func (a *App) CheckMassecLevel(snap model.PlantSnapshot) error {
	if snap.Massec.LevelPercent < model.MinMassecLevelPercent {
		return model.LevelChain(snap.Massec.LevelPercent)
	}
	return nil
}
