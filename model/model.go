// Presets is the configuration for the app-presets app. It consists of
package model

// A ChannelState represents the state of a single channel.
type ChannelState struct {
	ID    string      `json:"id,omitempty"`
	State interface{} `json:"state,omitempty"`
}

// A ThingState represents the state of a single thing. It consists of the id of the thing,
// a list of channel states and a boolean which indicates whether the thing is
// included in the scene.
type ThingState struct {
	ID       string         `json:"id,omitempty"`
	Channels []ChannelState `json:"channels"`
}

// A Scene encodes the state of multiple things within a scope. It has a UUID that is a unique
// identifier the scene, a slot number, which is the position of the scene within a
// UI menu, a label which provides a human readable label for a scene, a scope which restricts the
// set of selectable things and a list of thing states.
type Scene struct {
	ID     string       `json:"id,omitempty"`
	Slot   int          `json:"slot,omitempty"`
	Label  string       `json:"label,omitempty"`
	Scope  string       `json:"scope,omitempty"`
	Things []ThingState `json:"things"`
}

// A Presets object is a collection of Scenes.
type Presets struct {
	Version string   `json:"version,omitempty"`
	Scenes  []*Scene `json:"scenes"`
}
