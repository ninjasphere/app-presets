package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/ninjasphere/app-presets/model"
	"github.com/ninjasphere/app-presets/service"
	"github.com/ninjasphere/go-ninja/config"
)

type PresetsRouter struct {
	presets *service.PresetsService
}

func NewPresetsRouter() *PresetsRouter {
	return &PresetsRouter{}
}

func (pr *PresetsRouter) Register(r martini.Router) {
	r.Get("/:id", pr.GetScene)
	r.Get("/prototype/site", pr.GetSitePrototype)
	r.Get("/prototype/room/:roomID", pr.GetRoomPrototype)
	r.Put("/:id", pr.PutScene)
	r.Delete("/:id", pr.DeleteScene)
	r.Post("/:id/apply", pr.ApplyScene)
	r.Get("", pr.GetScenes)
	r.Post("", pr.PutScene)
	r.Delete("", pr.DeleteScenes)
}

func writeResponse(code int, w http.ResponseWriter, response interface{}, err error) {
	if err == nil {
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(code)
		w.Write([]byte(fmt.Sprintf("error: %v\n", err)))
	}
}

func query(r *http.Request) *model.Query {
	result := &model.Query{}
	r.ParseForm()
	if scopes, ok := r.Form["scope"]; ok {
		result.Scope = &scopes[0]
	}
	if ids, ok := r.Form["id"]; ok {
		result.ID = &ids[0]
	}
	if slots, ok := r.Form["slot"]; ok {
		slot := 0
		if n, err := fmt.Sscanf(slots[0], "%d", &slot); n == 1 && err == nil {
			result.Slot = &slot
		}
	}
	return result
}

func (pr *PresetsRouter) GetScene(r *http.Request, w http.ResponseWriter, params martini.Params) {
	id := params["id"]
	scenes, err := pr.presets.FetchScenes(&model.Query{ID: &id})
	if scenes != nil && len(*scenes) == 1 {
		writeResponse(400, w, (*scenes)[0], err)
	} else {
		writeResponse(404, w, nil, err)
	}
}

func (pr *PresetsRouter) ApplyScene(r *http.Request, w http.ResponseWriter, params martini.Params) {
	scene, err := pr.presets.ApplyScene(params["id"])
	writeResponse(400, w, scene, err)
}

func (pr *PresetsRouter) GetScenes(r *http.Request, w http.ResponseWriter) {
	q := query(r)
	scenes, err := pr.presets.FetchScenes(q)
	writeResponse(400, w, scenes, err)
}

func (pr *PresetsRouter) PutScene(r *http.Request, w http.ResponseWriter, params martini.Params) {
	scene := &model.Scene{}
	json.NewDecoder(r.Body).Decode(scene)
	scene.ID = params["id"]
	r.ParseForm()
	if slots, ok := r.Form["slot"]; ok {
		slot := 0
		if n, err := fmt.Sscanf(slots[0], "%d", &slot); n == 1 && err == nil {
			scene.Slot = slot
		}
	}
	if labels, ok := r.Form["label"]; ok {
		scene.Label = labels[0]
	}
	scene, err := pr.presets.StoreScene(scene)
	writeResponse(400, w, scene, err)
}

func (pr *PresetsRouter) DeleteScene(r *http.Request, w http.ResponseWriter, params martini.Params) {
	id := params["id"]
	scenes, err := pr.presets.DeleteScenes(&model.Query{ID: &id})
	if scenes != nil && len(*scenes) == 1 {
		writeResponse(400, w, (*scenes)[0], err)
	} else {
		writeResponse(404, w, nil, err)
	}
}

func (pr *PresetsRouter) DeleteScenes(r *http.Request, w http.ResponseWriter, params martini.Params) {
	q := query(r)
	scenes, err := pr.presets.DeleteScenes(q)
	writeResponse(400, w, scenes, err)
}

func (pr *PresetsRouter) GetSitePrototype(r *http.Request, w http.ResponseWriter) {
	siteID := config.MustString("siteId")
	prototype, err := pr.presets.FetchScenePrototype(fmt.Sprintf("site:%s", siteID))
	writeResponse(400, w, prototype, err)
}

func (pr *PresetsRouter) GetRoomPrototype(r *http.Request, w http.ResponseWriter, params martini.Params) {
	prototype, err := pr.presets.FetchScenePrototype(fmt.Sprintf("room:%s", params["roomID"]))
	writeResponse(400, w, prototype, err)
}
