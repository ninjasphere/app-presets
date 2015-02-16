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

// see: http://schema.ninjablocks.com/service/presets#fetchScenes
func (ps *PresetsService) FetchScenes(q *model.Query) (*[]*model.Scene, error) {
	ps.checkInit()
	if scope, _, _, err := ps.parseScope(q.Scope); err != nil {
		return nil, err
	} else {
		q.Scope = &scope
		found := ps.match(q)
		result := ps.copyScenes(found)
		return &result, nil
	}
}

// see: http://schema.ninjablocks.com/service/presets#deleteScenes
func (ps *PresetsService) DeleteScenes(q *model.Query) (*[]*model.Scene, error) {
	ps.checkInit()

	if scope, _, _, err := ps.parseScope(q.Scope); err != nil {
		return nil, err
	} else {
		q.Scope = &scope
		result := ps.deleteAll(ps.match(q))
		return &result, nil
	}
}

// see: http://schema.ninjablocks.com/service/presets#fetchScenePrototype
func (ps *PresetsService) FetchScenePrototype(scope string) (*model.Scene, error) {
	ps.checkInit()

	if scope == "" {
		scope = "site"
	}

	if scope, room, _, err := ps.parseScope(&scope); err != nil {
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
			ts := ps.createThingState(t)
			if ts != nil {
				result.Things = append(result.Things, *ts)
			}
		}
		return result, nil
	}
}

// see: http://schema.ninjablocks.com/service/presets#storeScene
func (ps *PresetsService) StoreScene(m *model.Scene) (*model.Scene, error) {
	ps.checkInit()

	if m.Scope == "" {
		m.Scope = "site"
	}

	if scope, _, _, err := ps.parseScope(&m.Scope); err != nil {
		return nil, err
	} else {
		m.Scope = scope
	}

	if m.ID == "" {
		m.ID = uuid.NewUUID().String()
	}

	if m.Slot <= 0 {
		m.Slot = 1
	}

	if m.Label == "" {
		m.Label = fmt.Sprintf("Preset %d", m.Slot)
	}

	found := ps.match(&model.Query{
		ID:    &m.ID,
		Scope: &m.Scope,
		Slot:  &m.Slot,
	})

	if len(found) > 1 {
		ps.deleteAll(found[1:])
	}

	if len(found) < 1 {
		ps.Model.Scenes = append(ps.Model.Scenes, m)
	} else {
		ps.Model.Scenes[found[0]] = m
	}

	ps.Save(ps.Model)
	return m, nil
}

// see: http://schema.ninjablocks.com/service/presets#applyScene
func (ps *PresetsService) ApplyScene(id string) (*model.Scene, error) {
	ps.checkInit()
	if id == "" {
		return nil, fmt.Errorf("illegal argument: id is empty")
	}
	if scenes, err := ps.FetchScenes(&model.Query{ID: &id}); err != nil || scenes == nil {
		return nil, err
	} else {
		thingClient := ps.Conn.GetServiceClient("$home/services/ThingModel")
		for _, scene := range *scenes {
			for i, t := range scene.Things {
				thing := &nmodel.Thing{}
				if err := thingClient.Call("fetch", []string{t.ID}, &thing, defaultTimeout); err != nil {
					ps.Log.Errorf("failed to obtain thing '%s': %v", id, err)
					continue
				}
				current := ps.createThingState(thing)
				t = *t.MergeUndoState(current)
				scene.Things[i] = t
				for _, c := range t.Channels {
					topic := fmt.Sprintf("$thing/%s/channel/%s", t.ID, c.ID)
					client := ps.Conn.GetServiceClient(topic)
					if err := client.Call("set", c.State, nil, defaultTimeout); err != nil {
						ps.Log.Warningf("Call to %s failed: %v", topic, err)
					}
				}
			}
			ps.Save(ps.Model)
			return scene, nil
		}
		return nil, fmt.Errorf("failed to find a matching scene: %s", id)
	}
}
