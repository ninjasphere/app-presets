package service

import (
	"fmt"
	"github.com/ninjasphere/app-presets/model"

	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"strings"
)

type task struct {
	topic   string
	method  string
	payload interface{}
}

func (ps *PresetsService) worker() {
	for {
		select {
		case w := <-ps.queue:
			if w == nil {
				return
			}
			client := ps.Conn.GetServiceClient(w.topic)
			if err := client.Call(w.method, w.payload, nil, defaultTimeout); err != nil {
				ps.Log.Warningf("Call to %s of %s with %v failed: %v", w.method, w.topic, w.payload, err)
			}
		}
	}
}

// check that the service has been initialized
func (ps *PresetsService) checkInit() {
	if ps.Log == nil {
		ps.Log = logger.GetLogger("com.ninja.app-presets")
	}
	if !ps.initialized {
		ps.Log.Fatalf("illegal state: the service is not initialized")
	}
}

// make a copy of the channel's state, or nil if there is no such state
func copyState(ch *nmodel.Channel) interface{} {
	if ch.LastState != nil {
		if state, ok := ch.LastState.(map[string]interface{}); ok {
			if payload, ok := state["payload"]; ok {
				return payload
			}
		}
	}
	return nil
}

// parse a scope parameter and return the normalized form and the components
func (ps *PresetsService) parseScope(scope *string) (string, string, string, error) {
	var err error
	room := ""
	siteID := ""
	resultScope := ""

	if scope == nil || *scope == "" {
		return "", "", "", nil
	}
	resultScope = *scope
	parts := strings.Split(resultScope, ":")
	if len(parts) > 2 {
		err = fmt.Errorf("illegal argument: scope has too many parts")
	} else {
		if len(parts) == 0 {
			parts = []string{"site"}
		}
		switch parts[0] {
		case "room":
			room = parts[1]
		case "site":
			siteID = config.MustString("siteId")
			if len(parts) == 2 && parts[1] != siteID {
				err = fmt.Errorf("cannot configure presets for foreign site")
			} else {
				resultScope = fmt.Sprintf("site:%s", siteID)
			}
		default:
			err = fmt.Errorf("illegal argument: scope has an unrecognized scheme")
		}
	}
	if err != nil {
		ps.Log.Errorf("bad scope: %s: %v", scope, err)
	}
	return resultScope, room, siteID, err

}

// find the indicies of all matching scenes
func (ps *PresetsService) match(spec *model.Query) []int {
	found := make([]int, 0, len(ps.Model.Scenes))

	if spec.Scope != nil && *spec.Scope == "" {
		spec.Scope = nil
	}

	matchAll := spec.Scope == nil && spec.ID == nil && spec.Slot == nil

	for i, m := range ps.Model.Scenes {

		if matchAll {
			found = append(found, i)
		} else {
			// look for the index of all matching scenes

			if spec.Scope != nil && m.Scope == *spec.Scope {
				if spec.Slot != nil {
					if m.Slot == *spec.Slot {
						found = append(found, i)
						continue
					}
				} else {
					found = append(found, i)
					continue
				}
			}

			if spec.ID != nil && m.ID == *spec.ID {
				found = append(found, i)
				continue
			}
		}
	}
	return found
}

// make a copy of the specified scenes
func (ps *PresetsService) copyScenes(selection []int) []*model.Scene {
	result := make([]*model.Scene, len(selection))
	for i, x := range selection {
		result[i] = ps.Model.Scenes[x]
	}
	return result
}

// delete all the matching scenes
func (ps *PresetsService) deleteAll(selection []int) []*model.Scene {
	// no two scenes can have the same slot,scope or id.
	// delete the duplicates

	result := make([]*model.Scene, len(selection))

	j := 0
	k := 0
	for i, e := range ps.Model.Scenes {
		if j == len(selection) || i != selection[j] {
			if i != k {
				ps.Model.Scenes[k] = e
			}
			k++
		} else {
			result[j] = e
			j++
		}
	}
	ps.Model.Scenes = ps.Model.Scenes[0:k]
	return result
}

// create a ThingState object from a thing.
func (ps *PresetsService) createThingState(t *nmodel.Thing) *model.ThingState {
	if t.Device == nil || t.Device.Channels == nil {
		return nil
	}
	thingState := model.ThingState{
		ID:       t.ID,
		Channels: make([]model.ChannelState, 0, len(*t.Device.Channels)),
	}
Channels:
	for _, c := range *t.Device.Channels {

		for _, x := range excludedChannels {
			// don't include channels with excluded schema
			if x == c.Schema {
				continue Channels
			}
		}

		if c.SupportedMethods == nil {
			// don't include channels with no supported methods
			continue
		}

		found := false
		for _, m := range *c.SupportedMethods {
			found = (m == "set")
			if found {
				break
			}
		}
		if !found {
			// don't include channels that do not support the set method
			continue
		}
		state := copyState(c)
		if state == nil {
			continue
		}
		channelState := model.ChannelState{
			ID:    c.ID,
			State: state,
		}
		thingState.Channels = append(thingState.Channels, channelState)
	}

	if len(thingState.Channels) == 0 {
		return nil
	}

	return &thingState
}
