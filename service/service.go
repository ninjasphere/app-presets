package service

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/ninjasphere/app-presets/model"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/rpc"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

var excludedChannels []string = []string{}

type Connection interface {
	ExportService(service interface{}, topic string, ann *nmodel.ServiceAnnouncement) (*rpc.ExportedService, error)
	GetServiceClient(serviceTopic string) *ninja.ServiceClient
}

type PresetsService struct {
	Model       *model.Presets
	Save        func(*model.Presets)
	Conn        Connection
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

	var err error
	siteID := config.MustString("siteId")
	topic := fmt.Sprintf("$site/%s/service/%s", siteID, "presets")
	announcement := &nmodel.ServiceAnnouncement{
		Schema: "http://schema.ninjablocks.com/service/presets",
	}
	if _, err = ps.Conn.ExportService(ps, topic, announcement); err != nil {
		return err
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

// see: http://schema.ninjablocks.com/service/presets#fetchScenes
func (ps *PresetsService) FetchScenes(scope string) (*[]*model.Scene, error) {
	ps.checkInit()
	if scope, _, _, err := ps.parseScope(scope); err != nil {
		return nil, err
	} else {
		collect := make([]*model.Scene, 0, 0)
		for _, m := range ps.Model.Scenes {
			if m.Scope == scope {
				collect = append(collect, m)
			}
		}
		return &collect, nil
	}
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

// see: http://schema.ninjablocks.com/service/presets#deleteScene
func (ps *PresetsService) DeleteScene(id string) (*model.Scene, error) {
	ps.checkInit()
	for i, m := range ps.Model.Scenes {
		if m.ID == id {
			ps.Model.Scenes = append(ps.Model.Scenes[0:i], ps.Model.Scenes[i+1:]...)
			return m, nil
		}
	}
	return nil, fmt.Errorf("No such scene: %s", id)
}

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

// see: http://schema.ninjablocks.com/service/presets#fetchScenePrototype
func (ps *PresetsService) FetchScenePrototype(scope string) (*model.Scene, error) {
	ps.checkInit()

	if scope, room, _, err := ps.parseScope(scope); err != nil {
		return nil, err
	} else {

		thingClient := ps.Conn.GetServiceClient("$home/services/ThingModel")
		things := make([]*nmodel.Thing, 0)
		keptThings := make([]*nmodel.Thing, 0, len(things))
		if err := thingClient.Call("fetchAll", nil, &things, defaultTimeout); err != nil {
			return nil, err
		}

		for _, t := range things {
			if !t.Promoted ||
				(room != "" && (t.Location == nil || *t.Location != room)) {
				continue
			}
			keptThings = append(keptThings, t)
		}
		result := &model.Scene{
			Scope:  scope,
			Things: make([]model.ThingState, 0, len(keptThings)),
		}
		for _, t := range keptThings {
			if t.Device == nil || t.Device.Channels == nil {
				continue
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

			if len(thingState.Channels) > 0 {
				result.Things = append(result.Things, thingState)
			}
		}
		return result, nil
	}
}

// see: http://schema.ninjablocks.com/service/presets#storeScene
func (ps *PresetsService) StoreScene(model *model.Scene) (*model.Scene, error) {
	ps.checkInit()

	if scope, _, _, err := ps.parseScope(model.Scope); err != nil {
		return nil, err
	} else {
		model.Scope = scope
	}

	if model.Label == "" {
		model.Label = fmt.Sprintf("Preset %d", model.Slot)
	}

	found := -1
	for i, m := range ps.Model.Scenes {
		if model.ID == "" {
			if m.Scope == model.Scope && m.Slot == model.Slot {
				found = i
				break
			}
		} else {
			if m.ID == model.ID {
				found = i
				break
			}
		}
	}

	if model.ID == "" {
		model.ID = uuid.NewUUID().String()
	}

	if found < 0 {
		ps.Model.Scenes = append(ps.Model.Scenes, model)
	} else {
		ps.Model.Scenes[found] = model
	}
	ps.Save(ps.Model)
	return model, nil
}

// see: http://schema.ninjablocks.com/service/presets#applyScene
func (ps *PresetsService) ApplyScene(id string) error {
	ps.checkInit()
	if scene, err := ps.FetchScene(id); err != nil {
		return err
	} else {
		for _, t := range scene.Things {
			for _, c := range t.Channels {
				topic := fmt.Sprintf("$thing/%s/channel/%s", t.ID, c.ID)
				client := ps.Conn.GetServiceClient(topic)
				if err := client.Call("set", c.State, nil, defaultTimeout); err != nil {
					ps.Log.Warningf("Call to %s failed: %v", topic, err)
				}
			}
		}

		return nil
	}
}
