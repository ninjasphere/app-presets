package service

import (
	"fmt"
	"github.com/ninjasphere/app-presets/model"
)

type PresetsService struct {
	Model *model.Presets
	Save  func(*model.Presets)
}

func (ps *PresetsService) Init() error {
	return nil
}

func (ps *PresetsService) Destroy() error {
	return nil
}

// see: http://schema.ninjablocks.com/service/presets#listPresetable
func (ps *PresetsService) ListPresetable(scope string) ([]model.ThingState, error) {
	return make([]model.ThingState, 0, 0), fmt.Errorf("unimplemented function: ListPresetable")
}

// see: http://schema.ninjablocks.com/service/presets#fetchScenes
func (ps *PresetsService) FetchScenes(scope string) ([]model.Scene, error) {
	return make([]model.Scene, 0, 0), fmt.Errorf("unimplemented function: FetchScenes")
}

// see: http://schema.ninjablocks.com/service/presets#fetchScene
func (ps *PresetsService) FetchScene(id string) (*model.Scene, error) {
	return nil, fmt.Errorf("unimplemented function: FetchScene")
}

// see: http://schema.ninjablocks.com/service/presets#storeScene
func (ps *PresetsService) StoreScene(model *model.Scene) error {
	return fmt.Errorf("unimplemented function: StoreScene")
}

// see: http://schema.ninjablocks.com/service/presets#applyScene
func (ps *PresetsService) ApplyScene(id string) error {
	return fmt.Errorf("unimplemented function: ApplyScene")
}
