package store

import "github.com/lacsar712/sugarvac/internal/model"

type MassecSnapshotView struct {
	UnitID   string
	Massec     model.MassecReading
	Alarms   []model.AlarmEvent
	Revision uint64
}

func CloneMassecSnapshot(s model.PlantSnapshot) MassecSnapshotView {
	out := MassecSnapshotView{
		UnitID:   s.UnitID,
		Massec:     s.Massec,
		Revision: s.Revision,
	}
	out.Alarms = make([]model.AlarmEvent, len(s.Alarms))
	copy(out.Alarms, s.Alarms)
	return out
}
