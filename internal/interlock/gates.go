package interlock

import (
	"fmt"

	"github.com/lacsar712/sugarvac/internal/model"
)

type PermissiveSet struct {
	steamOK       bool
	ignitionOK   bool
	massecOK       bool
	pressureOK   bool
	steamjetOK bool
}

func NewPermissiveSet() *PermissiveSet { return &PermissiveSet{} }

func (p *PermissiveSet) SetSteam(ok bool)       { p.steamOK = ok }
func (p *PermissiveSet) SetIgnition(ok bool)   { p.ignitionOK = ok }
func (p *PermissiveSet) SetMassec(ok bool)       { p.massecOK = ok }
func (p *PermissiveSet) SetPressure(ok bool)   { p.pressureOK = ok }
func (p *PermissiveSet) SetSteamjet(ok bool) { p.steamjetOK = ok }

func (p *PermissiveSet) SteamOK() bool       { return p.steamOK }
func (p *PermissiveSet) IgnitionOK() bool   { return p.ignitionOK }
func (p *PermissiveSet) MassecOK() bool       { return p.massecOK }
func (p *PermissiveSet) PressureOK() bool   { return p.pressureOK }
func (p *PermissiveSet) SteamjetOK() bool { return p.steamjetOK }

func (p *PermissiveSet) AllFiring() bool {
	return p.steamOK && p.ignitionOK && p.massecOK && p.pressureOK && p.steamjetOK
}

func (p *PermissiveSet) CheckIgnition() error {
	if !p.steamOK {
		return model.PermissiveChain("ignition")
	}
	if !p.ignitionOK {
		return fmt.Errorf("%w", model.ErrIgnitionBlocked)
	}
	return nil
}

func CheckBoilLoss(reading model.SteamjetReading) error {
	if reading.BurnerPhase == model.BurnerStable && reading.PanhouseTempF < 600 {
		return fmt.Errorf("%w", model.ErrBoilLoss)
	}
	return nil
}

func (p *PermissiveSet) CheckFiring() error {
	if err := p.CheckIgnition(); err != nil {
		return err
	}
	if !p.massecOK {
		return fmt.Errorf("%w", model.ErrMassecLevelTrip)
	}
	if !p.pressureOK {
		return fmt.Errorf("%w", model.ErrPressureTrip)
	}
	if !p.steamjetOK {
		return fmt.Errorf("%w", model.ErrSteamjetTrip)
	}
	return nil
}

type CoordinationLock struct {
	holder string
	held   bool
}

func NewCoordinationLock() *CoordinationLock { return &CoordinationLock{} }

func (c *CoordinationLock) Acquire(holder string) error {
	if c.held {
		return fmt.Errorf("%w", model.ErrCoordinationLock)
	}
	c.holder = holder
	c.held = true
	return nil
}

func (c *CoordinationLock) Release(holder string) {
	if c.held && c.holder == holder {
		c.held = false
		c.holder = ""
	}
}

func (c *CoordinationLock) Require(holder string) error {
	if !c.held || c.holder != holder {
		return fmt.Errorf("%w", model.ErrCoordinationLock)
	}
	return nil
}

func (c *CoordinationLock) Held() bool { return c.held }
