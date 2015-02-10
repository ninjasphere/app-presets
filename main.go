package main

import (
	"fmt"
	"github.com/ninjasphere/app-presets/model"
	"github.com/ninjasphere/app-presets/service"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
)

var (
	info = ninja.LoadModuleInfo("./package.json")
)

//SchedulerApp describes the scheduler application.
type PresetsApp struct {
	support.AppSupport
	service *service.PresetsService
}

// Start is called after the ExportApp call is complete.
func (a *PresetsApp) Start(m *model.Presets) error {
	if a.service != nil {
		return fmt.Errorf("Service has already been started - the action has been ignored.")
	} else {
		if m == nil || m.Version == "" {
			m = &model.Presets{
				Version: Version,
			}
		}
		service := &service.PresetsService{
			Model: m,
			Save: func(m *model.Presets) {
				a.SendEvent("config", m)
			},
			Conn: a.Conn,
			Log:  a.Log,
		}
		service.Save(m)
		if err := service.Init(); err != nil {
			return err
		} else {
			a.service = service
		}
	}
	return nil
}

// Stop the scheduler module.
func (a *PresetsApp) Stop() error {
	tmp := a.service
	a.service = nil

	if tmp == nil {
		return fmt.Errorf("The service is not started - action has been ignored.")
	} else {
		tmp.Destroy()
	}
	return nil
}

func main() {
	app := &PresetsApp{}
	err := app.Init(info)
	if err != nil {
		app.Log.Fatalf("failed to initialize app: %v", err)
	}

	err = app.Export(app)
	if err != nil {
		app.Log.Fatalf("failed to export app: %v", err)
	}

	support.WaitUntilSignal()
}
