package main

import (
	"github.com/ninjasphere/app-presets/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
)

var (
	info = ninja.LoadModuleInfo("./package.json")
)

//SchedulerApp describes the scheduler application.
type PresetsApp struct {
	support.AppSupport
}

// Start is called after the ExportApp call is complete.
func (a *PresetsApp) Start(m *model.Presets) error {
	return nil
}

// Stop the scheduler module.
func (a *PresetsApp) Stop() error {
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
