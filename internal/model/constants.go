package model

import "time"

const (
	DefaultLeaseTTL        = 30 * time.Second
	SeedWindow            = 5 * time.Minute
	IgnitionDelayWindow    = 15 * time.Second
	MassecSwellSettleWindow  = 45 * time.Second
	SteamjetWarmupWindow = 2 * time.Minute
	FeedwaterRampWindow    = 30 * time.Second
	MaxMassecLevelPercent    = 95.0
	MinMassecLevelPercent    = 15.0
	TripMassecLowPercent     = 10.0
	TripMassecHighPercent    = 98.0
	NormalSteamPressurePSI = 1800.0
	MaxSteamPressurePSI    = 2000.0
	MinPanhouseO2Percent    = 2.5
	MaxPanhouseO2Percent    = 6.0
	DefaultJournalCapacity = 512
)
