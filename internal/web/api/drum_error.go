package api

import (
	"errors"

	"github.com/lacsar712/sugarvac/internal/model"
)

func classifyMassecError(err error) (string, bool) {
	if errors.Is(err, model.ErrMassecLevelLow) {
		return "massec_level_low", true
	}
	return "", false
}
