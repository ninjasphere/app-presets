package service

import (
	"fmt"
	"github.com/ninjasphere/app-presets/model"

	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"strings"
)

type matchSpec struct {
	id    *string
	scope *string
	slot  *int
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
func (ps *PresetsService) parseScope(scope string) (string, string, string, error) {
	var err error
	room := ""
	siteID := ""

	parts := strings.Split(scope, ":")
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
				scope = fmt.Sprintf("site:%s", siteID)
			}
		default:
			err = fmt.Errorf("illegal argument: scope has an unrecognized scheme")
		}
	}
	if err != nil {
		ps.Log.Errorf("bad scope: %s: %v", scope, err)
	}
	return scope, room, siteID, err

}

// find the indicies of all matching scenes
func (ps *PresetsService) match(spec matchSpec) []int {
	found := make([]int, 0, len(ps.Model.Scenes))
	for i, m := range ps.Model.Scenes {

		// look for the index of all matching scenes

		if spec.scope != nil && m.Scope == *spec.scope {
			if spec.slot != nil {
				if m.Slot == *spec.slot {
					found = append(found, i)
					continue
				}
			} else {
				found = append(found, i)
				continue
			}
		}

		if spec.id != nil && m.ID == *spec.id {
			found = append(found, i)
			continue
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
