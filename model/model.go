// Package model describes a collection of tasks to be executed at the opening and closing
// of schedule windows.
package model

// A Schedule specifies a list of Tasks, a Location and a TimeZone
type Presets struct {
	Version string `json:"version,omitempty"`
}
