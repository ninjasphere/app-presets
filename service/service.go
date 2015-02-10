package service

import (
	"fmt"
	"github.com/ninjasphere/app-presets/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
)

type PresetsService struct {
	Model       *model.Presets
	Save        func(*model.Presets)
	Conn        *ninja.Connection
	Log         *logger.Logger
	initialized bool
}

func (ps *PresetsService) Init() error {
	if ps.Log == nil {
		return fmt.Errorf("illegal state: no logger")
	}
	if ps.Model == nil {
		return fmt.Errorf("illegal state: Model is nil")
	}
	if ps.Save == nil {
		return fmt.Errorf("illegal state: Save is nil")
	}
	if ps.Conn == nil {
		return fmt.Errorf("illegal state: Conn is nil")
	}
	ps.initialized = true
	return nil
}

func (ps *PresetsService) Destroy() error {
	ps.initialized = false
	return nil
}

func (ps *PresetsService) checkInit() {
	if ps.Log == nil {
		ps.Log = logger.GetLogger("com.ninja.app-presets")
	}
	if !ps.initialized {
		ps.Log.Fatalf("illegal state: the service is not initialized")
	}
}

// see: http://schema.ninjablocks.com/service/presets#listPresetable
func (ps *PresetsService) ListPresetable(scope string) ([]*model.ThingState, error) {
	ps.checkInit()
	return make([]*model.ThingState, 0, 0), fmt.Errorf("unimplemented function: ListPresetable")
}

// see: http://schema.ninjablocks.com/service/presets#fetchScenes
func (ps *PresetsService) FetchScenes(scope string) ([]*model.Scene, error) {
	ps.checkInit()
	collect := make([]*model.Scene, 0, 0)
	for _, m := range ps.Model.Scenes {
		if m.Scope == scope {
			collect = append(collect, m)
		}
	}
	return collect, nil
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
