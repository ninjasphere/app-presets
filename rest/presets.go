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
	r.Post("/:id/apply", pr.ApplyScene)
	r.Get("", pr.GetScenes)
	r.Post("", pr.PutScene)
}

func writeResponse(code int, w http.ResponseWriter, response interface{}, err error) {
	if err == nil {
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(code)
		w.Write([]byte(fmt.Sprintf("error: %v\n", err)))
	}
}

func (pr *PresetsRouter) GetScene(r *http.Request, w http.ResponseWriter, params martini.Params) {
	scene, err := pr.presets.FetchScene(params["id"])
	writeResponse(400, w, scene, err)
}

func (pr *PresetsRouter) ApplyScene(r *http.Request, w http.ResponseWriter, params martini.Params) {
	err := pr.presets.ApplyScene(params["id"])
	writeResponse(400, w, nil, err)
}

func (pr *PresetsRouter) GetScenes(r *http.Request, w http.ResponseWriter) {
	r.ParseForm()
	var scope string
	if scopes, ok := r.Form["scope"]; !ok || len(scopes) == 0 || scopes[0] == "site" {
		siteID := config.MustString("siteId")
		scope = fmt.Sprintf("site/%s", siteID)
	} else {
		scope = scopes[0]
	}
	scenes, err := pr.presets.FetchScenes(scope)
	writeResponse(400, w, *scenes, err)
}

func (pr *PresetsRouter) PutScene(r *http.Request, w http.ResponseWriter, params martini.Params) {
	scene := &model.Scene{}
	json.NewDecoder(r.Body).Decode(scene)
	scene.ID = params["id"]

	scene, err := pr.presets.StoreScene(scene)
	writeResponse(400, w, scene, err)
}

func (pr *PresetsRouter) GetSitePrototype(r *http.Request, w http.ResponseWriter) {
	siteID := config.MustString("siteId")
	prototype, err := pr.presets.FetchScenePrototype(fmt.Sprintf("site/%s", siteID))
	writeResponse(400, w, prototype, err)
}

func (pr *PresetsRouter) GetRoomPrototype(r *http.Request, w http.ResponseWriter, params martini.Params) {
	prototype, err := pr.presets.FetchScenePrototype(fmt.Sprintf("room/%s", params["roomID"]))
	writeResponse(400, w, prototype, err)
}
