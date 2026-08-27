package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lacsar712/sugarvac/internal/app"
	"github.com/lacsar712/sugarvac/internal/clock"
	"github.com/lacsar712/sugarvac/internal/config"
	"github.com/lacsar712/sugarvac/internal/model"
)

func TestCase(t *testing.T) {
	clk := clock.NewManual(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	cfg := config.Default("FLAME-1")
	a, err := app.BootstrapWithClock(cfg, clk)
	if err != nil {
		t.Fatal(err)
	}
	comb := a.Snapshot().Steamjet
	comb.BurnerPhase = model.BurnerStable
	comb.PanhouseTempF = 400
	if err := a.Store().UpdateSteamjet(cfg.UnitID, comb); err != nil {
		t.Fatal(err)
	}
	err = a.OnBoilLoss(context.Background(), "maint-op")
	if err == nil {
		t.Fatal("expected thermal loss error")
	}
	if !errors.Is(err, model.ErrBoilLoss) {
		t.Fatalf("expected ErrBoilLoss, got %v", err)
	}
}
