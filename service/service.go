package service

import (
	"code.google.com/p/go-uuid/uuid"
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
	ps.checkInit()
	for _, m := range ps.Model.Scenes {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, fmt.Errorf("No such scene: %s", id)
}

// see: http://schema.ninjablocks.com/service/presets#storeScene
func (ps *PresetsService) StoreScene(model *model.Scene) (*model.Scene, error) {
	ps.checkInit()
	var found int

	if model.Scope == "" {
		return nil, fmt.Errorf("illegal argument: model.Scope is empty")
	}

	for i, m := range ps.Model.Scenes {
		if model.ID == "" {
			if m.Scope == model.Scope && m.Slot == model.Slot {
				found = i
			}
		} else {
			if m.ID == model.ID {
				found = i
			}
		}
	}

	if model.ID == "" {
		model.ID = uuid.NewUUID().String()
	}

	if found >= len(ps.Model.Scenes) {
		ps.Model.Scenes = append(ps.Model.Scenes, model)
	} else {
		ps.Model.Scenes[found] = model
	}
	ps.Save(ps.Model)
	return model, nil
}

// see: http://schema.ninjablocks.com/service/presets#applyScene
func (ps *PresetsService) ApplyScene(id string) error {
	return fmt.Errorf("unimplemented function: ApplyScene")
}
