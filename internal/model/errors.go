package model

import "errors"

var (
	ErrContextDone      = errors.New("operation cancelled")
	ErrPlantNotFound    = errors.New("plant unit not found")
	ErrLeaseHeld        = errors.New("interlock lease held by another operator")
	ErrLeaseMissing     = errors.New("interlock lease missing or expired")
	ErrGateBlocked      = errors.New("safety gate blocked")
	ErrSteamPermissive   = errors.New("steam permissive not satisfied")
	ErrIgnitionBlocked  = errors.New("ignition sequence blocked")
	ErrMassecLevelTrip    = errors.New("massec level trip condition")
	ErrPressureTrip     = errors.New("steam pressure trip condition")
	ErrSteamjetTrip   = errors.New("steamjet trip condition")
	ErrIllegalState     = errors.New("illegal plant state transition")
	ErrSnapshotStale    = errors.New("snapshot revision stale")
	ErrWindowOpen       = errors.New("timing window still open")
	ErrSeedIncomplete  = errors.New("panhouse seed incomplete")
	ErrCoordinationLock = errors.New("coordination lock held")
	ErrMassecLevelLow     = errors.New("massec level below low limit")
	ErrBoilLoss        = errors.New("panhouse boil lost")
	ErrStrikeLimit    = errors.New("strike valve at limit")
)
